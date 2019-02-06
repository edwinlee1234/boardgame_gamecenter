package model

// GameState db table
type GameState struct {
	ID         int32
	Type       string
	State      int32
	Result     string
	Seat       int32
	InsertTime string
	UpdateTime string
}

// GameState

// NotOpen owner only
const NotOpen = 0

// Opening 開放玩家
const Opening = 1

// Playing 遊戲中
const Playing = 2

// Close 關
const Close = 4

// Abort 放棄
const Abort = 5

// ChangeGameStateDB 改變state
func ChangeGameStateDB(id int32, state int32) error {
	stmt, _ := DB.Prepare("UPDATE `game_state` set `state` = ? where `id` = ?")
	res, _ := stmt.Exec(state, id)

	affect, err := res.RowsAffected()
	if err != nil || affect == 0 {
		return err
	}

	return nil
}
