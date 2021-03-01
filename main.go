package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Char byte

type CommunityCards struct {
	Cards  []Card
}

// Tells you the status of the gare
func (t *CommunityCards) status() int {
	switch len(t.Cards) {
	case 0: // pre-flop
		return 0
	case 3: //flop
		return 1
	case 4: // turn
		return 2
	case 5: // river
		return 3
	default:
		panic("There is an unexpected number of cards on the table")
	}
}

type Card struct {
	Number int8
	Suit   Char
}

type ByNumber []Card

// Implement sort interface for cards
func (a ByNumber) Len() int           { return len(a) }
func (a ByNumber) Less(i, j int) bool {
	if a[i].Number == 1{
		return false
	} else if a[j].Number == 1{
		return true
	}
	return int(a[i].Number) < int(a[j].Number)
}
func (a ByNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }


type Hand struct {
	Cards [2]Card
}

type Game struct {
	Table CommunityCards
	Hands []Hand
	Deck []Card
}

func (s Char) String() string {
	return fmt.Sprintf("%c", s)
}

func getAllNumbers() []int8 {
	return []int8{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	}
}

func getAllSuits() []Char {
	return []Char{
		'H', 'D', 'C', 'S',
	}
}

func createDeck() []Card {
	var deck []Card
	for _, s := range getAllSuits() {
		for _, n := range getAllNumbers() {
			deck = append(deck, Card{n, s})
		}
	}
	return deck
}

// Removes a single card from the deck
func extractCardFromSlice(deck *[]Card, i int) {
	deckLen := len(*deck)
	(*deck)[i] = (*deck)[deckLen-1] // Copy last element to index i.
	(*deck)[deckLen-1] = Card{}     // Erase last element (write zero value).
	(*deck) = (*deck)[:deckLen-1]   // Truncate slice.
}

// Takes the 2 player cards out of the deck
func addHandToTable(hand Hand, deck *[]Card, hands *[]Hand) {
	*hands = append(*hands, hand)
	for _, card := range hand.Cards {
		for i, deck_card := range *deck {
			if card == deck_card {
				extractCardFromSlice(deck, i)
			}
		}
	}
}

// Extracts n amount of cards from the deck
func getRandomCardsFromDeck(deck *[]Card, nr int) []Card {
	var cards []Card
	for i := 0; i < nr; i++ {
		deckLen := len(*deck)
		pick := rand.Intn(deckLen)
		crd := (*deck)[pick]
		extractCardFromSlice(deck, pick)
		cards = append(cards, crd)
	}
	return cards
}

// find 2, 3, or 4 of the same numbers on a slice of cards
func findMultipleSameNumbers(cards []Card, nr int) (map[int8]int8, bool) {
	store := make(map[int8]int8)
	for _, card := range cards {
		store[card.Number]++
	}
	for i, count := range store {
		if count != int8(nr) {
			delete(store, i)
		}
	}
	found := false
	if len(store) > 0 {
		found = true
	}
	return store, found
}

// Tries to find one pair
func checkMultiples(cards []Card, nr int) (int8, []Card) {
	var kickers []Card
	values, found := findMultipleSameNumbers(cards, 2)
	if !found {
		return 0, kickers
	}

	// Only keep the highest pair
	var max int8 = 0
	for i, _ := range values {
		if i > max && max != 1 {
			max = i
		}
	}

	// Return kickers
	for _, c := range cards {
		if c.Number != max {
			kickers = append(kickers, c)
		}
	}

	// Order kickers descending
	sort.Sort(sort.Reverse(ByNumber(kickers)))
	return max, kickers
}


// Tries to find two pairs
func checkTwoPairs(cards []Card) ([]int8, []Card) {
	twoPairs := []int8{}
	found, kickers := checkMultiples(cards, 2)
	if found == 0 {
		return twoPairs, kickers
	}
	secondFound, kickers := checkMultiples(kickers, 2)
	if secondFound == 0 {
		return twoPairs, kickers
	}
	twoPairs = append(twoPairs, found, secondFound)
	return twoPairs, kickers
}

func checkOnePair(cards []Card) (int8, []Card) {
	return checkMultiples(cards, 2)
}

func checkTrips(cards []Card) (int8, []Card) {
	return checkMultiples(cards, 3)
}

func checkPoker(cards []Card) (int8, []Card) {
	return checkMultiples(cards, 5)
}

// Retrieves scenarios from the job queue and crunches them
func casinoWorker(results chan<- int, jobs <-chan Game) {
	// Retrieve a single job (= one game)
	for work := range jobs {
		communityCards := work.Table.Cards
		tableStatus := work.Table.status()
		mapping := getStatusMap()
		deck := work.Deck
		cardsLeftToPull := mapping[tableStatus]
		cardsPulled := getRandomCardsFromDeck(&deck, cardsLeftToPull)
		communityCards = append(communityCards, cardsPulled...)

		// Calculate the best combination each player holds
		for _, hand := range work.Hands {
			handCards := hand.Cards
			var playerCardPool []Card = communityCards
			playerCardPool = append(playerCardPool, handCards[0], handCards[1])

			if len(playerCardPool) != 7 {
				panic("Player should have 7 cards available in total")
			}

			// Check pair
			foundPair, kickersPair := checkOnePair(playerCardPool)
			//found2, kickers2 := checkTwoPairs(playerCardPool)
			//found, kickers := checkTrips(playerCardPool)
			//found, kickers := checkPoker(playerCardPool)

			fmt.Println(foundPair, hand, kickersPair)

		}
		results <- 2
	}
}

// Tells you how many community cards we still need to pull from deck
func getStatusMap() map[int]int {
	mapping := make(map[int]int)
	mapping[0] = 5
	mapping[1] = 2
	mapping[2] = 1
	mapping[3] = 0
	return mapping
}


func main() {
	var deck []Card = createDeck()
	var hands []Hand

	h1 := Hand{
		[2]Card{
			Card{1, 'H'},
			Card{1, 'S'},
		},
	}
	h2 := Hand{
		[2]Card{
			Card{10, 'C'},
			Card{11, 'C'},
		},
	}

	addHandToTable(h1, &deck, &hands)
	addHandToTable(h2, &deck, &hands)

	var table CommunityCards = CommunityCards{}

	workers := 1
	simulations := 10

	resultsChannel := make(chan int, simulations)
	jobsChannel := make(chan Game, simulations)
	var results []int

	for i := 0; i < workers; i++ {
		go casinoWorker(resultsChannel, jobsChannel)
	}

	for i := 0; i < simulations; i++ {
		var setting Game = Game{
			Table: table,
			Hands: hands,
			Deck: deck,
		}
		jobsChannel <- setting
	}
	close(jobsChannel)
	for i := 0; i < simulations; i++ {
		results = append(results, <- resultsChannel)
	}
}