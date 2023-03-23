package v1

import (
	"log"
	"tetris/constant"
	"tetris/redis"

	"gopkg.in/olahol/melody.v1"
)

func BroadcastToPlayers(m *melody.Melody, msg []byte, room string) error {
	p, err := redis.GetRoomPlayers(room)
	if err != nil {
		return err
	}

	log.Println(room, "members: ", p)
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
