package main

import (
	"fmt"
	"testing"
)

// use go test -cover to get code coverage

func TestCommunityCardStatus(t *testing.T) {
	cards := CommunityCards{[]Card{
		Card{1, 'C'},
		Card{1, 'D'},
	},
	}
	assertPanic(t, func() {
		cards.status()
	})

	cards = CommunityCards{[]Card{
		Card{1, 'C'},
		Card{1, 'D'},
		Card{1, 'D'},
	},
	}

	if cards.status() != 1 {
		t.Error("Community Card status should be 1, (flop)")
	}
}

func TestIfThereIsAnExtraAceInDeck(t *testing.T) {
	expected := []int8{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	}
	if !EqualInt8Slice(getAllNumbers(true), expected) {
		t.Error("Ace is missing from the deck")
	}
}

func TestDeckHealth(t *testing.T) {
	deck := createDeck()
	healthy := func() {
		checkDeckHealth(deck)
	}
	assertNoPanic(t, healthy)

	if len(deck) != 52 {
		t.Error("Deck is not healthy")
	}

	deck[1] = deck[0]
	healthy = func() {
		checkDeckHealth(deck)
	}
	assertPanic(t, healthy)
}

func TestNumberCompare(t *testing.T) {
	outcomes := getOutcomes()
	p1 := []int8{1,2}
	p2 := []int8{12,11}
	outcome := numberCompare(p1,p2)
	if outcome != outcomes.Win {
		t.Error("Outcome assessment is wrong")
	}
	p1 = []int8{2,3,4}
	p2 = []int8{10,11,12}
	outcome = numberCompare(p1,p2)
	if outcome != outcomes.Lose {
		t.Error("Outcome assessment is wrong")
	}
	p1 = []int8{5,6,7}
	p2 = []int8{5,6,7}
	outcome = numberCompare(p1,p2)
	if outcome != outcomes.Tie {
		t.Error("Outcome assessment is wrong")
	}
	p1 = []int8{}
	p2 = []int8{}
	outcome = numberCompare(p1,p2)
	if outcome != outcomes.Tie {
		t.Error("Outcome assessment is wrong")
	}
}

func TestRegisterPlayerHand(t *testing.T) {
	combos := getCombinations()
	//outcomes := getOutcomes()
	winner := 1
	candidate := PlayerCombination{
		CombinationID: combos.Poker,
		Data:          []int8{1},
		Kickers:       []Card{
			{2,'H'},
		},
	}
	existing := PlayerCombination{
		CombinationID: combos.Poker,
		Data:          []int8{3},
		Kickers:       []Card{
			{2,'H'},
		},
	}
	registerPlayerHand(2, candidate, &existing, &winner)
	if winner != 2 {
		t.Error("Player hand registered incorrectly")
	}

	winner = 1
	candidate = PlayerCombination{
		CombinationID: combos.TwoPairs,
		Data:          []int8{1,3},
		Kickers:       []Card{
			{2,'H'},
		},
	}
	existing = PlayerCombination{
		CombinationID: combos.Trips,
		Data:          []int8{2},
		Kickers:       []Card{
			{2,'H'},
		},
	}
	registerPlayerHand(2, candidate, &existing, &winner)
	if winner != 1 {
		t.Error("Player hand registered incorrectly")
	}

	winner = 1
	candidate = PlayerCombination{
		CombinationID: combos.TwoPairs,
		Data:          []int8{6,5},
		Kickers:       []Card{
			{2,'C'},
		},
	}
	existing = PlayerCombination{
		CombinationID: combos.TwoPairs,
		Data:          []int8{6, 5},
		Kickers:       []Card{
			{2,'S'},
		},
	}
	registerPlayerHand(2, candidate, &existing, &winner)
	if winner != -1 {
		t.Error("Player hand registered incorrectly")
	}
}

func TestAddingCardToDeck(t *testing.T) {
	deck := createDeck()
	removeCardFromSlice(&deck, 0)
	fullDeck := createDeck()
	if deck[0] != fullDeck[51] || len(deck) != 51 {
		t.Error("Card was not removed from deck")
	}
	healthy := func() {
		checkDeckHealth(deck)
	}
	assertNoPanic(t, healthy)
}

func TestAddingHandToTable(t *testing.T) {
	hand := Hand{
		[2]Card{
			Card{10, 'C'},
			Card{1, 'D'},
		},
	}
	deck := createDeck()
	var hands []Hand
	addHandToTable(hand, &deck, &hands)
	if hands[0] != hand {
		fmt.Errorf("Hand hasn't been added to the table")
	}
	if len(deck) != 50 {
		fmt.Errorf("Hand hasn't been removed from the deck")
	}
}

func TestGettingRandomCard(t *testing.T) {
	deck := createDeck()
	crds := getRandomCardsFromDeck(&deck, 2)
	if len(deck) != 50 {
		t.Errorf("Did not extract random cards")
	}
	for _, c := range deck {
		if crds[0] == c || crds[1] == c {
			t.Errorf("Did not extract random cards")
		}
	}

	deck2 := createDeck()
	crds = getRandomCardsFromDeck(&deck2, 0)
	if len(deck2) != 52 {
		t.Errorf("Did not extract random cards")
	}
	healthy := func() {
		checkDeckHealth(deck2)
	}
	assertNoPanic(t, healthy)
}

func TestTwoPairCheck(t *testing.T) {
	cards := []Card{
		{5, 'H'},
		{6, 'H'},
		{1, 'D'},
		{1, 'H'},
		{2, 'H'},
		{7, 'S'},
		{2, 'S'},
	}
	expectedKickers := []Card{
		{7, 'S'},
	}

	pairs, kickers := checkTwoPairs(cards)
	expectedPairs := []int8{1,2}
	if !cardSliceContainsSameCards(expectedKickers, kickers) || !EqualInt8Slice(expectedPairs, pairs) {
		t.Errorf("Two Pairs not found")
	}
}

func TestCheckMultipleValues(t *testing.T) {
	cards := []Card{
		{5, 'H'},
		{6, 'H'},
		{7, 'H'},
		{2, 'H'},
		{3, 'H'},
		{1, 'H'},
		{1, 'S'},
	}

	result, kickers := checkMultiples(cards, 2, 3)
	expectedKickers := []Card{
		{5, 'H'},
		{6, 'H'},
		{7, 'H'},
	}

	if result != 1 || len(kickers) != 3 || !cardSliceContainsSameCards(expectedKickers, kickers) {
		t.Errorf("Pair not found")
	}

	cards = []Card{
		{5, 'H'},
		{7, 'H'},
		{7, 'S'},
		{1, 'H'},
		{1, 'S'},
	}
	result, kickers = checkMultiples(cards, 2, 3)
	expectedKickers = []Card{
		{5, 'H'},
		{7, 'H'},
		{7, 'S'},
	}

	if result != 1 || len(kickers) != 3 || !cardSliceContainsSameCards(expectedKickers, kickers) {
		t.Errorf("Pair not found")
	}

	cards = []Card{
		{1, 'H'},
		{7, 'H'},
		{7, 'S'},
		{7, 'C'},
		{1, 'S'},
	}
	result, kickers = checkMultiples(cards, 3, 2)
	expectedKickers = []Card{
		{1, 'H'},
		{1, 'S'},
	}

	if result != 7 || len(kickers) != 2 || !cardSliceContainsSameCards(expectedKickers, kickers) {
		t.Errorf("Pair not found")
	}
}

func TestStraight(t *testing.T) {
	cards := []Card{
		{1, 'H'},
		{2, 'H'},
		{3, 'S'},
		{4, 'C'},
		{5, 'S'},
		{7, 'S'},
		{8, 'S'},
	}
	straight := checkStraight(cards)
	if straight != 5 {
		t.Errorf("Straight not found")
	}

	cards = []Card{
		{1, 'H'},
		{2, 'H'},
		{3, 'S'},
		{4, 'C'},
		{5, 'S'},
		{6, 'S'},
		{7, 'S'},
	}
	straight = checkStraight(cards)
	if straight != 7 {
		t.Errorf("Straight not found")
	}

	cards = []Card{
		{1, 'H'},
		{2, 'H'},
		{9, 'S'},
		{10, 'C'},
		{11, 'S'},
		{12, 'S'},
		{13, 'S'},
	}
	straight = checkStraight(cards)
	if straight != 14 {
		t.Errorf("Straight not found")
	}
}

func TestFlush(t *testing.T) {
	cards := []Card{
		{1, 'S'},
		{2, 'H'},
		{3, 'S'},
		{4, 'C'},
		{5, 'S'},
		{7, 'S'},
		{8, 'S'},
	}
	straight := checkFlush(cards)
	expected := []int8{1,3,5,7,8}
	if !EqualInt8Slice(expected, straight) {
		t.Errorf("Flush not found")
	}

	cards = []Card{
		{1, 'C'},
		{2, 'H'},
		{3, 'S'},
		{4, 'C'},
		{5, 'S'},
		{7, 'S'},
		{8, 'S'},
	}
	straight = checkFlush(cards)
	expected = []int8{}
	if !EqualInt8Slice(expected, straight) {
		t.Errorf("Flush not found")
	}
}

// Asserts that a function throws a panic
func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

// Asserts that a function throws a panic
func assertNoPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()
	f()
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func EqualInt8Slice(a, b []int8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func EqualCardSlice(a, b []Card) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func cardSliceContainsSameCards (a,b []Card) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		found := false
		for _, v2 := range b {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}