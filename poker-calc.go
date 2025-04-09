package main

import (
    "encoding/json"
    "fmt"
    "os"
    "math/rand"
    "time"
    "sort"
    "strconv"
    "math"
)

type TableData struct {
    CommunityCards []string `json:"community_cards"`
    HandCards      []string `json:"hand_cards"`
    PlayerCount    int      `json:"player_count"`
}

var faceToRank = map[string]int{
    "2": 2, "3": 3, "4": 4, "5": 5,
    "6": 6, "7": 7, "8": 8, "9": 9,
    "T": 10, "J": 11, "Q": 12, "K": 13, "A": 14,
}

func getDeck() []string {
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
    suits := []string{"s", "c", "d", "h"}
    var deck []string

    for _, suit := range suits {
        for _, value := range values {
            deck = append(deck, value+suit)
        }
    }

	return deck
}

func removeFromList(slice []string, value string) []string {

    new_list := []string{}

    for _, item := range slice {
        if item != value {
            new_list = append(new_list, item)
        }
    }

    return new_list
}

func containsStr(slice []string, value string) bool {
    for _, item := range slice {
        if item == value {
            return true
        }
    }
    return false
}

func containsInt(slice []int, value int) bool {
    for _, item := range slice {
        if item == value {
            return true
        }
    }
    return false
}

func max(nums []int) int {
    m := nums[0]
    for _, n := range nums {
        if n > m { m = n }
    }
    return m
}

func getMaxCardIndex(cards []string) int {
    card_ranks := []int{}
    for _, card := range cards {
        card_ranks = append(card_ranks, faceToRank[string(card[0])])
    }
    max_rank := max(card_ranks)
    max_rank_index := getIntIndex(card_ranks, max_rank)

    return max_rank_index
}

func removeByValueInt(slice []int, value int) []int {
    var result []int
    for _, v := range slice {
        if v != value {
            result = append(result, v)
        }
    }
    return result
}

func removeByValueString(slice []string, value string) []string {
    var result []string
    for _, v := range slice {
        if v != value {
            result = append(result, v)
        }
    }
    return result
}

func anyInListString(first []string, second []string) bool {
    for _, item1 := range first {
        for _, item2 := range second {
            if item1 == item2 {
                return true
            }
        }
    }
    return false
}

func getIntIndex(slice []int, value int) int {
    for i, item := range slice {
        if item == value {
            return i
        }
    }
    return -1 // return -1 if not found
}

func findCardByRank(cards []string, rank int) (int) {
    card_ranks := []int{}
    for _, card := range cards {
        card_ranks = append(card_ranks, faceToRank[string(card[0])])
    }

    for i, card_rank := range card_ranks {
        if card_rank == rank {
            return i
        }
    }

    return -1
}

func addToUsedCards(available_cards []string, used_cards []string, used_card string) ([]string, []string) {
    available_cards = removeByValueString(available_cards, used_card)
    used_cards = append(used_cards, used_card)

    return available_cards, used_cards
}

func boolToInt(b bool) int {
    if b {
        return 1 // return 1 for true
    }
    return 0 // return 0 for false
}

func sortMapByKeys(m map[int]int) map[int]int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	sortedMap := make(map[int]int)
	for _, k := range keys {
		sortedMap[k] = m[k]
	}

	return sortedMap
}

func roundToTwoDecimals(num float64) float64 {
	return math.Round(num*100) / 100
}

func evaluateHand(hand_cards []string, community_cards []string) [][]int {
    seven_cards := append(hand_cards, community_cards...)

    //fmt.Println(seven_cards)

    card_faces := []string{}
    card_colors := []string{}
    card_ranks := []int{}
    for _, card := range seven_cards {
        card_faces = append(card_faces, string(card[0]))
        card_colors = append(card_colors, string(card[1]))
    }

    colors_count := map[string]int{}
    for _, color := range card_colors {
        colors_count[color]++
    }

    rank_count := map[int]int{}
    for _, face := range card_faces {
        rank := faceToRank[string(face)]
        rank_count[rank]++
        card_ranks = append(card_ranks, rank)
    }
    rank_count = sortMapByKeys(rank_count)

    available_kickers := card_ranks
    available_cards := seven_cards
    used_cards := []string{}


    is_flush := false
    flush_high_card := -1
    flush_cards := []string{}
    //flush_color := ""
    for color, count := range colors_count {
        if count >= 5 {
            is_flush = true
            //flush_color = color

            for index, cardColor := range card_colors {
                if (cardColor == color) {
                    flush_cards = append(flush_cards, seven_cards[index])
                    if (card_ranks[index] > flush_high_card) {
                        flush_high_card = card_ranks[index]
                    }
                }
            }
        }
    }

    is_straight := false
    is_royal := false
    straight_high_card := -1
    straight_cards := []string{}
    for _, rank := range card_ranks {
        if (containsInt(card_ranks, rank+1) && containsInt(card_ranks, rank+2) && containsInt(card_ranks, rank+3) && containsInt(card_ranks, rank+4)) {
            is_straight = true
            if (rank+4 > straight_high_card) {
                straight_high_card = rank+4
                straight_cards = append(straight_cards, seven_cards[findCardByRank(seven_cards, rank)])
                straight_cards = append(straight_cards, seven_cards[findCardByRank(seven_cards, rank+1)])
                straight_cards = append(straight_cards, seven_cards[findCardByRank(seven_cards, rank+2)])
                straight_cards = append(straight_cards, seven_cards[findCardByRank(seven_cards, rank+3)])
                straight_cards = append(straight_cards, seven_cards[findCardByRank(seven_cards, rank+4)])
            }

            if (straight_high_card == 14) {
                is_royal = true
            }
        }
    } 

    pairs := []int{}
    three_kind := []int{}
    four_kind := []int{}

    pairs_cards := [][]string{}
    three_kind_cards := [][]string{}
    four_kind_cards := [][]string{}

    for rank, count := range rank_count {
        if (count == 2) {
            pairs = append(pairs, rank)

            card1 := seven_cards[findCardByRank(seven_cards, rank)]
            temp_list1 := removeByValueString(seven_cards, card1)
            card2 := seven_cards[findCardByRank(temp_list1, rank)]

            pairs_cards = append(pairs_cards, []string{card1, card2})
        } else if (count == 3) {
            three_kind = append(three_kind, rank)

            card1 := seven_cards[findCardByRank(seven_cards, rank)]
            temp_list1 := removeByValueString(seven_cards, card1)
            card2 := seven_cards[findCardByRank(temp_list1, rank)]
            temp_list2 := removeByValueString(temp_list1, card1)
            card3 := seven_cards[findCardByRank(temp_list2, rank)]

            three_kind_cards = append(three_kind_cards, []string{card1, card2, card3})
        } else if (count == 4) {
            four_kind = append(four_kind, rank)

            card1 := seven_cards[findCardByRank(seven_cards, rank)]
            temp_list1 := removeByValueString(seven_cards, card1)
            card2 := seven_cards[findCardByRank(temp_list1, rank)]
            temp_list2 := removeByValueString(temp_list1, card1)
            card3 := seven_cards[findCardByRank(temp_list2, rank)]
            temp_list3 := removeByValueString(temp_list2, card1)
            card4 := seven_cards[findCardByRank(temp_list3, rank)]

            four_kind_cards = append(four_kind_cards, []string{card1, card2, card3, card4})
        }
    }

    combos_list := [][]int{
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
    }

    //Evaluating final combos

    //Royal/Straight flush
    if (is_flush && is_straight) { 
        available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[0])
        available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[1])
        available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[2])
        available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[3])
        available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[4])

        if (is_royal) {
            combos_list[0] = append(combos_list[0], straight_high_card)
        } else {
            combos_list[1] = append(combos_list[1], straight_high_card)
        }
    }

    //Four of a kind 
    if (len(four_kind) > 0) {
        if (!anyInListString(four_kind_cards[0], used_cards)) {
            combos_list[2] = four_kind
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, four_kind_cards[0][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, four_kind_cards[0][1])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, four_kind_cards[0][2])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, four_kind_cards[0][3])
        }
    }

    //Full house
    if (len(three_kind) > 0 && len(pairs) > 0) {
        if (!anyInListString(three_kind_cards[0], used_cards) && !anyInListString(pairs_cards[0], used_cards)) {
            max_threekind := max(three_kind)
            max_pair := max(pairs)

            three_kind = removeByValueInt(three_kind, max_threekind)
            pairs = removeByValueInt(pairs, max_pair)

            available_kickers = removeByValueInt(available_kickers, max_threekind)
            available_kickers = removeByValueInt(available_kickers, max_pair)

            available_cards, used_cards = addToUsedCards(available_cards, used_cards, three_kind_cards[0][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, three_kind_cards[0][1])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, three_kind_cards[0][2])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[0][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[0][1])

            combos_list[3] = append(combos_list[3], max_threekind)
        }
    }

    // Flush
    if (is_flush && !is_straight) {
        if (!anyInListString(flush_cards, used_cards)) {
            combos_list[4] = append(combos_list[4], flush_high_card)
            
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[1])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[2])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[3])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, flush_cards[4])
        }
    }

    // Straight
    if (is_straight && !is_flush) {
        if (!anyInListString(straight_cards, used_cards)) {
            combos_list[5] = append(combos_list[5], straight_high_card)

            available_cards, used_cards = addToUsedCards(available_cards, used_cards, straight_cards[0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, straight_cards[1])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, straight_cards[2])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, straight_cards[3])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, straight_cards[4])
        }
    }

    //Three of a kind 
    if (len(three_kind) > 0) {
        if (!anyInListString(three_kind_cards[0], used_cards)) {
            combos_list[6] = three_kind

            available_cards, used_cards = addToUsedCards(available_cards, used_cards, three_kind_cards[0][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, three_kind_cards[0][1])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, three_kind_cards[0][2])
        }
    }

    //Two pair
    if (len(pairs) > 1) {
        if (!anyInListString(pairs_cards[0], used_cards) && !anyInListString(pairs_cards[1], used_cards)) {
            max_pair1 := max(pairs)
            pairs = removeByValueInt(pairs, max_pair1)
            max_pair2 := max(pairs)
            pairs = removeByValueInt(pairs, max_pair2)

            combos_list[7] = append(combos_list[7], max_pair1)
            combos_list[7] = append(combos_list[7], max_pair2)

            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[0][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[0][1])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[1][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[1][1])
        }
    }

    //Pair
    if (len(pairs) > 0) {
        if (!anyInListString(pairs_cards[0], used_cards)) {
            combos_list[8] = pairs

            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[0][0])
            available_cards, used_cards = addToUsedCards(available_cards, used_cards, pairs_cards[0][1])
        }
    }

    //High card/Kicker
    if (len(available_cards) > 0) {
        kicker := available_cards[getMaxCardIndex(available_cards)]
        combos_list[9] = append(combos_list[9], faceToRank[string(kicker[0])])
    } else {
        combos_list[9] = append(combos_list[9], 0)
    }

    return combos_list
}

func compareHands(hand1 [][]int, hand2 [][]int) int {
    // -1: Hand 1 wins
    // 0: Tie
    // 1: Hand 2 wins

    // Royal Flush
    if (len(hand1[0]) > 0 || len(hand2[0]) > 0) {
        if (len(hand2[0]) == 0){
            return -1
        } else {
            return 1
        }
    }

    // Straight flush (hand[1])
    if (len(hand1[1]) > 0 || len(hand2[1]) > 0) {
        if (len(hand1[1]) > 0 && len(hand2[1]) == 0) {
            return -1
        } else if (len(hand1[1]) == 0 && len(hand2[1]) > 0) {
            return 1
        } else if (len(hand1[1]) > 0 && len(hand2[1]) > 0) {
            if (hand1[1][0] > hand2[1][0]) {
                return -1
            } else if (hand1[1][0] < hand2[1][0]) {
                return 1
            } else if (hand1[1][0] == hand2[1][0]) {
                //continue
            }
        }
    }

    //Four of a kind
    if (len(hand1[2]) > 0 || len(hand2[2]) > 0) {
        if (len(hand1[2]) > 0 && len(hand2[2]) == 0) {
            return -1
        } else if (len(hand1[2]) == 0 && len(hand2[2]) > 0) {
            return 1
        } else if (len(hand1[2]) > 0 && len(hand2[2]) > 0) {
            if (hand1[2][0] > hand2[2][0]) {
                return -1
            } else if (hand1[2][0] < hand2[2][0]) {
                return 1
            } else if (hand1[2][0] == hand2[2][0]) {
                //continue
            }
        }
    }

    //Full house
    if (len(hand1[3]) > 0 || len(hand2[3]) > 0) {
        if (len(hand1[3]) > 0 && len(hand2[3]) == 0) {
            return -1
        } else if (len(hand1[3]) == 0 && len(hand2[3]) > 0) {
            return 1
        } else if (len(hand1[3]) > 0 && len(hand2[3]) > 0) {
            if (hand1[3][0] > hand2[3][0]) {
                return -1
            } else if (hand1[3][0] < hand2[3][0]) {
                return 1
            } else if (hand1[3][0] == hand2[3][0]) {
                //continue
            }
        }
    }

    //Flush
    if (len(hand1[4]) > 0 || len(hand2[4]) > 0) {
        if (len(hand1[4]) > 0 && len(hand2[4]) == 0) {
            return -1
        } else if (len(hand1[4]) == 0 && len(hand2[4]) > 0) {
            return 1
        } else if (len(hand1[4]) > 0 && len(hand2[4]) > 0) {
            if (hand1[4][0] > hand2[4][0]) {
                return -1
            } else if (hand1[4][0] < hand2[4][0]) {
                return 1
            } else if (hand1[4][0] == hand2[4][0]) {
                //continue
            }
        }
    }

    //Straight
    if (len(hand1[5]) > 0 || len(hand2[5]) > 0) {
        if (len(hand1[5]) > 0 && len(hand2[5]) == 0) {
            return -1
        } else if (len(hand1[5]) == 0 && len(hand2[5]) > 0) {
            return 1
        } else if (len(hand1[5]) > 0 && len(hand2[5]) > 0) {
            if (hand1[5][0] > hand2[5][0]) {
                return -1
            } else if (hand1[5][0] < hand2[5][0]) {
                return 1
            } else if (hand1[5][0] == hand2[5][0]) {
                //continue
            }
        }
    }

    //Three of a kind
    if (len(hand1[6]) > 0 || len(hand2[6]) > 0) {
        if (len(hand1[6]) > 0 && len(hand2[6]) == 0) {
            return -1
        } else if (len(hand1[6]) == 0 && len(hand2[6]) > 0) {
            return 1
        } else if (len(hand1[6]) > 0 && len(hand2[6]) > 0) {
            if (hand1[6][0] > hand2[6][0]) {
                return -1
            } else if (hand1[6][0] < hand2[6][0]) {
                return 1
            } else if (hand1[6][0] == hand2[6][0]) {
                //continue
            }
        }
    }

    //Two pair
    if (len(hand1[7]) > 0 || len(hand2[7]) > 0) {
        if (len(hand1[7]) > 0 && len(hand2[7]) == 0) {
            return -1
        } else if (len(hand1[7]) == 0 && len(hand2[7]) > 0) {
            return 1
        } else if (len(hand1[7]) > 0 && len(hand2[7]) > 0) {
            if (hand1[7][0] > hand2[7][0]) {
                return -1
            } else if (hand1[7][0] < hand2[7][0]) {
                return 1
            } else if (hand1[7][0] == hand2[7][0]) {
                if (hand1[7][1] > hand2[7][1]) {
                    return -1
                } else if (hand1[7][1] < hand2[7][1]) {
                    return 1
                } else if (hand1[7][1] == hand2[7][1]) {
                    //continue
                }
            }
        }
    }

    //Pair
    if (len(hand1[8]) > 0 || len(hand2[8]) > 0) {
        if (len(hand1[8]) > 0 && len(hand2[8]) == 0) {
            return -1
        } else if (len(hand1[8]) == 0 && len(hand2[8]) > 0) {
            return 1
        } else if (len(hand1[8]) > 0 && len(hand2[8]) > 0) {
            if (hand1[8][0] > hand2[8][0]) {
                return -1
            } else if (hand1[8][0] < hand2[8][0]) {
                return 1
            } else if (hand1[8][0] == hand2[8][0]) {
                //continue
            }
        }
    }

    //High card/Kicker
    if (hand1[9][0] > hand2[9][0]) {
        return -1
    } else if (hand1[9][0] < hand2[9][0]) {
        return 1
    } else if (hand1[9][0] == hand2[9][0]) {
        return 0
    }

    return 0
}

func getRemoveRandomCard(floatingDeck []string) (string, []string) {
    rand.Seed(time.Now().UnixNano())
    //rand.Seed(2)
    randomIndex := rand.Intn(len(floatingDeck))
    randomCard := floatingDeck[randomIndex]

    floatingDeck = removeFromList(floatingDeck, randomCard)

    return randomCard, floatingDeck
}

func findWinner(evaluated_hands [][][]int) []int {
    current_winner := []int{0}
    for i := 0; i < len(evaluated_hands)-1; i++ {
        result := compareHands(evaluated_hands[current_winner[0]], evaluated_hands[i+1])
        if (result == -1) {
            //Current winner stays
        } else if (result == 1) {
            current_winner = []int{i+1}
        } else if (result == 0) {
            current_winner = append(current_winner, i+1)
        }
    }

    return current_winner
}

func main() {
    file, err := os.Open("input_state.json")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    var tableData TableData
    if err := json.NewDecoder(file).Decode(&tableData); err != nil {
        fmt.Println(err)
        return
    }

    //TODO Verify for duplicate cards
    //TODO add cli arguments for input file, iterations and optional player data output

    //fmt.Println(tableData)

    max_iterations := 20000
    start := time.Now()

    player_stats := [][]int{}
    for i := 0; i < tableData.PlayerCount; i++ {
        player_stats = append(player_stats, []int{0,0,0,0,0,0,0,0,0,0,0})
    }

    //Simulating games
    for i := 0; i < max_iterations; i++ {
        floatingDeck := getDeck()

        //Remove hand cards
        for _, value := range tableData.HandCards {
            floatingDeck = removeFromList(floatingDeck, value)
        }

        //Remove community cards
        for _, value := range tableData.CommunityCards {
            floatingDeck = removeFromList(floatingDeck, value)
        }
        
        table_cards := make([][]string, 2)


        //Filling community cards, if needed
        filled_community_cards := []string{}
        for _, card := range tableData.CommunityCards {
            if (card == "not_drawn") {
                new_card := ""
                new_card, floatingDeck = getRemoveRandomCard(floatingDeck)
                filled_community_cards = append(filled_community_cards, new_card)
            } else {
                filled_community_cards = append(filled_community_cards, card)
            }
        } 
        table_cards[0] = filled_community_cards
        //fmt.Println(filled_community_cards)

        //Adding own hand
        table_cards[1] = tableData.HandCards

        //Filling player cards
        for i := 0; i < tableData.PlayerCount; i++ {
            card1, card2 := "", ""
            card1, floatingDeck = getRemoveRandomCard(floatingDeck)
            card2, floatingDeck = getRemoveRandomCard(floatingDeck)
            player_hand := []string{card1, card2}
            table_cards = append(table_cards, player_hand)
        }

        //Evaluating all hands
        evaluated_hands := [][][]int{}
        for i := 0; i < tableData.PlayerCount; i++ {
            evaluated_hand := evaluateHand(table_cards[0], table_cards[i+1])
            //fmt.Println(table_cards[i+1])
            //fmt.Println(evaluated_hand)
            evaluated_hands = append(evaluated_hands, evaluated_hand)
        }

        //fmt.Println(table_cards)
        //fmt.Println(evaluated_hands)

        //Find winner
        winning_hand := findWinner(evaluated_hands)
        //fmt.Println(winning_hand)
        //fmt.Println(i)

        //Log stats
        for i := 0; i < tableData.PlayerCount; i++ {
            player_stats[i][0] += boolToInt(len(evaluated_hands[i][0]) > 0)
            player_stats[i][1] += boolToInt(len(evaluated_hands[i][1]) > 0)
            player_stats[i][2] += boolToInt(len(evaluated_hands[i][2]) > 0)
            player_stats[i][3] += boolToInt(len(evaluated_hands[i][3]) > 0)
            player_stats[i][4] += boolToInt(len(evaluated_hands[i][4]) > 0)
            player_stats[i][5] += boolToInt(len(evaluated_hands[i][5]) > 0)
            player_stats[i][6] += boolToInt(len(evaluated_hands[i][6]) > 0)
            player_stats[i][7] += boolToInt(len(evaluated_hands[i][7]) > 0)
            player_stats[i][8] += boolToInt(len(evaluated_hands[i][8]) > 0)
            player_stats[i][9] += boolToInt(len(evaluated_hands[i][9]) > 0)
            player_stats[i][10] += boolToInt(getIntIndex(winning_hand, i) > -1)
        }

        //fmt.Println(evaluated_hands)

    }

    //fmt.Println(player_stats)

    duration := time.Since(start).Seconds()

    for i := 0; i < tableData.PlayerCount; i++ {
        fmt.Println("Player " + strconv.Itoa(i+1))
        fmt.Println("-------------------------")
        fmt.Println("Win: " + strconv.FormatFloat(float64(float64(player_stats[i][10])/float64(max_iterations)*100), 'f', 2, 64) + "%")
        fmt.Println()
        fmt.Println()
    }
    fmt.Printf("Execution time: %.2f seconds\n", duration)

}