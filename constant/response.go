package constant

import (
	"gopkg.in/olahol/melody.v1"
)

const (
	SUCCESS int = iota + 200
	
	GET_ENEMY
)
const (
	ERROR   int = 500
	ERROR_GET_MESSAGE = iota + 501
	ERROR_PARAMS
	ERROR_ACCCOUNT_OVER_OR_EMPTY
	ERROR_PASSWORD_OVER_OR_EMPTY
	ERROR_NICKNAME_OVER_OR_EMPTY
	ERROR_ACCOUNT_NO_EXIST
	ERROR_PASSWORD_NO_CURRECT
	ERROR_UPSERT
	ERROR_REGISTER
	ERROR_EXIST_ROOM
	ERROR_IN_ROOM
	ERROR_OUT_ROOM

	ERROR_GET_ROOM_PLAYERS
	ERROR_GET_ROOM_ENEMY
	ERROR_GET_ROOM_STATE
	ERROR_GET_ROOM_STATE_LIST
	ERROR_GET_ROOM_IS_PLAYING

	ERROR_SET_ROOM_OWNER

	ERROR_SESSION_ID

	ERROR_DATA_JSON_MARSHAL

	ERROR_BROADCASE_PLAYERS
)

var msgFlags = map[int]string{
	SUCCESS:                      "Ok",

	GET_ENEMY: "get enemy",

	ERROR:                        "Fail",
	ERROR_PARAMS:                 "Get params failed!",
	ERROR_GET_MESSAGE:            "Get message failed!",
	ERROR_ACCCOUNT_OVER_OR_EMPTY: "Account invalid, over or empty!",
	ERROR_PASSWORD_OVER_OR_EMPTY: "Password invalid, over or empty!",
	ERROR_NICKNAME_OVER_OR_EMPTY: "Nickname invalid, over or empty!",
	ERROR_ACCOUNT_NO_EXIST:       "Account no exist!",
	ERROR_PASSWORD_NO_CURRECT:    "Password no currect!",
	ERROR_UPSERT:                 "Upsert data failed!",
	ERROR_REGISTER:               "Register failed!",
	ERROR_EXIST_ROOM:             "Room doesn't exists!",
	ERROR_IN_ROOM:                "In room failed!",
	ERROR_OUT_ROOM:               "out room failed!",

	ERROR_GET_ROOM_PLAYERS:    "get room players failed!",
	ERROR_GET_ROOM_STATE_LIST: "get room state list failed!",
	ERROR_GET_ROOM_STATE:      "get room state failed!",
	ERROR_GET_ROOM_IS_PLAYING: "get room is playing failed!",
	ERROR_GET_ROOM_ENEMY:      "get room enemy failed!",

	ERROR_SESSION_ID: "get session id failed!",

	ERROR_DATA_JSON_MARSHAL: "set json marshal failed!",

	ERROR_BROADCASE_PLAYERS: "broadcast to players failed!",

	ERROR_SET_ROOM_OWNER: "set room owner failed!",
}

type ReturnMessage struct {
	ReturnCode int         `json:"return_code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

func ResponeWithData(m *melody.Melody, s *melody.Session, returnCode int, ids []string, data ...interface{}) {
	msg, ok := msgFlags[returnCode]
	if !ok {
		msg = msgFlags[ERROR]
	}
	response := ReturnMessage{
		ReturnCode: returnCode,
		Msg:        msg,
		Data:       data,
	}

	jsonMsg, _ := StringToByte(response)

	for _, id := range ids {
		BroadcastToPlayer(m, jsonMsg, id)
	}
}
