package redis

import (
	"errors"
	"fmt"
	"tetris/constant"
	"tetris/postgresql"

	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

const (
	PLAYER_SIGN_IN = iota + 601
)

var roomList = []string{"hall"}

var msgFlags = []string{
	PLAYER_SIGN_IN: "Player is in the hall",
}

func Initialize() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Panic("[info] redis link Error! err:", err)
	}

	// set redis empty
	redisClient.FlushAll(ctx)

	// set redis hall

	for i := 1; i <= constant.MAXROOMCOUNT; i++ {
		roomList = append(roomList, constant.ROOMPREFIX+fmt.Sprint(i))
	}

	for _, room := range roomList {
		redisClient.SAdd(ctx, room, room)
	}

}

// 初始化redis玩家列表資料
func InitializePlayerInfo(id string, playerInfo postgresql.PlayerInfo) error {
	_, err := redisClient.RPush(ctx, id, playerInfo.PlayerID, playerInfo.Icon,
		playerInfo.Account, playerInfo.Password, playerInfo.NickName, playerInfo.Win, constant.PLACEHALL, constant.OWNERSTATEFALSE).Result()
	if err != nil {
		return err
	}
	_, err = redisClient.SAdd(ctx, constant.PLACEHALL, id).Result()
	if err != nil {
		return err
	}
	return nil
}

// 刪除redis玩家資料
func DelPlayerInfo(id string) (int64, error) {
	return redisClient.Del(ctx, id).Result()
}

// 獲得redis玩家列表長度 list
func GetInfoLen(id string) (int, error) {
	leng, err := redisClient.LLen(ctx, id).Result()
	if err != nil {
		return 0, err
	}
	return int(leng), nil
}

// 找到現在所在的房間
func GetNowPlayerPlace(id string) (string, error) {
	place, err := redisClient.LIndex(ctx, id, int64(constant.NOWPLAYERPLACE)).Result()
	if err != nil {
		return "", err
	}
	return place, nil
}

// 檢查是不是房主
func IsRoomOwner(id, room string) (bool, error) {
	roomOwner, err := redisClient.LIndex(ctx, room, 0).Result()
	if err != nil {
		return false, err
	}

	return (roomOwner == id), nil
}

func SetPlayerInfoString(id string, setIndex int, setInfo string) (string, error) {
	return redisClient.LSet(ctx, id, int64(setIndex), setInfo).Result()
}

func ChangePlayerPlace(id, next string) error {
	now, err := GetNowPlayerPlace(id)
	if err != nil {
		return err
	}
	if now == next {
		return errors.New("[error] player wil go into the same Room")
	}

	_, err = SetPlayerInfoString(id, constant.NOWPLAYERPLACE, next)
	if err != nil {
		return err
	}

	err = changePlace(id, now, next)
	if err != nil {
		return err
	}
	return nil
}

func changePlace(id, now, next string) error {
	_, err := redisClient.SRem(ctx, now, id).Result()
	if err != nil {
		return err
	}

	_, err = redisClient.SAdd(ctx, next, id).Result()
	if err != nil {
		return err
	}

	return nil
}

func LeavePlace(id string) (int64, error) {
	now, err := GetNowPlayerPlace(id)
	if err != nil {
		return 0, err
	}

	return redisClient.SRem(ctx, now, id).Result()
}

func RoomCheck(nextRoom string) bool {
	for _, room := range roomList {
		if room == nextRoom {
			return true
		}
	}
	return false
}

func GetRoomPlayers(room string) ([]string, error) {
	ids, err := redisClient.SMembers(ctx, room).Result()
	if err != nil {
		return []string{}, err
	}
	return ids, nil
}
