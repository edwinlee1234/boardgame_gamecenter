package main

import (
	"errors"

	pb "./proto"

	jaipurClass "./games/jaipur"
)

func newCenter() *Center {
	return &Center{
		jaipurHub: jaipurClass.NewHub(wsURL),
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
	case gameType == "jaipur":
		c.jaipurHub.NewGame(gameID)
		if err := c.jaipurHub.Info(userID, gameID); err != nil {
			return err
		}
		break
	default:
		return errors.New("No this game")
	}

	return nil
}

// CreateGame Center 創立新遊戲
func (c *Center) CreateGame(gameID int32, gameType string, players *pb.Players) error {
	usersInfo := convertUsersInfo(players)

	switch {
	case gameType == "jaipur":
		c.jaipurHub.NewGame(gameID)
		if err := c.jaipurHub.Init(gameID, usersInfo); err != nil {
			return err
		}
		break
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
