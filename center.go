package main

import (
	"errors"
	"fmt"
	"log"

	lib "boardgame_gamecenter/lib"
	pb "boardgame_gamecenter/proto"

	jaipurClass "boardgame_gamecenter/games/jaipur"
)

// TODO 用這種方式改寫遊戲
type GamesHub interface {
	NewGame(gameID int32, usersInfo map[int32]string, extraInfo map[string]interface{}) error
	Info(userID []int32, gameID int32) error
	Action(userID int32, gameID int32, act interface{}) error
}

type CenterV2 struct {
	gameshub map[string]GamesHub
}

func newCenterV2(WS *lib.WS) *CenterV2 {
	center := &CenterV2{
		gameshub: map[string]GamesHub{
			JAIPUR: jaipurClass.NewHub(WS),
		},
	}

	return center
}

// GameInfo 取得遊戲資訊
func (c *CenterV2) GameInfo(userID []int32, gameID int32, gameType string) error {
	class, exist := c.gameshub[gameType]
	if !exist {
		return fmt.Errorf("no this game %v", gameType)
	}

	return class.Info(userID, gameID)
}

// CreateGame Center 創立新遊戲
func (c *CenterV2) CreateGame(gameID int32, gameType string, players *pb.Players, extraInfo map[string]interface{}) (err error) {
	usersInfo := convertUsersInfo(players)

	class, exist := c.gameshub[gameType]
	if !exist {
		return fmt.Errorf("no this game %v", gameType)
	}

	return class.NewGame(gameID, usersInfo, extraInfo)
}

// ActionProcess ActionProcess
func (c *CenterV2) ActionProcess(userID int32, gameID int32, gameType string, action interface{}) error {
	class, exist := c.gameshub[gameType]
	if !exist {
		return fmt.Errorf("no this game %v", gameType)
	}

	return class.Action(userID, gameID, action)
}

func newCenter(WS *lib.WS) *Center {
	return &Center{
		jaipurHub: jaipurClass.NewHub(WS),
	}
}

// Center 遊戲中心
// TODO 這樣寫不太好擴充，會一直壘加在switch裡面，擴充會動到舊有的程式
type Center struct {
	jaipurHub *jaipurClass.JaipurHub
	// 其他遊戲，等擴充
}

// GameInfo 取得遊戲資訊
func (c *Center) GameInfo(userID []int32, gameID int32, gameType string) error {
	switch {
	case gameType == JAIPUR:
		if err := c.jaipurHub.Info(userID, gameID); err != nil {
			return err
		}
		break
	default:
		return errors.New("No this game")
	}

	return nil
}

// ActionProcess ActionProcess
func (c *Center) ActionProcess(userID int32, gameID int32, gameType string, action interface{}) error {
	switch {
	case gameType == JAIPUR:
		if err := c.jaipurHub.Action(userID, gameID, action); err != nil {
			log.Printf("%v", err)
			return err
		}
		break
	default:
		return errors.New("No this game")
	}

	return nil
}

// CreateGame Center 創立新遊戲
func (c *Center) CreateGame(gameID int32, gameType string, players *pb.Players, extraInfo map[string]interface{}) (err error) {
	usersInfo := convertUsersInfo(players)
	switch {
	case gameType == JAIPUR:
		if err := c.jaipurHub.NewGame(gameID, usersInfo, extraInfo); err != nil {
			return err
		}
	default:
		return errors.New("No this game")
	}

	return nil
}

func convertUsersInfo(players *pb.Players) map[int32]string {
	usersInfo := make(map[int32]string)

	for _, player := range players.PlayerList {
		usersInfo[player.ID] = player.UUID
	}

	return usersInfo
}
