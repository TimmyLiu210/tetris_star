package constant

import (
	"encoding/json"
	"log"
	"tetris/config"

	"github.com/spf13/viper"
	"gopkg.in/olahol/melody.v1"
)

type Message struct {
	EventType int         `json:"event_type"`
	Data      interface{} `json:"data"`
}

func ReadConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault("PORT", ":5000")
	viper.SetDefault("RUN_MODE", "debug")
	viper.SetDefault("READ_TIMEOUT", 1000)
	viper.SetDefault("WRITE_TIMEOUT", 1000)
	viper.SetDefault("REQUEST_TIMEOUT", 1000)
	viper.SetDefault("SHUTDOWN_TIMEOUT", 1000)

	envs := []string{
		"PORT",
		"RUN_MODE",
		"READ_TIMEOUT",
		"WRITE_TIMEOUT",
		"REQUEST_TIMEOUT",
		"SHUTDOWN_TIMEOUT",
	}

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			log.Println(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	config.Initialize()
}

func ByteToMessage(msg []byte) Message {
	var msgInfo Message
	err := json.Unmarshal(msg, &msgInfo)
	if err != nil {
		log.Println("byte to message err", err)
	}
	return msgInfo
}

func StringToByte(s ReturnMessage) ([]byte, error) {
	msg, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func OverOrEmpty(info string, length int) bool {
	if info == "" || len(info) > length {
		return false
	}
	return true
}

func BroadcastToPlayer(m *melody.Melody, msg []byte, id string) {
	m.BroadcastFilter(msg, func(s *melody.Session) bool {
		compareID, _ := s.Get(SESSIONPREFIX)
		return compareID == id
	})
}
