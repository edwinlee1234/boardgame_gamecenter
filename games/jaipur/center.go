package jaipur

import (
	lib "boardgame_gamecenter/lib"
	model "boardgame_gamecenter/model"
	pb "boardgame_gamecenter/proto"
	redis "boardgame_gamecenter/redis"
	"errors"
	"log"
)

// NewHub NewHub
func NewHub(WS *lib.WS) *JaipurHub {
	return &JaipurHub{
		Hub: make(map[int32]*Jaipur),
		WS:  WS,
	}
}

// JaipurHub 放Jaipur的遊戲
type JaipurHub struct {
	Hub map[int32]*Jaipur
	WS  *lib.WS
}

// CheckUserValid CheckUserValid
func (j *JaipurHub) CheckUserValid(usersInfo map[int32]string) error {
	if len(usersInfo) != 2 {
		return errors.New("much be two player")
	}

	return nil
}

// NewGame NewGame
func (j *JaipurHub) NewGame(gameID int32) {
	// 沒有這遊戲ID才開新的
	if err := j.checkGame(gameID); err != nil {
		j.Hub[gameID] = NewJaipur(gameID, j)
	}
}

func (j *JaipurHub) checkGame(gameID int32) error {
	if _, exist := j.Hub[gameID]; !exist {
		return errors.New("No Game")
	}

	return nil
}

// Init Init
func (j *JaipurHub) Init(gameID int32, usersInfo map[int32]string) error {
	if err := j.checkGame(gameID); err != nil {
		return errors.New("No Game")
	}

	JaipurClass := j.Hub[gameID]
	JaipurClass.Init(usersInfo)

	return nil
}

// Info BroadcastInfo
func (j *JaipurHub) Info(userID []int32, gameID int32) error {
	if err := j.checkGame(gameID); err != nil {
		return errors.New("No Game")
	}

	JaipurClass := j.Hub[gameID]
	JaipurClass.BroadcastInfo(userID)

	return nil
}

// Action Action
func (j *JaipurHub) Action(userID int32, gameID int32, act interface{}) error {
	if err := j.checkGame(gameID); err != nil {
		return errors.New("No Game")
	}

	actionPd, res := act.(*pb.JaipurActionStruct)
	if !res {
		return errors.New("Action struct error")
	}

	action := Action{
		Type:             actionPd.Type,
		Take:             actionPd.Take,
		Sell:             actionPd.Sell,
		SwitchSelfCard:   actionPd.SwitchSelfCard,
		SwitchTargetCard: actionPd.SwitchTargetCard,
	}

	JaipurClass := j.Hub[gameID]
	if err := JaipurClass.Action(userID, action); err != nil {
		return err
	}

	return nil
}

// BroadcastChannel BroadcastChannel
func (j *JaipurHub) BroadcastChannel(channelID int32, data []byte) {
	j.WS.BroadcastChannel(channelID, data)
}

// BroadcastUser BroadcastUser
func (j *JaipurHub) BroadcastUser(channelID int32, UUID string, data []byte) {
	j.WS.BroadcastUser(channelID, UUID, data)
}

// GameOver GameOver
func (j *JaipurHub) GameOver(gameID int32) {
	// 刪掉遊戲
	delete(j.Hub, gameID)
	redis.DelGame(gameID)

	err := model.ChangeGameStateDB(gameID, model.Close)
	if err != nil {
		log.Println(err)
	}
}

// RecordResult RecordResult
func (j *JaipurHub) RecordResult(gameID int32, player1ID int32, player2ID int32, winnerID int32, extraInfo []byte) {
	_, err := model.InsertJaipurResult(gameID, player1ID, player2ID, winnerID, extraInfo)

	if err != nil {
		log.Println(err)
	}
}
