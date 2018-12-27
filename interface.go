package main

// IGameClass 遊戲物件的介面
type IGameClass interface {
	Init(userID []int) // 初始化遊戲
	Action()           // 動行
	JudgeWinOrLoss()   // 判斷輸贏
	Broadcast()        // 推播
}
