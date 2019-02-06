package redis

import "strconv"

// DelGame DelGame
func DelGame(gameID int32) {
	key := gameInfoRedisPrefix(gameID)
	Client.Del(key)
}

func gameInfoRedisPrefix(gameID int32) string {
	return "game_info:" + strconv.Itoa(int(gameID))
}
