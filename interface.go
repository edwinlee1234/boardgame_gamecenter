package main

// GameClass 遊戲物件的介面
type GameClass interface {
	Init()           // 初始化遊戲
	Action()         // 動行
	JudgeWinOrLoss() // 判斷輸贏
	Broadcast()      // 推播
}
