package model

// JaipurResult db table
type JaipurResult struct {
	GameID          int32
	Player1ID       int32
	Player2ID       int32
	WinnerID        int32
	CreateTimestamp string
}

// InsertJaipurResult 新增一局遊戲
func InsertJaipurResult(gameID int32, player1ID int32, player2ID int32, winnerID int32, extraInfo []byte) (int32, error) {
	stmt, err := DB.Prepare("INSERT INTO jaipur_result (game_id, player1_id, player2_id, winner_id, extra_info) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	val, err := stmt.Exec(gameID, player1ID, player2ID, winnerID, extraInfo)
	if err != nil {
		return 0, err
	}

	id, err := val.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}
