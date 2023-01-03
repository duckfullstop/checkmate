package main

import (
	"bufio"
	"fmt"
	"github.com/duckfullstop/checkmate/pkg/blackjack"
	"os"
	"strings"
)

var hitKeywords = []string{
	"hit",
	"take",
	"deal",
	"h",
}
var stickKeywords = []string{
	"stick",
	"stand",
	"stay",
	"s",
}

// contains is a helper function: searches sl for any instance of target, returning a boolean truthfulness value.
// Capitalisation normalised.
func contains(sl []string, target string) bool {
	// Normalise.
	target = strings.ToLower(target)
	for _, value := range sl {
		if value == target {
			return true
		}
	}
	return false
}

// PlayBlackjackSP plays a single round of Blackjack on stdout, with the player simply playing for themselves (single player).
// A win condition is simply getting 21 or less; the only way the player loses is if they go bust.
func PlayBlackjackSP(table *blackjack.Table, player *blackjack.Player) (err error) {
	if table == nil || player == nil {
		return ErrNilReference
	}
	fmt.Print("Dealing new table...")
	// Calling Deal() resets the table automatically, which for our use case is absolutely fine
	errs := table.Deal()
	if len(errs) != 0 {
		return fmt.Errorf("errors returned when dealing new table: %v", errs)
	}
	fmt.Print("...Table ready to play!\n")

	reader := bufio.NewReader(os.Stdin)

	// Gameplay loop - breaks out when the hand is completed.
	// SWEng: It may make more sense for this to loop over each hand. Discussion point?
	for {
		// Always print the hand first
		fmt.Print("Your hand:\n")
		// Unsafe lookup! Potentially check that this is safe to do / whether there is more than one hand.
		for _, v := range player.Hands[0].Cards {
			fmt.Printf(" - %s\n", v.String())
		}
		fmt.Print("----------\n")

		// Now check score and see if the hand is still valid for further play
		score, minScore, locked, valid := player.Hands[0].Score()
		if locked {
			if !valid {
				fmt.Printf("Bust! Your hand is worth %v\n", score)
				break
			}
			// This shouldn't ever fire as a stick in the previous iteration breaks the loop, but just in case...
			fmt.Printf("Sticking with a score of %v\n", score)
			break
		}
		// Special secret flow for if you get a natural 21
		// lint: gocritic suggests rewriting this to switch, I disagree and think this is more readable as an if statement imo
		if len(player.Hands[0].Cards) == 2 && score == 21 {
			fmt.Printf("Blackjack! Score: %v (you should probably stick, just saying)\n", score)
		} else if score != minScore {
			fmt.Printf("Score: %v (%v with aces counting as 1)\n", score, minScore)
		} else {
			fmt.Printf("Score: %v\n", score)
		}

		// Accept user input
		var endHand bool
		fmt.Print("Action ([h]it, [s]tick): ")
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			// Quick massage to replace CRLF with LF - makes this platform independent
			input = strings.ReplaceAll(input, "\n", "")
			if contains(hitKeywords, input) {
				err := player.Hands[0].Hit()
				if err != nil {
					return err
				}
				break
			} else if contains(stickKeywords, input) {
				fmt.Printf("Sticking with a score of %v\n", score)
				err := player.Hands[0].Stick()
				if err != nil {
					return err
				}
				endHand = true
				break
			}
			// We didn't get a valid input, be sad with the user and loop again
			fmt.Print("Invalid action! Choose one of [h]it, [s]tick:")
		}
		if endHand {
			break
		}
		// If the hand isn't over, continue execution by repeating the for loop
	}
	err = table.EndRound()
	if err != nil {
		return err
	}
	score, _, _, valid := player.Hands[0].Score()
	if valid {
		fmt.Printf("Congratulations on a score of %v!\n", score)
	} else {
		fmt.Printf("Commiserations on a score of %v!\n", score)
	}

	return
}
