package jaipur

import "fmt"

// NewJaipur new Jaipur物件
func NewJaipur() *Jaipur {
	return &Jaipur{
		DeskCard:       []string{},
		unfoldCard:     []string{},
		player1Card:    []string{},
		player2Card:    []string{},
		remainPoinCard: make(map[string][]string),
	}
}

// Jaipur 遊戲的物件
type Jaipur struct {
	DeskCard       []string
	unfoldCard     []string
	player1Card    []string
	player2Card    []string
	remainPoinCard map[string][]string
}

// Init 初始化
func (j *Jaipur) Init() {
	fmt.Println("init")
}

// Action 行為
func (j *Jaipur) Action() {

}

// JudgeWinOrLoss 判斷輸贏
func (j *Jaipur) JudgeWinOrLoss() {

}

// Broadcast 推播
func (j *Jaipur) Broadcast() {

}
