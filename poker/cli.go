package poker

import "io"

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func NewCLI(playerStore PlayerStore, in io.Reader) *CLI {
	return &CLI{playerStore: playerStore, in: in}
}

func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Chris")
}
