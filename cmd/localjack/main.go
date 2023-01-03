package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/duckfullstop/checkmate/pkg/blackjack"
	"os"
	"strings"
)

func Initialise(deckCount int) (table *blackjack.Table, player *blackjack.Player, err error) {
	table = blackjack.NewTable(deckCount)
	if err != nil {
		return nil, nil, err
	}
	player = blackjack.NewPlayer()
	err = table.Join(player)
	// Shorthand on the return, saves a needless if err check
	return table, player, err
}

func main() {
	var decks int
	flag.IntVar(&decks, "decks", 1, "Number of decks to draw from.")

	flag.Parse()

	if decks < 1 {
		fmt.Printf("invalid number of decks - must be more than one!")
		os.Exit(2)
	}

	deck, player, err := Initialise(decks)
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	for {
		err = PlayBlackjackSP(deck, player)
		if err != nil {
			fmt.Printf("execution error! %s", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		var playAgain bool
		fmt.Print("Do you want to play again? [y]es, [n]o: ")
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("error: %s", err)
				os.Exit(1)
			}
			// Quickly massage input
			input = strings.ToLower(strings.ReplaceAll(input, "\n", ""))

			// Yes, we're cheating to determine the boolean outcome
			if strings.ContainsAny(input, "y") {
				// Continue the loop!
				playAgain = true
				break
			} else if strings.ContainsAny(input, "n") {
				// Redefinition not needed, done for clarity.
				playAgain = false
				break
			}
			// We didn't get a valid input, be sad with the user and loop again
			fmt.Print("Invalid action! Choose one of [y]es, [n]o: ")
		}
		if !playAgain {
			break
		}
	}
	fmt.Print("Thanks for playing! 💙")
}
