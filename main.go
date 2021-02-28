package main

import (
	"fmt"
	"math/rand"
)

type Char byte

type CommunityCards struct {
	Cards  []Card
}

func (t *CommunityCards) status() int {
	switch len(t.Cards) {
	case 0:
		return 0
	case 3:
		return 1
	case 4:
		return 2
	case 5:
		return 3
	default:
		panic("There is an unexpected number of cards on the table")
	}
}

type Card struct {
	Number int8
	Suit   Char
}

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

func extractCardFromDeck(deck *[]Card, i int) {
	deckLen := len(*deck)
	(*deck)[i] = (*deck)[deckLen-1] // Copy last element to index i.
	(*deck)[deckLen-1] = Card{}     // Erase last element (write zero value).
	(*deck) = (*deck)[:deckLen-1]   // Truncate slice.
}

func addHandToTable(hand Hand, deck *[]Card, hands *[]Hand) {
	*hands = append(*hands, hand)
	for _, card := range hand.Cards {
		for i, deck_card := range *deck {
			if card == deck_card {
				extractCardFromDeck(deck, i)
			}
		}
	}
}

func getRandomCardsFromDeck(deck *[]Card, nr int) []Card {
	var cards []Card
	for i := 0; i < nr; i++ {
		deckLen := len(*deck)
		pick := rand.Intn(deckLen)
		crd := (*deck)[pick]
		extractCardFromDeck(deck, pick)
		cards = append(cards, crd)
	}
	return cards
}

func casinoWorker(results chan<- int, jobs <-chan Game) {
	for work := range jobs {
		tableStatus := work.Table.status()
		mapping := getStatusMap()
		cardsLeftToPull := mapping[tableStatus]
		deck := work.Deck
		cardsPulled := getRandomCardsFromDeck(&deck, cardsLeftToPull)
		fmt.Println(cardsPulled)
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

	workers := 4
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

	fmt.Println(h1, deck, table)
}