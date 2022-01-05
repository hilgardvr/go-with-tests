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

func (cli *Cli)PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *Cli) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000} 
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + 10 * time.Minute
	}
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