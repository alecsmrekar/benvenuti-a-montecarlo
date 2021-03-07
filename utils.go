package main

import "math/rand"

func getAllNumbers(doubleAce bool) []int8 {
	cards := []int8{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	}
	if doubleAce {
		cards = append(cards, 14)
	}
	return cards
}

func getAllSuits() []Char {
	return []Char{
		'H', 'D', 'C', 'S',
	}
}

// Creates a 52 cards deck
func createDeck() []Card {
	var deck []Card
	for _, s := range getAllSuits() {
		for _, n := range getAllNumbers(false) {
			deck = append(deck, Card{n, s})
		}
	}
	return deck
}

// Removes one card from the passed deck of cards
func removeCardFromSlice(s *[]Card, i int) {
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}

// Takes the 2 player cards out of the deck
func addHandToTable(hand Hand, deck *[]Card, hands *[]Hand) {
	*hands = append(*hands, hand)
	for _, card := range hand.Cards {
		addCardToTable(card, deck)
	}
}

func addCardToTable(card Card, deck *[]Card) {
	for i, deck_card := range *deck {
		if card == deck_card {
			removeCardFromSlice(deck, i)
		}
	}
}

// Extracts n amount of cards from the deck
func getRandomCardsFromDeck(deck *[]Card, nr int) ([]Card) {
	var cards []Card
	for i := 0; i < nr; i++ {
		deckLen := len(*deck)
		pick := rand.Intn(deckLen)
		crd := (*deck)[pick]
		removeCardFromSlice(deck, pick)
		cards = append(cards, crd)
	}
	return cards
}

// Gets a human-readable combination name
func getCombinationName (input int8) (string) {
	combos := getCombinations()
	mapping := map[int8]string{
		combos.StraightFlush: "Straight Flush",
		combos.Poker: "Poker",
		combos.FullHouse: "Full House",
		combos.Flush: "Flush",
		combos.Straight: "Straight",
		combos.Trips: "Trips",
		combos.TwoPairs: "Two Pairs",
		combos.OnePair: "One Pair",
		combos.HighCard: "High Card",
	}
	return mapping[input]
}

func getCombinations () Combinations {
	return  Combinations{
		StraightFlush: 1,
		Poker:         2,
		FullHouse:     3,
		Flush:         4,
		Straight:      5,
		Trips:         6,
		TwoPairs:      7,
		OnePair:       8,
		HighCard:      9,
	}
}

// Checks if the deck has duplicate cards
func checkDeckHealth(deck []Card) {
	store := make(map[Card]int)
	for _, crd := range deck {
		store[crd]++
		if store[crd] > 1 {
			panic("Deck has duplicate cards")
		}
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