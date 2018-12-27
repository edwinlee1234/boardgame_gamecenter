package jaipur

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// BroadcastRequest 推播的request格式
type BroadcastRequest struct {
	ChannelID int32  `json:"channed_id"`
	Data      []byte `json:"data"`
}

// BroadcastUserRequest 推播單一user的request格式
type BroadcastUserRequest struct {
	ChannelID int32  `json:"channed_id"`
	UUID      string `json:"UUID"`
	Data      []byte `json:"data"`
}

// NewHub NewHub
func NewHub(WsAPI string) *JaipurHub {
	return &JaipurHub{
		Hub:   make(map[int32]*Jaipur),
		WsAPI: WsAPI,
	}
}

// JaipurHub 放Jaipur的遊戲
type JaipurHub struct {
	Hub   map[int32]*Jaipur
	WsAPI string
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
func (j *JaipurHub) Action(gameID int32, userID int32, act Action) error {
	if err := j.checkGame(gameID); err != nil {
		return errors.New("No Game")
	}

	JaipurClass := j.Hub[gameID]
	if err := JaipurClass.Action(userID, act); err != nil {
		return err
	}

	return nil
}

// BroadcastChannel BroadcastChannel
func (j *JaipurHub) BroadcastChannel(channelID int32, data []byte) {
	var req BroadcastRequest
	req.ChannelID = channelID
	req.Data = data

	jsonValue, _ := json.Marshal(req)
	_, err := http.Post(j.WsAPI+"/broadcast", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Println(err)
	}
}

// BroadcastUser BroadcastUser
func (j *JaipurHub) BroadcastUser(channelID int32, UUID string, data []byte) {
	var req BroadcastUserRequest
	req.ChannelID = channelID
	req.UUID = UUID
	req.Data = data

	jsonValue, _ := json.Marshal(req)
	_, err := http.Post(j.WsAPI+"/broadcastUser", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Println(err)
	}
}
