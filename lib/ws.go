package lib

import (
	"bytes"
	"encoding/json"
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

// NewWS NewWS
func NewWS(WsURL string) *WS {
	return &WS{
		WsURL: WsURL,
	}
}

// WS WS
type WS struct {
	WsURL string
}

// BroadcastChannel BroadcastChannel
func (w *WS) BroadcastChannel(channelID int32, data []byte) {
	var req BroadcastRequest
	req.ChannelID = channelID
	req.Data = data

	jsonValue, _ := json.Marshal(req)
	_, err := http.Post(w.WsURL+"/broadcast", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Println(err)
	}
}

// BroadcastUser BroadcastUser
func (w *WS) BroadcastUser(channelID int32, UUID string, data []byte) {
	var req BroadcastUserRequest
	req.ChannelID = channelID
	req.UUID = UUID
	req.Data = data

	jsonValue, _ := json.Marshal(req)
	_, err := http.Post(w.WsURL+"/broadcastUser", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Println(err)
	}
}
