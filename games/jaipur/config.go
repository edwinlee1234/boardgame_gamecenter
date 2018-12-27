package jaipur

type gameInit struct {
	Cards       map[string]int32 `json:"cards"`
	CenterCards []string         `json:"centerCards"`
	CardsTotal  int              `json:"CardsTotal"`
}

// 開局 5卡0 camel
// or 3卡2 camel

var cardTotal = 55

var allSellCard = []string{"sliver", "gold", "diamond", "leather", "spice", "cloth"}

var cardsNum = map[string]int32{
	"sliver":  6,
	"gold":    6,
	"diamond": 6,
	"leather": 10,
	"spice":   8,
	"cloth":   8,
	"camel":   11,
}

var cardsPoint = map[string][]int32{
	"sliver":  {5, 5, 5, 5, 5},
	"gold":    {6, 5, 5, 5, 5},
	"diamond": {7, 7, 5, 5, 5},
	"leather": {4, 3, 2, 1, 1, 1, 1, 1, 1, 1},
	"spice":   {5, 3, 3, 2, 2, 1, 1},
	"cloth":   {5, 3, 3, 2, 2, 1, 1},
}

var bonusNum = map[string][]int32{
	"ten_bonus":   {10, 10},
	"six_bonus":   {6, 6},
	"three_bonus": {3, 3},
	"camel_bonus": {5},
}
