package namegen

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// The ranks and suits of this name generator.
var (
	Ranks = []string{
		"one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
		"ace", "king", "page", "princess", "queen", "jack",
		"king", "magus", "prince", "knight", "challenger",
		"daughter", "son", "priestess", "shaman",
	}

	Suits = []string{
		"clubs", "hearts", "spades", "diamonds", // common playing cards
		"swords", "cups", "pentacles", "wands", // tarot
		"disks",                         // thoth tarot
		"coins",                         // karma
		"earth", "wind", "water", "air", // classical elements
		"aether", "spirits", "nirvana", // new age sounding things
		"chakras", "dilutions", "goings",
	}
)

// Next creates a name.
func Next() string {
	rank := Ranks[rand.Int()%len(Ranks)]
	suit := Suits[rand.Int()%len(Suits)]

	return fmt.Sprintf("%s-of-%s-%d", rank, suit, rand.Int63()%100000)
}
