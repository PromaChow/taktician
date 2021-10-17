package importptn

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/jmoiron/sqlx"
	"github.com/nelhage/taktician/playtak"
	"github.com/nelhage/taktician/ptn"

	"github.com/google/subcommands"

	_ "github.com/mattn/go-sqlite3" // we assume sqlite
)

type Command struct{}

func (*Command) Name() string     { return "import-ptn" }
func (*Command) Synopsis() string { return "Import PTNs from playtak DB" }
func (*Command) Usage() string {
	return `import-ptn GAMES.db
`
}

func (c *Command) SetFlags(flags *flag.FlagSet) {
}

const (
	ReportInterval = 1000
	Workers        = 4
)

func (c *Command) Execute(ctx context.Context, flag *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(flag.Args()) != 1 {
		log.Println("Must supply a game database")
		return subcommands.ExitUsageError
	}

	sql, err := sqlx.Open("sqlite3", flag.Arg(0))
	if err != nil {
		log.Fatal("open: ", err)
	}

	_, err = sql.Exec(createPTNTable)
	if err != nil {
		log.Fatal("create schema: ", err)
	}

	tx := sql.MustBegin()
	defer tx.Commit()
	cur, err := tx.Queryx(selectTODO)
	if err != nil {
		log.Fatal("query: ", err)
	}

	type result struct {
		game gameRow
		ptn  string
	}

	todo := make(chan gameRow)
	results := make(chan result)

	go func() {
		for cur.Next() {
			var game gameRow
			err := cur.StructScan(&game)
			if err != nil {
				log.Fatal("scan:", err)
			}
			todo <- game
		}
		close(todo)
	}()

	var grp errgroup.Group

	for j := 0; j < Workers; j++ {
		grp.Go(
			func() error {
				for game := range todo {
					ptn, err := importOne(&game)
					if ptn == "" {
						continue
					}
					if err != nil {
						log.Printf("could not import: id=%d err=%v", game.Id, err)
						continue
					}
					results <- result{game, ptn}
				}
				return nil
			})
	}

	go func() {
		if err := grp.Wait(); err != nil {
			log.Fatalf("import: %v", err)
		}
		close(results)
	}()

	i := 0
	for result := range results {
		_, err = tx.NamedExec(
			insertPTN, &ptnRow{Id: result.game.Id, PTN: result.ptn})
		if err != nil {
			log.Fatalf("insert id=%d err=%v ", result.game.Id, err)
		}
		i = i + 1
		if i%ReportInterval == 0 {
			log.Printf("%d...", i)
		}
	}

	return subcommands.ExitSuccess
}

func formatTags(g *gameRow) []ptn.Tag {
	/*
	 * [Site "PlayTak.com"]
	 * [Event "Online Play"]
	 * [Date "2017.02.05"]
	 * [Time "20:31:18"]
	 * [Player1 "nelhage"]
	 * [Player2 "Guest3179"]
	 * [Clock "20:0"]
	 * [Result "R-0"]
	 * [Size "5"]
	 */

	t := time.Unix(int64(g.Date)/1000, int64(g.Date%1000)*int64(time.Millisecond))
	tags := []ptn.Tag{
		ptn.Tag{"Site", "playtak.com"},
		ptn.Tag{"Date", t.Format("2006.01.02")},
		ptn.Tag{"Time", t.Format("15:04:05")},
		ptn.Tag{"Player1", g.PlayerWhite},
		ptn.Tag{"Player2", g.PlayerBlack},
		ptn.Tag{"Result", g.Result},
		ptn.Tag{"Size", strconv.Itoa(g.Size)},
	}
	if g.TimerTime != 0 {
		timer := time.Duration(g.TimerTime) * time.Second
		secs := int64(timer.Seconds())
		timestr := fmt.Sprintf("%02d:%02d", secs/60, secs%60)
		if g.TimerInc != 0 {
			timestr = fmt.Sprintf("%s +%d", timestr, g.TimerInc)
		}
		tags = append(tags, ptn.Tag{"Clock", timestr})
	}
	return tags
}

func importOne(g *gameRow) (string, error) {
	if g.Notation == "" {
		return "", nil
	}

	var out ptn.PTN
	out.Tags = formatTags(g)

	moves := strings.Split(g.Notation, ",")
	for i, mv := range moves {
		if i%2 == 0 {
			out.Ops = append(out.Ops, &ptn.MoveNumber{Number: i/2 + 1})
		}
		mv, err := playtak.ParseServer(strings.Trim(mv, " "))
		if err != nil {
			return "", fmt.Errorf("move %d: %v", i, err)
		}
		out.Ops = append(out.Ops, &ptn.Move{Move: mv})
	}

	return out.Render(), nil
}
