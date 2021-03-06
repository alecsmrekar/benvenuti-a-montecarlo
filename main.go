package main

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type Char byte

type CommunityCards struct {
	Cards  []Card
}

// Tells you the status of the game
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

type PlayerCombination struct {
	CombinationID int8
	Data []int8
	Kickers []Card
}

// Formats the players hand
func (combo PlayerCombination) print() string {
	if len(combo.Kickers) == 0 {
		return fmt.Sprintf("%v with cards %v\n", getCombinationType(combo.CombinationID), combo.Data)
	}
	return fmt.Sprintf("%v with cards %v and with kickers %v\n", getCombinationType(combo.CombinationID), combo.Data, combo.Kickers)
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

// Check if [nr] cards with the same value are in the input slice
func checkMultiples(cards []Card, nr int, kickerNr int) (int8, []Card) {
	var kickers []Card
	values, found := findMultipleSameNumbers(cards, nr)
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

	// Return the leftover cards
	kickers = kickers[:kickerNr]
	return max, kickers
}


// Tries to find two pairs
func checkTwoPairs(cards []Card) ([]int8, []Card) {
	twoPairs := []int8{}
	found, kickers := checkMultiples(cards, 2, 3)
	if found == 0 {
		return twoPairs, kickers
	}
	secondFound, kickers := checkMultiples(kickers, 2, 1)
	if secondFound == 0 {
		return twoPairs, kickers
	}
	twoPairs = append(twoPairs, found, secondFound)
	return twoPairs, kickers
}

func checkOnePair(cards []Card) (int8, []Card) {
	result, kickers := checkMultiples(cards, 2, 3)
	return result, kickers
}

func checkTrips(cards []Card) (int8, []Card) {
	result, kickers := checkMultiples(cards, 3, 2)
	return result, kickers
}

func checkPoker(cards []Card) (int8, []Card) {
	result, kickers := checkMultiples(cards, 4, 1)
	return result, kickers
}

func checkStraight(cards []Card) (int8) {
	store := make(map[int8]int8)
	for _, card := range cards {
		store[card.Number]++
		if card.Number == 1 { // An ace also counts as last card
			store[int8(14)]++
		}
	}

	var consecutive, found int8 = 0, 0
	for _, nr := range getAllNumbers(true) {
		if _, ok := store[nr]; ok {
			consecutive++
			if consecutive >=5 {
				found = nr
			}
		} else {
			consecutive = 0
		}
	}
	return found
}

func checkFullHouse(cards []Card) ([]int8) {
	trips, kickers := checkTrips(cards)
	if trips > 0 {
		pair, _ := checkMultiples(kickers, 2, 0)
		if pair > 0 {
			return []int8{trips, pair}
		}
	}
	return []int8{}
}

func checkStraightFlush(cards []Card) int8 {
	highValue := checkStraight(cards)
	if highValue == 0 {
		return 0
	}

	// Keep only cards which are in the straight range
	var keepCards []Card
	for i := 6; i >= 0; i-- {
		card := cards[i]
		if (card.Number >= (highValue-4) && card.Number <= highValue) || (highValue == 14 && card.Number == 1) {
			keepCards = append(keepCards, card)
		}
	}

	foundFlush := checkFlush(keepCards)
	if len(foundFlush) > 1 {
		return highValue
	}
	return 0
}

func checkFlush(cards []Card) ([]int8) {
	store := make(map[Char][]int8)
	for _, card := range cards {
		store[card.Suit] = append(store[card.Suit], card.Number)
	}

	var found []int8
	for i, item := range store {
		if len(item) >= 5 {
			for _, nr := range store[i] {
				found = append(found, nr)
			}
		}
	}

	if len(found) == 0 {
		var emptyResult []int8
		return emptyResult
	}

	// Sort ascending, but take into account the ace
	sort.Slice(found, func(i, j int) bool {
		return (found[i] < found[j] && found[i] > 1)
	})
	// Keep only the 5 highest ones
	found = found[len(found)-5:]
	return found
}

type Outcome struct {
	Win int
	Tie int
	Lose int
}

func getOutcomes() Outcome {
	return Outcome{1,2,3}
}

// Registers a players best hand and determines if it beats the previous best
func registerPlayerHand(id int, candidate PlayerCombination, lastBest *PlayerCombination, winners *int) {
	fmt.Printf("Player %v has: %v", id, candidate.print())

	// If there is not previous hand, this hand wins automatically
	// If this hand has the better combinations, it wins
	if (*lastBest).CombinationID == 0 || candidate.CombinationID < (*lastBest).CombinationID {
		// clear win for the candidate
		*lastBest = candidate
		*winners = id
		return
	} else if candidate.CombinationID > (*lastBest).CombinationID {
		// Loss for the candidate
		return
	}

	// From here down, the previous best and the candidate have the best combination
	// We need to compare in more detail

	outcomes := getOutcomes()

	greaterEqualOrLower := func(c1 int8, c2 int8) int {
		if c1 == c2 {
			return outcomes.Tie // equal
		} else if c1 == 1 {
			return outcomes.Win // greater
		} else if c2 == 1 {
			return outcomes.Lose // greater
		} else if c1 > c2 {
			return outcomes.Win // greater
		}
		return outcomes.Lose	// lower
	}

	numberCompare := func(k1, k2 []int8) int {
		if len(k1) != len(k2) {
			panic("Kicker length should be the same")
		}
		for i, _ := range k1 {
			result := greaterEqualOrLower(k1[i], k2[i])
			if result != outcomes.Tie {
				return result
			}
		}
		return outcomes.Tie
	}

	kickerCompare := func(k1, k2 []Card) int {
		var k1n, k2n []int8
		for _, v := range k1 {
			k1n	= append(k1n, v.Number)
		}
		for _, v := range k2 {
			k2n	= append(k2n, v.Number)
		}
		return numberCompare(k1n, k2n)
	}

	var outcome int
	switch candidate.CombinationID {
	case 1, 3, 4, 5:	// straight flush, normal straight, full house or flush
		outcome = numberCompare(candidate.Data, (*lastBest).Data)
	case 2, 6, 7, 8: // poker,  trips, two pair, on pair
		outcome = numberCompare(candidate.Data, (*lastBest).Data)
		if outcome == outcomes.Tie {
			outcome = kickerCompare(candidate.Kickers, (*lastBest).Kickers)
		}
	case 9: // high card
		outcome = kickerCompare(candidate.Kickers, (*lastBest).Kickers)
	default:
		outcome = outcomes.Tie
	}

	if outcome == outcomes.Win {
		// we have a clear winner
		*winners = id
		*lastBest = candidate
	} else if outcome == outcomes.Tie {
		// Of there is a tie, we don't have a current single winner
		*winners = -1
	} else if outcome == 0 {
		panic("Outcome hasn't been asserted")
	}
}

// Retrieves scenarios from the job queue and crunches them
func casinoWorker(results chan<- int, jobs <-chan Game) {
	fmt.Println("Starting worker")

	// Retrieve a single job (= one game)
	for work := range jobs {
		communityCards := work.Table.Cards
		tableStatus := work.Table.status()
		mapping := getStatusMap()
		deck := work.Deck
		cardsLeftToPull := mapping[tableStatus]
		cardsPulled := getRandomCardsFromDeck(&deck, cardsLeftToPull)
		communityCards = append(communityCards, cardsPulled...)
		lastBest := PlayerCombination{}
		var weHaveAWinner int = -1

		// Calculate the best combination each player holds
		for playerIndex, hand := range work.Hands {
			var playerCardPool []Card = communityCards
			playerCardPool = append(playerCardPool, hand.Cards[:]...)
			var foundInt int8
			var foundSlice []int8
			var kickers []Card
			if len(playerCardPool) != 7 {
				panic("Player should have 7 cards available in total")
			}
			checkDeckHealth(append(deck, communityCards...))

			// The best hand rank returns the lower value
			foundInt = checkStraightFlush(playerCardPool)
			if foundInt > 0{
				combo := PlayerCombination{1, []int8{foundInt}, kickers}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundInt, kickers = checkPoker(playerCardPool)
			if foundInt > 0{
				combo := PlayerCombination{2, []int8{foundInt}, kickers}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundSlice = checkFullHouse(playerCardPool)
			if len(foundSlice) > 0{
				combo := PlayerCombination{3, foundSlice, []Card{}}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundSlice = checkFlush(playerCardPool)
			if len(foundSlice) > 0{
				combo := PlayerCombination{4, foundSlice, []Card{}}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundInt = checkStraight(playerCardPool)
			if foundInt > 0{
				combo := PlayerCombination{5, []int8{foundInt}, []Card{}}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundInt, kickers = checkTrips(playerCardPool)
			if foundInt > 0{
				combo := PlayerCombination{6, []int8{foundInt}, kickers}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundSlice, kickers = checkTwoPairs(playerCardPool)
			if len(foundSlice) == 2{
				combo := PlayerCombination{7, foundSlice, kickers}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			foundInt, kickers = checkOnePair(playerCardPool)
			if foundInt > 0{
				combo := PlayerCombination{8, []int8{foundInt}, kickers}
				registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
				continue
			}
			sort.Sort(sort.Reverse(ByNumber(playerCardPool)))
			playerCardPool = playerCardPool[:5]
			combo := PlayerCombination{9, []int8{}, playerCardPool}
			registerPlayerHand(playerIndex, combo, &lastBest, &weHaveAWinner)
		}

		if weHaveAWinner >= 0 {
			fmt.Printf("Player %v wins\n\n", weHaveAWinner)
		} else {
			fmt.Println("No winner")
		}
		results <- weHaveAWinner
	}
	fmt.Println("Worker done")
}

func main() {

	deck := createDeck()
	var hands []Hand

	h1 := Hand{
		[2]Card{
			Card{1, 'H'},
			Card{2, 'H'},
		},
	}
	h2 := Hand{
		[2]Card{
			Card{1, 'C'},
			Card{2, 'C'},
		},
	}
	h3 := Hand{
		[2]Card{
			Card{10, 'C'},
			Card{1, 'D'},
		},
	}

	addHandToTable(h1, &deck, &hands)
	addHandToTable(h2, &deck, &hands)
	addHandToTable(h3, &deck, &hands)

	var insertCards = []Card{
		//Card{6, 'H'},
	}
	table := CommunityCards{
		insertCards,
	}

	// Remove community cards from deck
	for _, c := range table.Cards {
		addCardToTable(c, &deck)
	}

	var workers, simulations int
	for ;workers==0; {
		fmt.Println("Number of threads to use: ")
		fmt.Scanf("%d", &workers)
	}

	for ;simulations==0; {
		fmt.Println("Number of simulated games to run: ")
		fmt.Scanf("%d", &simulations)
	}

	start := time.Now()
	resultsChannel := make(chan int, simulations)
	jobsChannel := make(chan Game, simulations)

	for i := 0; i < workers; i++ {
		go casinoWorker(resultsChannel, jobsChannel)
	}
	for i := 0; i < simulations; i++ {
		// Make a new deck slice for each worker
		deckDestination := make([]Card, len(deck))
		copy(deckDestination, deck)

		setting := Game{
			Table: table,
			Hands: hands,
			Deck: deckDestination,
		}
		jobsChannel <- setting
	}

	close(jobsChannel)
	results := make(map[int]int)
	winCount := 0

	for i := 0; i < simulations; i++ {
		winner := <- resultsChannel
		if winner >= 0 {
			results[winner]++
			winCount++
		}
	}
	fmt.Println("-------")
	simulationsF := float64(simulations)

	for i, wins := range results {
		winProbability := float64(wins)/simulationsF*100
		fmt.Printf("Player ID %v win probability: %f%% \n", i, winProbability)
	}

	splitProbability := (simulationsF-float64(winCount))/simulationsF*100
	fmt.Printf("Split probability: %f%% \n\n", splitProbability)

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}