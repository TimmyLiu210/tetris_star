package v1

import (
	"context"
	"encoding/json"
	"tetris/constant"
	"tetris/postgresql"
	"tetris/redis"

	"github.com/lithammer/shortuuid"
	"gopkg.in/olahol/melody.v1"
)

// 註冊
func SignUp(m *melody.Melody, s *melody.Session, msg []byte) {
	var params struct {
		EventType  int                   `json:"event_type"`
		PlayerInfo postgresql.PlayerInfo `json:"data"`
	}

	ids := []string{GetSessionID(s)}
	err := json.Unmarshal(msg, &params)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_PARAMS, ids)
	}

	if !constant.OverOrEmpty(params.PlayerInfo.Account, constant.MAXACCOUNT) {
		constant.ResponeWithData(m, s, constant.ERROR_ACCCOUNT_OVER_OR_EMPTY, ids)
		return
	}

	if !constant.OverOrEmpty(params.PlayerInfo.Password, constant.MAXPASSWORD) {
		constant.ResponeWithData(m, s, constant.ERROR_PASSWORD_OVER_OR_EMPTY, ids)
		return
	}

	if !constant.OverOrEmpty(params.PlayerInfo.NickName, constant.MAXNICKNAME) {
		constant.ResponeWithData(m, s, constant.ERROR_NICKNAME_OVER_OR_EMPTY, ids)
		return
	}
	params.PlayerInfo.PlayerID = constant.PLAYERPREFIX + shortuuid.New()
	if err = params.PlayerInfo.Register(context.Background(), postgresql.PoolWr.Write()); err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_REGISTER, ids, err.Error())
		return
	}

	redis.InitializePlayerInfo(ids[0], params.PlayerInfo)

	constant.ResponeWithData(m, s, constant.SUCCESS, nil)
	return
}

// 登入
func SignIn(m *melody.Melody, s *melody.Session, msg []byte) {
	var params struct {
		EventType  int                   `json:"event_type"`
		PlayerInfo postgresql.PlayerInfo `json:"data"`
	}
	ids := []string{GetSessionID(s)}
	err := json.Unmarshal(msg, &params)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_PARAMS, ids)
		return
	}

	if !constant.OverOrEmpty(params.PlayerInfo.Account, constant.MAXACCOUNT) {
		constant.ResponeWithData(m, s, constant.ERROR_ACCCOUNT_OVER_OR_EMPTY, ids)
		return
	}

	if !constant.OverOrEmpty(params.PlayerInfo.Password, constant.MAXPASSWORD) {
		constant.ResponeWithData(m, s, constant.ERROR_PASSWORD_OVER_OR_EMPTY, ids)
		return
	}

	playerInfo, err := params.PlayerInfo.GetPlayer(context.Background(), postgresql.PoolWr.Read())
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR, ids, err.Error())
		return
	}

	redis.InitializePlayerInfo(ids[0], *playerInfo)

	constant.ResponeWithData(m, s, constant.SUCCESS, ids)
	return
}

// 登出
func SignOut(m *melody.Melody, s *melody.Session, msg []byte) {
	var params struct {
		EventType  int                   `json:"event_type"`
		PlayerInfo postgresql.PlayerInfo `json:"data"`
	}
	ids := []string{GetSessionID(s)}
	err := json.Unmarshal(msg, &params)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_PARAMS, ids)
		return
	}

	_, err = redis.LeavePlace(ids[0])
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR, ids)
		return
	}

	_, err = redis.DelPlayerInfo(ids[0])
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR, ids)
		return
	}
	
	constant.ResponeWithData(m, s, constant.SUCCESS, ids)
	return
}

func InRoom(m *melody.Melody, s *melody.Session, msg []byte) {
	var params struct {
		EventType int    `json:"event_type"`
		Room      string `json:"room"`
	}
	var (
		ids = []string{GetSessionID(s)}
		err = json.Unmarshal(msg, &params)
	)

	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_PARAMS, ids)
		return
	}

	if !redis.RoomCheck(params.Room) {
		constant.ResponeWithData(m, s, constant.ERROR_EXIST_ROOM, ids)
		return
	}

	if err := redis.ChangePlayerPlace(ids[0], params.Room); err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_IN_ROOM, ids)
		return
	}

	enemy, ownerCheck, err := redis.GetEnemy(ids[0], params.Room)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_GET_ROOM_ENEMY, ids)
		return
	}

	var (
		enemyID            string
		enemyBroadcastInfo postgresql.PlayerInfo
	)
	if !ownerCheck {
		_, err = redis.SetPlayerInfoString(ids[0], constant.PLAYERISROOMOWNER, constant.OWNERSTATETRUE)
		if err != nil {
			constant.ResponeWithData(m, s, constant.ERROR_SET_ROOM_OWNER, ids)
			return
		}
	} else {
		enemyID = enemy.PlayerID
		enemyBroadcastInfo, _, err = redis.GetEnemy(enemyID, params.Room)
		if err != nil {
			constant.ResponeWithData(m, s, constant.ERROR_GET_ROOM_ENEMY, ids)
			return
		}
	}

	err = BroadcastRoomStateList(m, s)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_BROADCASE_PLAYERS, ids)
		return
	}
	if !ownerCheck {
		constant.ResponeWithData(m, s, constant.SUCCESS, ids, nil)
		return
	} else {
		constant.ResponeWithData(m, s, constant.GET_ENEMY, []string{enemyID}, enemyBroadcastInfo)
		constant.ResponeWithData(m, s, constant.SUCCESS, ids, enemy)
		return
	}
}

func OutRoom(m *melody.Melody, s *melody.Session, msg []byte) {
	var params struct {
		EventType int    `json:"event_type"`
		Room      string `json:"room"`
	}
	ids := []string{GetSessionID(s)}

	err := json.Unmarshal(msg, &params)

	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_PARAMS, ids)
		return
	}

	if !redis.RoomCheck(params.Room) {
		constant.ResponeWithData(m, s, constant.ERROR_EXIST_ROOM, ids)
		return
	}

	if err := redis.ChangePlayerPlace(ids[0], params.Room); err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_IN_ROOM, ids)
		return
	}

	err = BroadcastRoomStateList(m, s)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_BROADCASE_PLAYERS, ids)
		return
	}

	constant.ResponeWithData(m, s, constant.SUCCESS, ids)
	return
}

func StartGame(m *melody.Melody, s *melody.Session, msg []byte) {
	return
}
