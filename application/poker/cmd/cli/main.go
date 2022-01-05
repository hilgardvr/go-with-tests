package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hilgardvr/go-with-tests/application/poker"
)

const dbFilename = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFilename)

	if err != nil {
		log.Fatal(err)
	}
	defer close()
	game := poker.NewCli(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter))
	game.PlayPoker()
}

