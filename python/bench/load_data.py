import tak.train

import sys

def main(args):
  tak.train.load_corpus(args[1])

if __name__ == '__main__':
  main(sys.argv)
