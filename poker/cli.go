package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func NewCLI(playerStore PlayerStore, in io.Reader) *CLI {
	return &CLI{playerStore: playerStore, in: in}
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}
func extractWinner(input string) string {

	return strings.Replace(input, " wins", "", 1)
}
