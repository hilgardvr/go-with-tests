package poker

import (
	"bufio"
	"io"
	"strings"
)

type Cli struct {
	playerStore PlayerStore
	in *bufio.Scanner 
}

func (cli *Cli)PlayPoker() {
	// reader := bufio.NewScanner(cli.in)
	// reader.Scan()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func NewCli(store PlayerStore, in io.Reader) *Cli {
	return &Cli{
		playerStore: store,
		in: bufio.NewScanner(in),
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *Cli)readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}