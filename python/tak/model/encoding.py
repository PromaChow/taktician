#!/usr/bin/env python
import torch

from .. import game, moves, pieces

MAX_RESERVES = 50
MAX_CAPSTONES = 2

MAX_SLIDES = 256


class Token:
    EMPTY = 0

    MY_TOP_FLAT = 1
    MY_FLAT = 2
    MY_STANDING = 3
    MY_CAPSTONE = 4

    THEIR_TOP_FLAT = 5
    THEIR_FLAT = 6
    THEIR_STANDING = 7
    THEIR_CAPSTONE = 8

    WHITE_TO_PLAY = 9
    BLACK_TO_PLAY = 10

    LAST_CAPSTONE_VALUE = 254
    CAPSTONES = list(
        range(LAST_CAPSTONE_VALUE - MAX_CAPSTONES + 1, LAST_CAPSTONE_VALUE + 1)
    )
    FIRST_CAPSTONES_VALUE = CAPSTONES[0]
    LAST_RESERVES_VALUE = FIRST_CAPSTONES_VALUE - 1
    RESERVES = list(
        range(LAST_RESERVES_VALUE - MAX_RESERVES + 1, LAST_RESERVES_VALUE + 1)
    )

    FIRST_RESERVES_VALUE = RESERVES[0]

    OUTPUT_SENTINEL = 255

    # [to_play
    #   my_reserves my_caps
    #   their_reserves their_caps
    #   board...
    # ]


TOP_PIECES = {
    (True, pieces.Kind.FLAT): Token.MY_TOP_FLAT,
    (False, pieces.Kind.FLAT): Token.THEIR_TOP_FLAT,
    (True, pieces.Kind.STANDING): Token.MY_STANDING,
    (False, pieces.Kind.STANDING): Token.THEIR_STANDING,
    (True, pieces.Kind.CAPSTONE): Token.MY_CAPSTONE,
    (False, pieces.Kind.CAPSTONE): Token.THEIR_CAPSTONE,
}


def encode(p: game.Position, include_sentinel: bool = True) -> list[int]:
    data = []
    if include_sentinel:
        data.append(Token.OUTPUT_SENTINEL)

    if p.to_move() == pieces.Color.WHITE:
        data.append(Token.WHITE_TO_PLAY)
    else:
        data.append(Token.BLACK_TO_PLAY)
    stones = p.stones

    data.append(Token.RESERVES[stones[p.to_move().value].stones])
    data.append(Token.CAPSTONES[stones[p.to_move().value].caps])
    data.append(Token.RESERVES[stones[p.to_move().flip().value].stones])
    data.append(Token.CAPSTONES[stones[p.to_move().flip().value].caps])

    for square in p.board:
        if len(square) == 0:
            data.append(Token.EMPTY)
            continue
        top, *stack = square
        data.append(TOP_PIECES[(top.color == p.to_move(), top.kind)])
        for flat in stack:
            data.append(
                Token.MY_FLAT if flat.color == p.to_move() else Token.THEIR_FLAT
            )
    return data


def decode(board: torch.Tensor) -> game.Position:
    i = 0
    if board[i] == Token.OUTPUT_SENTINEL:
        i += 1

    if board[i] == Token.WHITE_TO_PLAY:
        to_play = pieces.Color.WHITE
    else:
        to_play = pieces.Color.BLACK
    i += 1
    for _ in range(2):
        assert board[i].item() in Token.RESERVES
        i += 1
        assert board[i].item() in Token.CAPSTONES
        i += 1

    squares = []
    this_sq = None

    for sq in board[i:].numpy():
        if sq == Token.MY_FLAT:
            this_sq.append(pieces.Piece.cached(to_play, pieces.Kind.FLAT))
        elif sq == Token.THEIR_FLAT:
            this_sq.append(pieces.Piece.cached(to_play.flip(), pieces.Kind.FLAT))
        else:
            if this_sq is not None:
                squares.append(this_sq)

            if sq == Token.EMPTY:
                this_sq = []
            else:
                kind = {
                    Token.MY_CAPSTONE: pieces.Kind.CAPSTONE,
                    Token.THEIR_CAPSTONE: pieces.Kind.CAPSTONE,
                    Token.MY_STANDING: pieces.Kind.STANDING,
                    Token.THEIR_STANDING: pieces.Kind.STANDING,
                    Token.MY_TOP_FLAT: pieces.Kind.FLAT,
                    Token.THEIR_TOP_FLAT: pieces.Kind.FLAT,
                }[sq]
                color = to_play
                if sq in [
                    Token.THEIR_CAPSTONE,
                    Token.THEIR_STANDING,
                    Token.THEIR_TOP_FLAT,
                ]:
                    color = color.flip()
                this_sq = [pieces.Piece.cached(color, kind)]

    if this_sq is not None:
        squares.append(this_sq)
    size = int(len(squares) ** (1 / 2))
    assert size * size == len(squares), f"Got a bad number of squares: {len(squares)}"
    return game.Position.from_squares(
        game.Config(size=size), squares, 2 if to_play == pieces.Color.WHITE else 3
    )


def _encode_batch(
    inputs, encode_one, dtype=torch.float
) -> (torch.Tensor, torch.Tensor):
    lens = torch.empty((len(inputs),), dtype=torch.int)
    out = torch.zeros((len(inputs), 0), dtype=dtype)
    for (i, p) in enumerate(inputs):
        encoded = encode_one(p)
        if len(encoded) > out.size(1):
            tmp = torch.zeros((out.size(0), len(encoded)), dtype=out.dtype)
            tmp[:, : out.size(1)] = out
            out = tmp
        out[i, : len(encoded)] = torch.tensor(encoded, dtype=out.dtype)
        lens[i] = len(encoded)
    mask = torch.zeros_like(out, dtype=torch.bool)
    for i, l in enumerate(lens):
        mask[i, :l] = 1
    return out, mask


def encode_batch(
    positions, include_sentinel: bool = True
) -> (torch.Tensor, torch.Tensor):
    return _encode_batch(
        positions, lambda p: encode(p, include_sentinel), dtype=torch.uint8
    )


MOVES_BY_SIZE = [moves.all_moves_for_size(s) for s in range(7)]
MOVES_TO_ID = [{m: i for (i, m) in enumerate(moves)} for moves in MOVES_BY_SIZE]

MAX_MOVE_ID = len(MOVES_BY_SIZE[-1])


def encode_move(size: int, m: moves.Move) -> int:
    return MOVES_TO_ID[size][m]


def decode_move(size: int, m: int) -> moves.Move:
    moves = MOVES_BY_SIZE[size]
    return moves[m]


def n_moves_for_size(size: int) -> int:
    return len(MOVES_BY_SIZE[size])


def encode_moves_batch(size, moves) -> torch.Tensor:
    encode = MOVES_TO_ID[size].__getitem__
    return torch.tensor([encode(m) for m in moves], dtype=torch.long)
