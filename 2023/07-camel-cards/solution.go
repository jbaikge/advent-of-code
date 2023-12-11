package camelcards

import (
	"bufio"
	"bytes"
	_ "embed"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/solutions"
)

const (
	TypeHighCard int = iota
	TypeOnePair
	TypeTwoPair
	TypeThreeKind
	TypeFullHouse
	TypeFourKind
	TypeFiveKind
)

//go:embed test.txt
var testData []byte

//go:embed input.txt
var inputData []byte

var cardMap = map[byte]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

func init() {
	solutions.Register(new(Solution))
}

type Hand struct {
	Cards      string
	Bid        int
	Values     [5]int
	WildValues [5]int
}

func NewHand(cards string, bid int) (hand Hand) {
	hand.Cards = cards
	hand.Bid = bid
	for i := range cards {
		value := cardMap[cards[i]]
		hand.Values[i] = value

		if value == cardMap['J'] {
			value = 1
		}
		hand.WildValues[i] = value

	}
	return
}

func (hand Hand) Type() int {
	counts := make(map[byte]int)
	for i := range hand.Cards {
		ch := hand.Cards[i]
		if _, found := counts[ch]; !found {
			counts[ch] = 0
		}
		counts[ch]++
	}

	switch len(counts) {
	case 1:
		return TypeFiveKind
	case 2:
		for _, count := range counts {
			if count == 4 {
				return TypeFourKind
			}
			if count == 3 {
				return TypeFullHouse
			}
		}
		log.Fatalf("len = 2 but not 4 of a kind or full house: %v", hand.Cards)
	case 3:
		pairs := 0
		for _, count := range counts {
			if count == 3 {
				return TypeThreeKind
			}
			if count == 2 {
				pairs++
			}
		}
		if pairs == 2 {
			return TypeTwoPair
		}
		log.Fatalf("len = 3 but not 3 of a kind: %v", hand.Cards)
	case 4:
		return TypeOnePair
	}
	return TypeHighCard
}

func (hand Hand) WildType() int {
	jokers := 0
	for i := range hand.Cards {
		if hand.Cards[i] == 'J' {
			jokers++
		}
	}

	switch hand.Type() {
	case TypeFiveKind:
		return TypeFiveKind
	case TypeFourKind:
		if jokers > 0 {
			return TypeFiveKind
		}
		return TypeFourKind
	case TypeFullHouse:
		if jokers == 2 || jokers == 3 {
			return TypeFiveKind
		}
		return TypeFullHouse
	case TypeThreeKind:
		if jokers == 1 {
			return TypeFourKind
		}
		if jokers == 2 {
			return TypeFiveKind
		}
		if jokers == 3 {
			return TypeFourKind
		}
		return TypeThreeKind
	case TypeTwoPair:
		if jokers == 1 {
			return TypeFullHouse
		}
		if jokers == 2 {
			return TypeFourKind
		}
		return TypeTwoPair
	case TypeOnePair:
		if jokers == 1 {
			return TypeThreeKind
		}
		if jokers == 2 {
			return TypeThreeKind
		}
		return TypeOnePair
	case TypeHighCard:
		if jokers == 1 {
			return TypeOnePair
		}
		return TypeHighCard
	}
	return TypeHighCard
}

func (hand Hand) Less(other Hand) int {
	if type1, type2 := hand.Type(), other.Type(); type1 != type2 {
		return type1 - type2
	}
	for i := 0; i < 5; i++ {
		if v1, v2 := hand.Values[i], other.Values[i]; v1 != v2 {
			return v1 - v2
		}
	}
	return 0
}

func (hand Hand) WildLess(other Hand) int {
	if type1, type2 := hand.WildType(), other.WildType(); type1 != type2 {
		return type1 - type2
	}
	for i := 0; i < 5; i++ {
		if v1, v2 := hand.WildValues[i], other.WildValues[i]; v1 != v2 {
			return v1 - v2
		}
	}
	return 0
}

type Solution struct {
	Hands []Hand
}

func (*Solution) Meta() solutions.Meta {
	return solutions.Meta{
		Name:    "Camel Cards",
		Year:    2023,
		Problem: 7,
		Datas: []solutions.Data{
			{
				Name:    "Test",
				Input:   testData,
				Expect1: 6440,
				Expect2: 5905,
			},
			{
				Name:  "Input",
				Input: inputData,
			},
		},
	}
}

func (s *Solution) Parse(data []byte) (err error) {
	s.Hands = make([]Hand, 0, 1000)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var bid int

		fields := strings.Fields(scanner.Text())
		bid, err = strconv.Atoi(fields[1])
		if err != nil {
			return
		}
		s.Hands = append(s.Hands, NewHand(fields[0], bid))
	}
	return
}

func (s *Solution) Part1() (answer int, err error) {
	hands := make([]Hand, len(s.Hands))
	copy(hands, s.Hands)

	slices.SortFunc[[]Hand, Hand](hands, func(a Hand, b Hand) int {
		return a.Less(b)
	})

	for i, hand := range hands {
		answer += (i + 1) * hand.Bid
	}
	return
}

func (s *Solution) Part2() (answer int, err error) {
	hands := make([]Hand, len(s.Hands))
	copy(hands, s.Hands)

	slices.SortFunc[[]Hand, Hand](hands, func(a Hand, b Hand) int {
		return a.WildLess(b)
	})

	for i, hand := range hands {
		answer += (i + 1) * hand.Bid
	}
	return
}
