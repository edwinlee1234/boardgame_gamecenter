package main

// 遊戲狀態
var notOpen = 0 // owner only
var opening = 1 // 開放玩家
var playing = 2 // 遊戲中

const (
	JAIPUR = "jaipur"
	LOBBU  = "lobby"
)

var channelSupport = map[string]bool{
	JAIPUR: true,
	LOBBU:  true,
}

// 支援的遊戲
var gameSupport = map[string]bool{
	JAIPUR: true,
}
