package constant

import (
	"gopkg.in/olahol/melody.v1"
)

const (
	SUCCESS int = 200
	ERROR   int = 500
)
const (
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
)

var msgFlags = map[int]string{
	SUCCESS:                      "Ok",
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
	ERROR_GET_ROOM_PLAYERS:       "get room players failed!",
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
