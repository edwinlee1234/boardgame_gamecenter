package jaipur

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Bonus []int32  `json:"bonus"`
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

// TakeEventPublic 拿牌Event
type TakeEventPublic struct {
	Event string              `json:"event"`
	Info  TakeEventPublicInfo `json:"info"`
}

// TakeEventPublicInfo 拿牌Event Info
type TakeEventPublicInfo struct {
	TakeDeskCardNum []int32  `json:"takeDeskCardNum"`
	NewDeskCard     []string `json:"newDeskCard"`
	DeskCard        []string `json:"deskCard"`
	FoldCardNum     int32    `json:"foldCardNum"`
	PlayerID        int32    `json:"playerID"`
	ActionPlayer    int32    `json:"actionPlayer"`
}

// TakeEventPrivate 拿牌Event user (動作的user專用)
type TakeEventPrivate struct {
	Event string               `json:"event"`
	Info  TakeEventPrivateInfo `json:"info"`
}

// TakeEventPrivateInfo 拿牌Event Info (動作的user專用)
type TakeEventPrivateInfo struct {
	Card  []string `json:"card"`
	Camel int32    `json:"camel"`
}

// SellEventPublic 拿牌Event
type SellEventPublic struct {
	Event string              `json:"event"`
	Info  SellEventPublicInfo `json:"info"`
}

// SellEventPublicInfo 拿牌Event Info
type SellEventPublicInfo struct {
	CardsPoint   map[string][]int32 `json:"cardsPoint"`
	SellCards    []string           `json:"sellcards"`
	PlayerID     int32              `json:"playerID"`
	ActionPlayer int32              `json:"actionPlayer"`
}

// SellEventPrivate 賣牌 (動作的user專用)
type SellEventPrivate struct {
	Event string               `json:"event"`
	Info  SellEventPrivateInfo `json:"info"`
}

// SellEventPrivateInfo 賣牌 (動作的user專用)
type SellEventPrivateInfo struct {
	Bonus int32    `json:"bonus"`
	Card  []string `json:"card"`
	Point int32    `json:"point"`
}

// SwitchEventPublic 拿牌Event
type SwitchEventPublic struct {
	Event string                `json:"event"`
	Info  SwitchEventPublicInfo `json:"info"`
}

// SwitchEventPublicInfo 拿牌Event Info
type SwitchEventPublicInfo struct {
	TakeDeskCardNum []int32  `json:"takeDeskCardNum"`
	NewDeskCard     []string `json:"newDeskCard"`
	DeskCard        []string `json:"deskCard"`
	PlayerID        int32    `json:"playerID"`
	ActionPlayer    int32    `json:"actionPlayer"`
}

// SwitchEventPrivate 換牌 (動作的user專用)
type SwitchEventPrivate struct {
	Event string                 `json:"event"`
	Info  SwitchEventPrivateInfo `json:"info"`
}

// SwitchEventPrivateInfo 換牌 (動作的user專用)
type SwitchEventPrivateInfo struct {
	Card  []string `json:"card"`
	Camel int32    `json:"camel"`
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

// OpponentChangeEvent OpponentChangeEvent
type OpponentChangeEvent struct {
	Event    string         `json:"event"`
	Opponent OpponentPlayer `json:"opponentPlayer"`
}

// GameOverEvent GameOverEvent
type GameOverEvent struct {
	Event string            `json:"event"`
	Info  GameOverEventInfo `json:"info"`
}

// GameOverEventInfo GameOverEventInfo
type GameOverEventInfo struct {
	WinnerID           int32   `json:"winnerID"`
	CamelBonusWinnerID int32   `json:"camelBonusWinnerID"`
	Players            Players `json:"players"`
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
	j.cardsPoint = getCardPoint()
	nowCardsNum := getCardNum()
	var gameTotalCard = cardTotal

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

		gameTotalCard -= 5
	}

	// 全部的牌random一次
	i := 0
	for i < gameTotalCard {
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
	var newCard []string
	for _, cardPosition := range takeCardKeys {
		newCard = append(newCard, j.takeOneAndPushNewOneCard(cardPosition))
	}

	// 換行動玩家
	j.switchPlayer()

	fmt.Println("********")

	// 推播 遊戲內的全部人
	publicInfo := TakeEventPublic{
		Event: "TakeCardPublic",
		Info: TakeEventPublicInfo{
			TakeDeskCardNum: takeCardKeys,
			NewDeskCard:     newCard,
			DeskCard:        j.deskCard,
			FoldCardNum:     int32(len(j.foldCard)),
			PlayerID:        player.ID,
			ActionPlayer:    j.players[j.actionPlayer].ID,
		},
	}
	publicJSON, _ := json.Marshal(publicInfo)
	j.hub.BroadcastChannel(j.gameID, publicJSON)

	// 推播 拿卡的玩家
	privateInfo := TakeEventPrivate{
		Event: "TakeCardPrivate",
		Info: TakeEventPrivateInfo{
			Card:  player.Card,
			Camel: player.Camel,
		},
	}
	privateJSON, _ := json.Marshal(privateInfo)
	j.hub.BroadcastUser(j.gameID, player.UUID, privateJSON)

	// 推給對手，手牌的變動
	j.opponentChange(int32(len(player.Card)), player.Camel > 0, player.ID)

	// 檢查遊戲結束了沒
	if j.checkGameOver() {
		// 遊戲結束
		j.gameOverEvent()
		fmt.Println("GameOver!!!!")
	}

	return nil
}

func (j *Jaipur) sellCard(userID int32, cards []int32) error {
	fmt.Println("********")
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
	log.Printf("sellCards: %v", sellCards)
	var point int32
	for _, cardType := range sellCards {
		// 沒分數了
		if len(j.cardsPoint[cardType]) <= 0 {
			break
		}

		point += j.cardsPoint[cardType][0]
		j.cardsPoint[cardType] = j.cardsPoint[cardType][1:]
	}

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
		}

		if notDel {
			newPlayerCards = append(newPlayerCards, cardV)
		}
	}
	log.Printf("old plyer cards: %v", player.Card)
	log.Printf("new player cards: %v", newPlayerCards)
	player.Card = newPlayerCards

	// Player 加分數
	player.Point += point

	// 如果賣>=3張就抽bonus
	var bonus int32
	var bonusType string
	if len(cards) >= 3 {
		switch len(cards) {
		case 3:
			bonusType = "three_bonus"
			break
		case 4:
			bonusType = "four_bonus"
			break
		default:
			bonusType = "five_bonus"
			break
		}

		// TODO bonus有限，這邊要把拿掉的刪掉
		rand.Seed(time.Now().UnixNano()) // run seed
		log.Printf("bonus: %s", bonusType)
		bonusList := bonusNum[bonusType]

		log.Printf("bonusList: %v", bonusList)
		randomNum := rand.Intn(len(bonusList))
		bonus = bonusList[randomNum]
		log.Printf("bonus: %d", bonus)

		player.Bonus = append(player.Bonus, bonus)
		// 加分
		player.Point += bonus
	}

	// 換行動玩家
	j.switchPlayer()

	publicEvent := SellEventPublic{
		Event: "SellEventPublic",
		Info: SellEventPublicInfo{
			CardsPoint:   j.cardsPoint,
			SellCards:    sellCards,
			PlayerID:     player.ID,
			ActionPlayer: j.players[j.actionPlayer].ID,
		},
	}

	publicJSON, _ := json.Marshal(publicEvent)
	j.hub.BroadcastChannel(j.gameID, publicJSON)

	privateEvent := SellEventPrivate{
		Event: "SellEventPrivate",
		Info: SellEventPrivateInfo{
			Bonus: bonus,
			Card:  player.Card,
			Point: player.Point,
		},
	}

	privateJSON, _ := json.Marshal(privateEvent)
	j.hub.BroadcastUser(j.gameID, player.UUID, privateJSON)

	// 推給對手，手牌的變動
	j.opponentChange(int32(len(player.Card)), player.Camel > 0, player.ID)

	// 檢查遊戲結束了沒
	if j.checkGameOver() {
		// 遊戲結束
		j.gameOverEvent()
		fmt.Println("GameOver!!!!")
	}

	fmt.Println("********")

	return nil
}

func (j *Jaipur) switchCard(userID int32, selfCards []int32, targetCards []int32) error {
	// 一定要大於1張
	if len(selfCards) < 2 || len(targetCards) < 2 {
		return errors.New("At least switch 2 cards")
	}

	// 檢查交換的數量有沒有一樣
	if len(selfCards) != len(targetCards) {
		return errors.New("selfCards targetCards num must be equal")
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

	// 檢查交換的手牌跟場上有沒有重復
	for _, selfVal := range selfCards {
		for _, tarVal := range targetCards {
			if selfVal == -1 && j.deskCard[tarVal] == "camel" {
				return errors.New("can't not switch same card!!!!")
			}
			if selfVal != -1 && player.Card[selfVal] == j.deskCard[tarVal] {
				return errors.New("can't not switch same card!!!!")
			}
		}
	}

	// 交換
	// selfCards帶-1就是camel
	var newCard []string
	for i := 0; i < len(selfCards); i++ {
		if selfCards[i] == -1 {
			if player.Camel-1 < 0 {
				return errors.New("camel num Error")
			}

			if len(player.Card)+1 > 7 {
				return errors.New("full cards")
			}

			// camel推到桌上的牌
			player.Card = append(player.Card, j.deskCard[targetCards[i]])
			newCard = append(newCard, "camel")
			j.deskCard[targetCards[i]] = "camel"

			// 減掉玩家的camel數
			player.Camel--
		} else {
			// 不是camel就用手上的牌直接交換
			deskCardTmp := j.deskCard[targetCards[i]]
			newCard = append(newCard, deskCardTmp)
			j.deskCard[targetCards[i]] = player.Card[selfCards[i]]
			player.Card[selfCards[i]] = deskCardTmp
		}
	}

	// 換行動玩家
	j.switchPlayer()

	// 推播 遊戲內的全部人
	publicInfo := SwitchEventPublic{
		Event: "SwitchEventPublic",
		Info: SwitchEventPublicInfo{
			TakeDeskCardNum: targetCards,
			NewDeskCard:     newCard,
			DeskCard:        j.deskCard,
			PlayerID:        player.ID,
			ActionPlayer:    j.players[j.actionPlayer].ID,
		},
	}
	publicJSON, _ := json.Marshal(publicInfo)
	j.hub.BroadcastChannel(j.gameID, publicJSON)

	// 推播 拿卡的玩家
	privateInfo := SwitchEventPrivate{
		Event: "SwitchEventPrivate",
		Info: SwitchEventPrivateInfo{
			Card:  player.Card,
			Camel: player.Camel,
		},
	}
	privateJSON, _ := json.Marshal(privateInfo)
	j.hub.BroadcastUser(j.gameID, player.UUID, privateJSON)

	// 推給對手，手牌的變動
	j.opponentChange(int32(len(player.Card)), player.Camel > 0, player.ID)

	return nil
}

// 判斷輸贏
func (j *Jaipur) judgeWinOrLoss() {

}

// 看遊戲結束了沒
func (j *Jaipur) checkGameOver() bool {
	// 達到條件後，那一回合後馬上game over
	// 1.牌拿光
	if len(j.foldCard) <= 0 {
		return true
	}

	// 2.三個類型的牌被賣光
	three := 0
	for _, list := range j.cardsPoint {
		if len(list) <= 0 {
			three++
		}
	}

	if three == 3 {
		return true
	}

	return false
}

// 推播給對方的卡牌數變化，跟有沒有cammel
func (j *Jaipur) opponentChange(cardsNum int32, haveCamel bool, actionUserID int32) {
	// 找出這個userID的對手
	var pushPlayer PlayerInfo
	for _, v := range j.players {
		if v.ID != actionUserID {
			pushPlayer = v
		}
	}

	event := OpponentChangeEvent{
		Event: "OpponentChangeEvent",
		Opponent: OpponentPlayer{
			cardsNum,
			haveCamel,
		},
	}

	eventJSON, _ := json.Marshal(event)
	j.hub.BroadcastUser(j.gameID, pushPlayer.UUID, eventJSON)
}

// 贏的時候
func (j *Jaipur) gameOverEvent() {
	// 算那個cammel比較多，會多加5分
	player1 := j.players[0]
	player2 := j.players[1]

	var camelWinner int32
	if player1.Camel > player2.Camel {
		player1.Bonus = append(player1.Bonus, 5)
		player1.Point += 5
		camelWinner = player1.ID
	}

	if player2.Camel > player1.Camel {
		player2.Bonus = append(player2.Bonus, 5)
		player2.Point += 5
		camelWinner = player2.ID
	}

	// 把兩邊的分數都公開跟贏家是誰
	var winnerID int32
	if player1.Point > player2.Point {
		winnerID = player1.ID
	}

	if player2.Point > player1.Point {
		winnerID = player2.ID
	}

	// 如果分數都一樣，那有camel bonus的人贏
	if player1.Point == player2.Point {
		if camelWinner == player1.ID {
			winnerID = player1.ID
		}

		if camelWinner == player2.ID {
			winnerID = player2.ID
		}
	}

	// 如果都分不出誰贏，就算場主贏 (第一個玩家)
	if winnerID == 0 {
		winnerID = j.players[0].ID
	}

	// 推播
	event := GameOverEvent{
		Event: "GameOverEvent",
		Info: GameOverEventInfo{
			WinnerID:           winnerID,
			CamelBonusWinnerID: camelWinner,
			Players:            j.players,
		},
	}

	eventJSON, _ := json.Marshal(event)
	j.hub.BroadcastChannel(j.gameID, eventJSON)

	// 寫db
	j.hub.RecordResult(j.gameID, player1.ID, player2.ID, winnerID, []byte{})
	// destory game
	j.hub.GameOver(j.gameID)
}

func (j *Jaipur) takeOneAndPushNewOneCard(cardPosition int32) string {
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

	return newCard
}
