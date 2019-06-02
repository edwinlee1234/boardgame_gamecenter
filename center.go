package main

import (
	"errors"
	"log"

	lib "boardgame_gamecenter/lib"
	pb "boardgame_gamecenter/proto"

	jaipurClass "boardgame_gamecenter/games/jaipur"
)

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
	case gameType == "jaipur":
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
	case gameType == "jaipur":
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
func (c *Center) CreateGame(gameID int32, gameType string, players *pb.Players) (err error) {
	usersInfo := convertUsersInfo(players)
	switch {
	case gameType == "jaipur":
		// TODO 這些邏輯都帶到c.jaipurHub.NewGame裡面，不要放在這裡，用err來判斷對不對就好了
		// 檢查人數
		if err = c.jaipurHub.CheckUserValid(usersInfo); err != nil {
			log.Printf("%v", err)
			return err
		}

		c.jaipurHub.NewGame(gameID)
		if err = c.jaipurHub.Init(gameID, usersInfo); err != nil {
			log.Printf("%v", err)
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
