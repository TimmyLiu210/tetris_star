package main

import (
	"log"
	"tetris/config"
	"tetris/constant"
	v1 "tetris/controller/v1"
	"tetris/game"
	"tetris/postgresql"
	"tetris/redis"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func init() {
	constant.ReadConfig(".env")
	postgresql.Initialize()
	redis.Initialize()
	game.Initialize()
}

func main() {
	gin.SetMode(config.RunMode)
	port := config.Port

	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	m.HandleMessage(func(session *melody.Session, msg []byte) {
		v1.CheckCommond(m, session, msg)

	})
	m.HandleConnect(func(session *melody.Session) {
		id := v1.InitSession(session)
		log.Println("[info] Sessoin connect:", id)
	})
	m.HandleClose(func(session *melody.Session, i int, s string) error {
		id := v1.GetSessionID(session)
		_, err := redis.LeavePlace(id)
		if err != nil {
			log.Println("[error] ", err)
			return err
		}
		_, err = redis.DelPlayerInfo(id)
		if err != nil {
			log.Println("[error] ", err)
			return err
		}
		log.Println("[info] Sessoin close:", id)
		return nil
	})

	r.Run(port)
}
