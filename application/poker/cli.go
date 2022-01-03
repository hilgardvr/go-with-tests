package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type Cli struct {
	playerStore PlayerStore
	in *bufio.Scanner 
	alerter BlindAlerter
}

type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int)
}

func (cli *Cli)PlayPoker() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000} 
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + 10 * time.Minute
	}
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func NewCli(store PlayerStore, in io.Reader, alerter BlindAlerter) *Cli {
	return &Cli{
		playerStore: store,
		in: bufio.NewScanner(in),
		alerter: alerter,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *Cli)readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}