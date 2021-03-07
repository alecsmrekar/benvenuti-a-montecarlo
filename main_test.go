package main

import (
	"testing"
)

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
	if !EqualCardSlice(getAllNumbers(true), expected) {
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
func EqualCardSlice(a, b []int8) bool {
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
