package v1

import (
	"encoding/json"
	"log"
	"tetris/constant"
	"tetris/redis"

	"gopkg.in/olahol/melody.v1"
)

func GetCommond(m *melody.Melody, s *melody.Session, msg []byte) {
	var params struct {
		EventType int `json:"event_type"`
		Commond   int `json:"commond"`
	}

	ids := []string{GetSessionID(s)}

	err := json.Unmarshal(msg, &params)
	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR_PARAMS, ids)
		return
	}

	room, err := redis.GetNowPlayerPlace(ids[0])

	if err != nil {
		constant.ResponeWithData(m, s, constant.ERROR, ids)
	}

	log.Println("room:", room)
}
