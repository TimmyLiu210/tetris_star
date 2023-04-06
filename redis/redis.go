package redis

import (
	"errors"
	"fmt"
	"strconv"
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

var RoomList = []string{"hall"}

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
		RoomList = append(RoomList, constant.ROOMPREFIX+fmt.Sprint(i))
	}

	for _, room := range RoomList {
		redisClient.SAdd(ctx, room, "waiting")
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
	for _, room := range RoomList {
		if room == nextRoom {
			return true
		}
	}
	return false
}

// 取得房間玩家列表
func GetRoomPlayers(room string) ([]string, error) {
	ids, err := redisClient.SMembers(ctx, room).Result()
	if err != nil {
		return []string{}, err
	}
	return ids, nil
}

// 取得對手資料
func GetEnemy(id, room string) (postgresql.PlayerInfo, bool, error) {
	var enemyID string
	var info postgresql.PlayerInfo
	players, err := GetRoomPlayers(room)
	if err != nil {
		return info, false, err
	}

	for _, playerID := range players {
		if playerID != id && playerID != constant.ROOMSTATEPLAYING && playerID != constant.ROOMSTATEWAITING {
			enemyID = playerID
		}
	}
	log.Println("enemy:", enemyID)
	if enemyID == "" {
		return info, false, nil
	}

	icon, err := GetPlayerInfo(id, constant.PLAYERICON)
	if err != nil {
		return postgresql.PlayerInfo{}, false, err
	}

	playerID, err := GetPlayerInfo(id, constant.PLAYERID)
	if err != nil {
		return postgresql.PlayerInfo{}, false, err
	}

	nickName, err := GetPlayerInfo(id, constant.PLAYERNICKNAME)
	if err != nil {
		return postgresql.PlayerInfo{}, false, err
	}

	win, err := GetPlayerInfo(id, constant.PLAYERWIN)
	if err != nil {
		return postgresql.PlayerInfo{}, false, err
	}

	info.PlayerID = playerID
	info.NickName = nickName
	info.Win, err = strconv.Atoi(win)
	if err != nil {
		return postgresql.PlayerInfo{}, false, err
	}
	info.Icon, err = strconv.Atoi(icon)
	if err != nil {
		return postgresql.PlayerInfo{}, false, err
	}

	return info, true, nil
}

// 取得玩家單項資料
func GetPlayerInfo(id string, infoType int) (string, error) {
	info, err := redisClient.LIndex(ctx, id, int64(infoType)).Result()
	if err != nil {
		return "", err
	}
	return info, nil
}

// 取得房間資料+狀態
func GetRoomState(room string) (postgresql.RoomState, error) {
	var roomState postgresql.RoomState

	exist, err := redisClient.SIsMember(ctx, room, constant.ROOMSTATEPLAYING).Result()
	if err != nil {
		return postgresql.RoomState{}, err
	}

	if exist {
		roomState.IsPlaying = true
	} else {
		roomState.IsPlaying = false
	}

	playerList, err := GetRoomPlayers(room)
	if err != nil {
		return postgresql.RoomState{}, err
	}

	roomState.PlayerList = playerList

	return roomState, nil
}

func GetAllRoomState() ([]postgresql.RoomState, error) {
	var allRoomList []postgresql.RoomState

	for _, room := range RoomList {
		if room != constant.PLACEHALL {
			roomState, err := GetRoomState(room)
			if err != nil {
				return []postgresql.RoomState{}, err
			}

			allRoomList = append(allRoomList, roomState)
		}
	}

	return allRoomList, nil
}

func SetRoomPlayerState(room string, state string) {

}

func UpdatePlayerWin(id string) {

}
