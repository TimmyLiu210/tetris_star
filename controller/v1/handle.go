package v1

import (
	"encoding/json"
	"log"
	"tetris/constant"

	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
)

type Message struct {
	EventType int         `json:"event_type"`
	Data      interface{} `json:"data"`
}

func InitSession(s *melody.Session) string {
	id := uuid.New().String()
	s.Set(constant.SESSIONPREFIX, id)
	return id
}

func GetSessionID(s *melody.Session) string {
	if id, isExist := s.Get(constant.SESSIONPREFIX); isExist {
		return id.(string)
	}
	return InitSession(s)
}

func SetPlayerPlace(s *melody.Session, place string) {
	s.Set(string(constant.PLACE), place)
	return
}

func CheckCommond(m *melody.Melody, s *melody.Session, msg []byte) {
	var (
		message Message
		ids     = []string{GetSessionID(s)}
	)

	json.Unmarshal(msg, &message)
	log.Println("event_type", message.EventType)
	switch message.EventType {
	case constant.SIGN_UP:
		SignUp(m, s, msg)
	case constant.SIGN_IN:
		SignIn(m, s, msg)
	case constant.SIGN_OUT:
		SignOut(m, s, msg)
	case constant.IN_ROOM:
		InRoom(m, s, msg)
	case constant.OUT_ROOM:
		OutRoom(m, s, msg)
	case constant.START_GAME:
		StartGame(m, s, msg)
	case constant.GAME_COMMOND:
		GetCommond(m, s, msg)
	default:
		constant.ResponeWithData(m, s, constant.ERROR_GET_MESSAGE, ids, nil)
	}
	return
}
