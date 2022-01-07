package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const PlayerPrompt = "Please enter the number of players: "

type Game struct {
	alerter BlindAlerter
	store PlayerStore
}

func NewGame(alerter BlindAlerter, store PlayerStore) *Game {
	return &Game{
		alerter: alerter,
		store: store,
	}
}

func (p *Game) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5 + numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000} 
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (g *Game) Finish(winner string) {
	g.store.RecordWin(winner)
}

type Cli struct {
	in *bufio.Scanner 
	out io.Writer
	game *Game
}

func (cli *Cli)PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	numberOfPlayerInput := cli.readLine()
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayerInput, "\n"))
	cli.game.Start(numberOfPlayers)
	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)
	cli.game.Finish(winner)
}

func NewCli(in io.Reader, out io.Writer, game *Game) *Cli {
	return &Cli{
		in: bufio.NewScanner(in),
		out: out,
		game: game,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *Cli)readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}