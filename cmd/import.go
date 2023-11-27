package cmd

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No file(s) specified!")
			return
		}

		var games []*Game
		for _, arg := range args {
			g, err := parseFile(arg)
			if err != nil {
				log.Fatal(err)
			}
			games = append(games, g)
		}

		currentSelections := make(map[string][]string)

		for _, g := range games {
			// After each game, add the new selections to the next game's
			// exclusions to avoid matching the same people twice.
			for p, s := range currentSelections {
				g.Exclusions[p] = append(g.Exclusions[p], s...)
			}

			for _, p := range g.Players {
				// Get the available options by combining the player's
				// exclusions with any previous selections for the current
				// game.
				exclusions := g.Exclusions[p]
				for _, s := range g.Selections {
					exclusions = append(exclusions, s)
				}

				var opts []string
				for _, o := range g.Players {
					if !contains(exclusions, o) && p != o {
						opts = append(opts, o)
					}
				}

				if len(opts) == 0 {
					log.Fatalf("No available options for player %s", p)
				}

				i := rand.Intn(len(opts))
				g.Selections[p] = opts[i]

				currentSelections[p] = append(currentSelections[p], opts[i])
			}
		}

		for _, g := range games {
			fmt.Println(fmt.Sprintf("%s", g.ID))
			for p, s := range g.Selections {
				fmt.Println(fmt.Sprintf("\t%s\t%s", p, s))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

type Game struct {
	ID         string
	Players    []string
	Selections map[string]string
	Exclusions map[string][]string
}

func parseFile(path string) (*Game, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	if !sc.Scan() {
		return nil, fmt.Errorf("File is empty")
	}
	id := sc.Text()

	exclusions := make(map[string][]string)
	for sc.Scan() {
		line := strings.Split(sc.Text(), ",")
		exclusions[line[0]] = line[1:]
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	var players []string
	for p := range exclusions {
		players = append(players, p)
	}

	return &Game{
		ID:         id,
		Players:    players,
		Selections: make(map[string]string),
		Exclusions: exclusions,
	}, nil
}

func contains(s []string, a string) bool {
	for _, v := range s {
		if v == a {
			return true
		}
	}
	return false
}
