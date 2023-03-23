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

	BroadcastToPlayers(m,[]byte("log in"), constant.PLACEHALL)

	constant.ResponeWithData(m, s, constant.SUCCESS, ids)
	return
}

// 登出
func SignOut(m *melody.Melody, s *melody.Session, msg []byte) {
	return
}

func InRoom(m *melody.Melody, s *melody.Session, msg []byte) {
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
		constant.ResponeWithData(m, s, constant.ERROR_EXIST_ROOM, nil)
		return
	}

	if err := redis.ChangePlayerPlace(ids[0], params.Room); err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_IN_ROOM, nil)
		return
	}
	constant.ResponeWithData(m, s, constant.SUCCESS, nil)
	return
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
		constant.ResponeWithData(m, s,constant.ERROR_EXIST_ROOM, nil)
		return
	}

if err := redis.ChangePlayerPlace(ids[0], params.Room); err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_OUT_ROOM, nil)
		return
	}
	constant.ResponeWithData(m, s, constant.SUCCESS, nil)
	return
}

func StartGame(m *melody.Melody, s *melody.Session, msg []byte) {
	return
}
