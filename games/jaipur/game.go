package jaipur

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// NewJaipur new Jaipur物件
func NewJaipur(gameID int32, hub *JaipurHub) *Jaipur {
	return &Jaipur{
		gameID:       gameID,
		hub:          hub,
		deskCard:     []string{},
		foldCard:     []string{},
		players:      Players{},
		cardsPoint:   make(map[string][]int32),
		action:       []string{},
		actionPlayer: 0,
	}
}

// Players 玩家
type Players []PlayerInfo

// Jaipur 遊戲的物件
type Jaipur struct {
	gameID     int32
	hub        *JaipurHub
	deskCard   []string
	foldCard   []string
	players    Players
	cardsPoint map[string][]int32
	// 行動的記錄
	action []string
	// 可以行動的玩家
	actionPlayer int32
}

// PlayerInfo 玩家資訊
type PlayerInfo struct {
	// 會員ID
	ID    int32    `json:"id"`
	UUID  string   `json:"uuid"`
	Card  []string `json:"card"`
	Camel int32    `json:"camel"`
	Point int32    `json:"point"`
	Bouns []int32  `json:"bouns"`
}

// OpponentPlayer info
type OpponentPlayer struct {
	CardNum   int32 `json:"cardNum"`
	HaveCamel bool  `json:"haveCamel"`
}

// Action 行動的struct
type Action struct {
	Type             string
	Take             int32
	Sell             []int32
	SwitchSelfCard   []int32
	SwitchTargetCard []int32
}

// TakeEvent 拿牌Event
type TakeEvent struct {
	Event string        `json:"event"`
	Info  TakeEventInfo `json:"info"`
}

// TakeEventInfo 拿牌Event Info
type TakeEventInfo struct {
	DeskCardNum int32 `json:"deskCardNum"`
	PlayerID    int32 `json:"playerID"`
}

// TakeEventPrivate 拿牌Event user (動作的user專用)
type TakeEventPrivate struct {
	Event string               `json:"event"`
	Info  TakeEventInfoPrivate `json:"info"`
}

// TakeEventInfoPrivate 拿牌Event Info (動作的user專用)
type TakeEventInfoPrivate struct {
	DeskCardNum int32
	Card        []string `json:"card"`
}

// SellEvent 拿牌Event
type SellEvent struct {
	Event string        `json:"event"`
	Info  TakeEventInfo `json:"info"`
}

// SellEventInfo 拿牌Event Info
type SellEventInfo struct {
	DeskCardNum int32              `json:"deskCardNum"`
	NewCard     string             `json:"newCard"`
	FoldCardNum int32              `json:"foldCardNum"`
	CardsPoint  map[string][]int32 `json:"cardsPoint"`
}

// SellEventPrivate 賣牌 (動作的user專用)
type SellEventPrivate struct {
	Event string               `json:"event"`
	Info  SellEventInfoPrivate `json:"info"`
}

// SellEventInfoPrivate 賣牌 (動作的user專用)
type SellEventInfoPrivate struct {
	Bonus int32    `json:"bonus"`
	Card  []string `json:"card"`
}

// InfoEvent InfoEvent
type InfoEvent struct {
	Event string `json:"event"`
	Info  Info   `json:"info"`
}

// Info game Info
type Info struct {
	DeskCard     []string           `json:"deskCard"`
	CardsPoint   map[string][]int32 `json:"cardsPoint"`
	FoldCardNum  int32              `json:"foldCardNum"`
	ActionPlayer int32              `json:"actionPlayer"`
	Players      Players            // 之後拿掉
}

// UserInfoEvent UserInfoEvent
type UserInfoEvent struct {
	Event    string         `json:"event"`
	UserInfo PlayerInfo     `json:"player_info"`
	Opponent OpponentPlayer `json:"opponentPlayer"`
}

// Info 遊戲的資料，包括某玩家的牌
func (j *Jaipur) Info() Info {
	var info Info
	info.DeskCard = j.deskCard
	info.CardsPoint = j.cardsPoint
	info.FoldCardNum = int32(len(j.foldCard))
	info.ActionPlayer = j.players[j.actionPlayer].ID
	info.Players = j.players // 之後拿掉

	return info
}

// BroadcastInfo 推info
func (j *Jaipur) BroadcastInfo(userID []int32) {
	var infoEvent InfoEvent
	gameInfo := j.Info()
	infoEvent.Event = "Info"
	infoEvent.Info = gameInfo

	var player1Event UserInfoEvent
	player1Event.Event = "UserInfo"
	player1Event.UserInfo = j.players[0]
	// 對手的資訊
	player1Event.Opponent = OpponentPlayer{
		CardNum:   int32(len(j.players[1].Card)),
		HaveCamel: (j.players[1].Camel != 0),
	}

	var player2Event UserInfoEvent
	player2Event.Event = "UserInfo"
	player2Event.UserInfo = j.players[1]
	// 對手的資訊
	player2Event.Opponent = OpponentPlayer{
		CardNum:   int32(len(j.players[0].Card)),
		HaveCamel: (j.players[0].Camel != 0),
	}

	infoJSON, _ := json.Marshal(infoEvent)
	player1JSON, _ := json.Marshal(player1Event)
	player2JSON, _ := json.Marshal(player2Event)

	// 沒有指定某一個會員就全推
	if userID == nil {
		j.hub.BroadcastChannel(j.gameID, infoJSON)
		j.hub.BroadcastUser(j.gameID, j.players[0].UUID, player1JSON)
		j.hub.BroadcastUser(j.gameID, j.players[1].UUID, player2JSON)

		return
	}

	// 指定某會員，用於某會員重新刷新畫面的init，用推播推
	for _, ID := range userID {
		if ID == j.players[0].ID {
			j.hub.BroadcastUser(j.gameID, j.players[0].UUID, player1JSON)
			j.hub.BroadcastUser(j.gameID, j.players[0].UUID, infoJSON)
		}

		if ID == j.players[1].ID {
			j.hub.BroadcastUser(j.gameID, j.players[1].UUID, player2JSON)
			j.hub.BroadcastUser(j.gameID, j.players[1].UUID, infoJSON)
		}
	}
}

// Init 初始化
func (j *Jaipur) Init(usersInfo map[int32]string) {
	rand.Seed(time.Now().UnixNano()) // 這一行一定要加，初始化隨機數

	cardType := []string{"sliver", "gold", "diamond", "leather", "spice", "cloth", "camel"}
	var totalCard []string
	j.cardsPoint = cardsPoint
	nowCardsNum := cardsNum

	// 兩方玩家發牌
	for ID, UUID := range usersInfo {
		var playerCard []string
		var playerStartCard string
		var camelNum int32

		// 決定開局手牌 5：0 還是 3：2 （兩個玩家可以不一樣）
		if rand.Intn(2) == 0 {
			playerStartCard = "no-camel"
		} else {
			playerStartCard = "have-camel"
		}

		if playerStartCard == "no-camel" {
			for i := 0; i < 5; i++ {
				num := rand.Intn(6)
				card := cardType[num]
				nowCardsNum[card]--

				playerCard = append(playerCard, card)
			}
		} else {
			for i := 0; i < 3; i++ {
				num := rand.Intn(6)
				card := cardType[num]
				nowCardsNum[card]--

				playerCard = append(playerCard, card)
			}

			nowCardsNum["camel"] -= 2
			camelNum = 2
		}

		// 看誰先行動
		if rand.Intn(2) == 0 {
			j.actionPlayer = 0
		} else {
			j.actionPlayer = 1
		}

		j.players = append(j.players, PlayerInfo{
			ID,
			UUID,
			playerCard,
			camelNum,
			0,
			[]int32{},
		})

		cardTotal -= 5
	}

	// 全部的牌random一次
	i := 0
	for i < cardTotal {
		num := rand.Intn(7)
		card := cardType[num]
		if nowCardsNum[card] <= 0 {
			continue
		}

		nowCardsNum[card]--
		totalCard = append(totalCard, card)

		i++
	}

	// 把桌面的牌抽5張出來
	var deskCard []string
	for i := 0; i < 5; i++ {
		deskCard = append(deskCard, totalCard[len(totalCard)-1])
		totalCard = totalCard[0 : len(totalCard)-1]
	}

	j.foldCard = totalCard
	j.deskCard = deskCard

	// 推播
	j.BroadcastInfo(nil)

	// DONE
	fmt.Printf("init Done!! \n")

	// Debug
	fmt.Printf("Foldcards: %v \n", j.foldCard)
	fmt.Printf("DeskCard: %v \n", j.deskCard)
	fmt.Printf("CardsPoint: %v \n", j.cardsPoint)
	fmt.Printf("PlayerCard: %v \n", j.players)
	fmt.Printf("ActionPlayer: %d \n", j.actionPlayer)
	fmt.Printf("FoldCardNum: %d \n", int32(len(j.foldCard)))
	fmt.Printf("*********\n")
}

// Action 行為
func (j *Jaipur) Action(userID int32, act Action) error {
	// 判斷是不是這個人的局合
	if j.players[j.actionPlayer].ID != userID {
		return errors.New("Not actionplayer")
	}

	// 拿牌
	if act.Type == "take" {
		if err := j.takeCard(userID, act.Take); err != nil {
			return err
		}
	}
	// 賣牌
	if act.Type == "sell" {
		if err := j.sellCard(userID, act.Sell); err != nil {
			return err
		}
	}
	// 換牌 1.手牌換 2.camel換
	if act.Type == "switch" {
		if err := j.switchCard(userID, act.SwitchSelfCard, act.SwitchTargetCard); err != nil {
			return err
		}
	}

	// 換牌以外，都要判斷輸贏
	// 換另外一個的人局合
	j.switchPlayer()

	return nil
}

// 換下一個玩家
func (j *Jaipur) switchPlayer() {
	if j.actionPlayer == 0 {
		j.actionPlayer = 1
	} else {
		j.actionPlayer = 0
	}
}

func (j *Jaipur) getPlayer(userID int32) (*PlayerInfo, error) {
	for k, player := range j.players {
		if player.ID == userID {
			return &j.players[k], nil
		}
	}

	err := errors.New("No User")

	return nil, err
}

func (j *Jaipur) takeCard(userID int32, takeCardKey int32) error {
	// 判斷有沒有這張卡
	if int32(len(j.deskCard)) <= takeCardKey {
		return errors.New("No this Card!")
	}

	// 把卡插入到那個玩家的牌組裡面
	takeCards := []string{}
	takeCardKeys := []int32{}
	takeCard := j.deskCard[takeCardKey]

	// 如果卡是camel就全拿
	if takeCard == "camel" {
		for k, v := range j.deskCard {
			if v == "camel" {
				takeCards = append(takeCards, v)
				takeCardKeys = append(takeCardKeys, int32(k))
			}
		}
	} else {
		takeCards = append(takeCards, takeCard)
		takeCardKeys = append(takeCardKeys, takeCardKey)
	}
	fmt.Printf("takeCards: %v \n", takeCards)
	fmt.Printf("takeCardKeys: %v \n", takeCardKeys)

	player, err := j.getPlayer(userID)
	if err != nil {
		return err
	}

	// camel不會放到手牌上
	if takeCard == "camel" {
		player.Camel = player.Camel + int32(len(takeCards))
	} else {
		// 7張就滿
		if (len(player.Card) + len(takeCards)) > 7 {
			return errors.New("Full Card!")
		}
		player.Card = append(player.Card, takeCards...)
	}

	// 一張一張抽新的，把舊的換掉或刪掉
	for _, cardPosition := range takeCardKeys {
		j.takeOneAndPushNewOneCard(cardPosition)
	}
	fmt.Println("********")
	return nil
}

func (j *Jaipur) sellCard(userID int32, cards []int32) error {
	player, err := j.getPlayer(userID)
	if err != nil {
		return err
	}

	if len(cards) <= 0 {
		return errors.New("At least sell 1 card")
	}

	var sellCards []string
	// 檢查有沒有超出範圍
	for _, v := range cards {
		if v >= int32(len(player.Card)) {
			return errors.New("CardNum Error")
		}

		sellCards = append(sellCards, player.Card[v])
	}

	// 檢查卡是否同一個類型
	sellCardType := sellCards[0]
	for _, v := range sellCards {
		if sellCardType != v {
			return errors.New("Card mush Be same")
		}
	}

	// 這是否合法的卡
	errorCard := true
	for _, v := range allSellCard {
		if v == sellCardType {
			errorCard = false
			break
		}
	}
	if errorCard {
		return errors.New("This card can't sell: " + sellCardType)
	}

	// 如果是sliver gold diamond 一定要賣兩張以上
	if (sellCardType == "sliver" || sellCardType == "gold" || sellCardType == "diamond") && len(cards) < 2 {
		return errors.New("sliver or gold or diamond must be selled two or more in same time")
	}

	// 賣掉換分
	var point int32
	pointList := j.cardsPoint[sellCardType]
	pointNum := 0
	i := 0
	safePoint := 0
	for safePoint <= 20 {
		if pointList[i] != 0 {
			point += pointList[i]
			// point被拿走了，歸0
			pointList[i] = 0
			pointNum++
		}

		if pointNum >= len(cards) {
			break
		}

		i++
		safePoint++
	}

	fmt.Println(point)

	// 卡拿掉
	var newPlayerCards []string
	for valK, cardV := range player.Card {
		notDel := true
		// 如果是要賣的卡，就不要推到新的牌組
		for _, delK := range cards {
			if int32(valK) == delK {
				notDel = false
				break
			}

			if notDel {
				newPlayerCards = append(newPlayerCards, cardV)
			}
		}
	}
	player.Card = newPlayerCards

	// Player 加分數
	player.Point += point

	// 如果賣>=3張就抽bonus
	if len(cards) >= 3 {
		var bouns int32
		var bounsType string

		switch len(cards) {
		case 3:
			bounsType = "three_bonus"
			break
		case 4:
			bounsType = "four_bonus"
			break
		case 5:
			bounsType = "five_bonus"
			break
		}

		bonusList := bonusNum[bounsType]
		randomNum := rand.Intn(len(bonusList))
		bouns = bonusList[randomNum]

		player.Bouns = append(player.Bouns, bouns)
	}

	// TODO 贏輸判斷

	return nil
}

func (j *Jaipur) switchCard(userID int32, selfCards []int32, targetCards []int32) error {
	// 一定要大於1張
	if len(selfCards) < 2 || len(targetCards) < 2 {
		return errors.New("At least switch 2 cards")
	}

	// 檢查交換的數量有沒有一樣
	if len(selfCards) != len(targetCards) {
		return errors.New("selfCards targetCards num must be equal ")
	}

	player, err := j.getPlayer(userID)
	if err != nil {
		return err
	}

	// 檢查有沒有超出範圍
	for i := 0; i < len(selfCards); i++ {
		if targetCards[i] >= int32(len(j.deskCard)) || selfCards[i] >= int32(len(player.Card)) {
			return errors.New("CardNum Error")
		}
	}

	// 交換
	// selfCards帶-1就是camel
	for i := 0; i < len(selfCards); i++ {
		if selfCards[i] == -1 {
			// 檢查camel的數量夠不夠
			if player.Camel < int32(len(selfCards)) {
				return errors.New("camel num Error")
			}
			// 有沒有滿手牌
			if len(player.Card)+1 > 7 {
				return errors.New("full cards")
			}

			player.Card = append(player.Card, j.deskCard[targetCards[i]])
			j.takeOneAndPushNewOneCard(targetCards[i])

			// 減掉玩家的camel數
			player.Camel -= int32(len(selfCards))
		} else {
			// 直接交換
			deskCardTmp := j.deskCard[targetCards[i]]
			j.deskCard[targetCards[i]] = player.Card[selfCards[i]]
			player.Card[selfCards[i]] = deskCardTmp
		}
	}

	return nil
}

// JudgeWinOrLoss 判斷輸贏
func (j *Jaipur) JudgeWinOrLoss() {
	// 1.牌拿光
	// 2.三個類型的牌被賣光
	// 達到條件後，那一回合後馬上game over
}

func (j *Jaipur) takeOneAndPushNewOneCard(cardPosition int32) {
	// 抽一張新的
	var newCard string
	if len(j.foldCard) > 0 {
		newCard = j.foldCard[(len(j.foldCard) - 1)]
		j.foldCard = append(j.foldCard[:(len(j.foldCard) - 1)])
	}

	// 如果還有新片，就把拿掉的卡換掉，否則就刪掉就好
	if newCard != "" {
		fmt.Printf("new card: %s & position: %d \n", newCard, cardPosition)
		j.deskCard[cardPosition] = newCard
	} else {
		j.deskCard = append(j.deskCard[:cardPosition], j.deskCard[cardPosition+1:]...)
	}
}
