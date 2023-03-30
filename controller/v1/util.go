package v1

import (
	"encoding/json"
	"tetris/constant"
	"tetris/redis"

	"gopkg.in/olahol/melody.v1"
)

func BroadcastToPlayers(m *melody.Melody, msg []byte, room string) error {
	p, err := redis.GetRoomPlayers(room)
	if err != nil {
		return err
	}

	for _, player := range p {
		BroadcastToPlayer(m, msg, player)
	}
	return nil
}

func BroadcastToPlayer(m *melody.Melody, msg []byte, id string) {
	m.BroadcastFilter(msg, func(s *melody.Session) bool {
		compareID, _ := s.Get(constant.SESSIONPREFIX)
		return compareID == id
	})
}

func BroadcastRoomStateList(m *melody.Melody, s *melody.Session) error {
	roomStateList, err := redis.GetAllRoomState()
	if err != nil {
		return err
	}

	rsJsonMsg, err := json.Marshal(roomStateList)
	if err != nil {
		return err
	}

	err = BroadcastToPlayers(m, rsJsonMsg, constant.PLACEHALL)
	if err != nil {
		return err
	}
	return nil
}
