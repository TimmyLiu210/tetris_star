package game

import (
	"errors"
	"math/rand"
	"strconv"
	"tetris/constant"
	"tetris/redis"

	"gopkg.in/olahol/melody.v1"
)

type TetrisCommond struct {
	Player  string
	Commond int
}

type TetrisSite struct {
	TetrisName string

	TetrisType       string
	TetrisRotateType int

	X [constant.TETRIS_COUNT]int
	Y [constant.TETRIS_COUNT]int

	Coordinate [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int
}

var (
	gameC [constant.MAXROOMCOUNT]chan int
	// tetris的地圖
	tetrisMap [][][]string

	// tetris 運轉規則
	tetrisRule [constant.TETRISTYPELENGTH][constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int

	// 接收玩家指令的chan
	playerACh []chan int
	playerBCh []chan int

	//
	gameStartC chan string
	gameEndC   chan string

	// 正在移動的tetris 位置
	tetrisNow [][]TetrisSite

	// 遊戲玩家列表
	playerList [][]string

	tetrisIndex [][][]int

	tetrisStore [][]TetrisSite

	tetrisList = []string{"I", "J", "L", "O", "S", "T", "Z"}

	// 每個tetris 初始位置  [tetrisIndex][X,Y][中心點順時針順序]
	tetrisStartSite [constant.TETRISTYPELENGTH][constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int

	tetrisMidInitialPointX = constant.TETRISWIDTH / 2
	tetrisMidInitialPointY = constant.TETRISLENGTH
)

func Initialize() {
	roomLen := len(redis.RoomList)
	playerACh = make([]chan int, 0)
	playerBCh = make([]chan int, 0)

	InitializeTetrisRule(&tetrisRule)

	for i := 0; i < roomLen; i++ {
		playerACh = append(playerACh, make(chan int))
		playerBCh = append(playerBCh, make(chan int))

		tetrisMap = append(tetrisMap, InitializeMap())
	}
	/*
		for i := range gameC{
			GameServer(gameC[i])
		}*/

}

// initialize the tetris rotate rule
func InitializeTetrisRule(tetrisRotate *[constant.TETRISTYPELENGTH][constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int) {

	for tetrisType := 0; tetrisType < len(tetrisList); tetrisType++ {
		switch tetrisType {
		case constant.TETRIS_I:
			tetrisStartSite[constant.TETRIS_I] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX}, {tetrisMidInitialPointY, tetrisMidInitialPointY - 1, tetrisMidInitialPointY - 2, tetrisMidInitialPointY - 3},
			}
			tetrisRotate[constant.TETRIS_I] = [constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int{
				//向右旋轉 直的最下面為 index [0]
				{
					{0, 1, 2, 3}, {0, -1, -2, -3},
				},
				{
					{0, -1, -2, -3}, {0, 1, 2, 3},
				},
			}
		case constant.TETRIS_J:
			tetrisStartSite[constant.TETRIS_J] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX - 1, tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX}, {tetrisMidInitialPointY, tetrisMidInitialPointY - 1, tetrisMidInitialPointY - 2, tetrisMidInitialPointY - 2},
			}

			tetrisRotate[constant.TETRIS_J] = [constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int{
				//向右旋轉 J的最左邊為 index [0]
				{
					{0, -1, 0, 1}, {1, 0, -1, -2},
				},
				{
					{1, 0, -1, -2}, {1, 2, 1, 0},
				},
				{
					{1, 2, 1, 0}, {-2, -1, 0, 1},
				},
				{
					{-2, -1, 0, 1}, {0, -1, 0, 1},
				},
			}
		case constant.TETRIS_L:
			tetrisStartSite[constant.TETRIS_L] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX + 1}, {tetrisMidInitialPointY, tetrisMidInitialPointY - 1, tetrisMidInitialPointY - 2, tetrisMidInitialPointY - 2},
			}
			tetrisRotate[constant.TETRIS_J] = [constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int{
				//向右旋轉 L的最右邊為 index [0]
				{
					{-1, 0, 1, 2}, {0, 1, 0, -1},
				},
				{
					{0, 1, 0, -1}, {2, 1, 0, -1},
				},
				{
					{2, 1, 0, -1}, {-1, -2, -1, 0},
				},
				{
					{-1, -2, -1, 0}, {-1, 0, 1, 2},
				},
			}
		case constant.TETRIS_O:
			tetrisStartSite[constant.TETRIS_O] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX, tetrisMidInitialPointX + 1, tetrisMidInitialPointX, tetrisMidInitialPointX + 1}, {tetrisMidInitialPointY - 1, tetrisMidInitialPointY - 1, tetrisMidInitialPointY, tetrisMidInitialPointY},
			}

		case constant.TETRIS_S:
			tetrisStartSite[constant.TETRIS_S] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX, tetrisMidInitialPointX + 1, tetrisMidInitialPointX + 1, tetrisMidInitialPointX + 2}, {tetrisMidInitialPointY - 1, tetrisMidInitialPointY, tetrisMidInitialPointY, tetrisMidInitialPointY + 1},
			}

			tetrisRotate[constant.TETRIS_S] = [constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int{
				//向右旋轉 S的最左下為 index [0]
				{
					{1, 0, -1, -2}, {0, 1, 0, 1},
				},
				{
					{-1, 0, 1, 2}, {0, -1, 0, -1},
				},
			}
		case constant.TETRIS_T:
			tetrisStartSite[constant.TETRIS_I] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX - 1, tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX + 1}, {tetrisMidInitialPointY, tetrisMidInitialPointY, tetrisMidInitialPointY, tetrisMidInitialPointY - 1},
			}
			tetrisRotate[constant.TETRIS_T] = [constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int{
				//向右旋轉 直的最下面為 index [0]
				{
					{1, 0, -1, -1}, {1, 0, -1, 1},
				},
				{
					{1, 0, -1, 1}, {-2, -1, 0, 0},
				},
				{
					{-2, -1, 0, 0}, {0, 1, 2, 0},
				},
				{
					{0, 1, 2, 0}, {1, 0, -1, -1},
				},
			}
		case constant.TETRIS_Z:
			tetrisStartSite[constant.TETRIS_I] = [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int{
				{tetrisMidInitialPointX - 1, tetrisMidInitialPointX, tetrisMidInitialPointX, tetrisMidInitialPointX + 1}, {tetrisMidInitialPointY, tetrisMidInitialPointY, tetrisMidInitialPointY - 1, tetrisMidInitialPointY - 1},
			}

			tetrisRotate[constant.TETRIS_Z] = [constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int{
				//向右旋轉 直的最下面為 index [0]
				{
					{0, -1, 0, -1}, {-1, 0, 1, 2},
				},
				{
					{0, 1, 0, 1}, {1, 0, -1, -2},
				},
			}
		}
	}

	return
}

func InitializeMap() [][]string {
	var emptyMap [][]string

	for i := 0; i < constant.TETRISWIDTH; i++ {
		mapW := []string{}
		for j := 0; j < constant.TETRISLENGTH; j++ {
			mapW = append(mapW, "empty")
		}
		emptyMap = append(emptyMap, mapW)
	}
	return emptyMap
}

/*
func GameServer(gameChannel chan int) {
	for {
		gameStart := <-gameStartC
		gameEnd := <-gameEndC
		go func() {
			for {
				log.Println(gameStart, "start game")
				TetrisMoving(gameStart)
				time.Sleep(constant.TETRIS_FALL_WAITING * time.Second)
				if gameEnd == gameStart {
					break
				}
			}
		}()
	}
}*/

func Commond(m *melody.Melody, s *melody.Session, msg []byte) error {
	return nil
}

func StartGame(room string) error {
	gameStartC <- room
	return nil
}

func EndGame(room, winner string) error {
	gameEndC <- room
	redis.SetRoomPlayerState(room, constant.ROOMSTATEWAITING)
	redis.UpdatePlayerWin(winner)
	return nil
}

func RestartGame(room string) error {
	gameStartC <- room
	return nil
}

func TetrisMoving(room string) error {
	roomIndex, _, _, err := GetIndex(room, constant.REQUEST_EMPTY_TYPE_STRING)
	if err != nil {
		return err
	}

	for _, id := range playerList[roomIndex] {
		TetrisFall(room, id)
	}
	return nil
}

// [Summary] Set tetris fall
func TetrisFall(room string, id string) (bool, error) {

	roomIndex, idIndex, tetrisIndex, err := GetIndex(room, id)
	if err != nil {
		return false, err
	}

	crashCheck, err := CrashCheck(roomIndex, idIndex, tetrisIndex, constant.TETRIS_CRASH_TYPE_FALL)
	if err != nil {
		return false, err
	}
	if crashCheck {
		for index := range tetrisNow[roomIndex][idIndex].Coordinate[constant.TETRIS_COORDINATE_Y] {
			tetrisNow[roomIndex][idIndex].Coordinate[constant.TETRIS_COORDINATE_Y][index] -= constant.TETRIS_MOVE_SPEED
		}
	} else {
		return false, err
	}

	return true, nil
}

func TetrisRotate(room string, id string) (bool, error) {

	roomIndex, idIndex, tetrisIndex, err := GetIndex(room, id)
	if err != nil {
		return false, err
	}

	crashCheck, err := CrashCheck(roomIndex, idIndex, tetrisIndex, constant.TETRIS_CRASH_TYPE_FALL)
	if err != nil {
		return false, err
	}
	if crashCheck {
		for coordinate, rotateRule := range tetrisRule[tetrisIndex][tetrisNow[roomIndex][idIndex].TetrisRotateType] {
			for index, rule := range rotateRule {
				tetrisNow[roomIndex][idIndex].Coordinate[coordinate][index] += rule
			}
		}

		tetrisNow[roomIndex][idIndex].TetrisRotateType = (tetrisNow[roomIndex][idIndex].TetrisRotateType + 1) % 4
	} else {
		return false, nil
	}

	return true, nil
}

func CreateNewTetris(room string, id string) error {
	var (
		tetrisN        string
		newTetrixIndex = rand.Intn(constant.TETRISTYPELENGTH - 1)
	)

	roomIndex, player, _, err := GetIndex(room, id)
	if err != nil {
		return err
	}

	tetrisIndex[roomIndex][player][newTetrixIndex] += 1
	tetrisN = strconv.Itoa(tetrisIndex[roomIndex][player][newTetrixIndex])

	tetrisNow[roomIndex][player] = TetrisSite{
		TetrisName:       tetrisList[newTetrixIndex] + tetrisN,
		TetrisType:       tetrisList[newTetrixIndex],
		TetrisRotateType: constant.TETRIS_ROTATE_INITIAL_TYPE,
		Coordinate:       tetrisStartSite[newTetrixIndex],
	}

	return nil
}

// 碰撞檢查
func CrashCheck(roomIndex int, playerIndex int, tetrisIndex int, crashType int) (bool, error) {
	var monitorTetris TetrisSite = tetrisNow[roomIndex][playerIndex]

	switch crashType {
	case constant.TETRIS_CRASH_TYPE_ROTATE:
		for coordinate, rotateRule := range tetrisRule[tetrisIndex][tetrisNow[roomIndex][playerIndex].TetrisRotateType] {
			for index, rule := range rotateRule {
				monitorTetris.Coordinate[coordinate][index] += rule
			}
		}

		for i := 0; i < constant.TETRIS_COUNT; i++ {
			if tetrisMap[roomIndex][monitorTetris.Coordinate[constant.TETRIS_COORDINATE_X][i]][monitorTetris.Coordinate[constant.TETRIS_COORDINATE_Y][i]] != constant.TETRIS_MAP_EMPTY {
				return false, nil
			}
		}
	case constant.TETRIS_CRASH_TYPE_FALL:
		for index := range monitorTetris.Coordinate[constant.TETRIS_COORDINATE_Y] {
			monitorTetris.Coordinate[constant.TETRIS_COORDINATE_Y][index] -= constant.TETRIS_FALL_SPEED
		}

		for i := range monitorTetris.Coordinate[constant.TETRIS_COORDINATE_Y] {
			if tetrisMap[roomIndex][monitorTetris.Coordinate[constant.TETRIS_COORDINATE_X][i]][monitorTetris.Coordinate[constant.TETRIS_COORDINATE_Y][i]] != constant.TETRIS_MAP_EMPTY {
				return false, nil
			}
		}

	default:
		return false, errors.New("get crash check failed!")
	}
	return true, nil
}

// 暫存方塊
func TetrisStore(room string, id string) (bool, error) {
	roomIndex, playerIndex, _, err := GetIndex(room, id)
	if err != nil {
		return false, err
	}

	tetrisStore[roomIndex][playerIndex] = tetrisNow[roomIndex][playerIndex]
	tetrisStore[roomIndex][playerIndex].TetrisRotateType = constant.TETRIS_ROTATE_INITIAL_TYPE

	err = CreateNewTetris(room, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// 找到房名和ID對應的index
func GetIndex(room, id string) (int, int, int, error) {
	var (
		roomIndex   int
		playerIndex int
		tetrisIndex int
	)

	if room == constant.REQUEST_EMPTY_TYPE_STRING && id == constant.REQUEST_EMPTY_TYPE_STRING {
		return constant.RESPONSE_ERROR_TYPE_INT, constant.RESPONSE_ERROR_TYPE_INT, constant.RESPONSE_ERROR_TYPE_INT, errors.New("get index failed!")
	}

	if room != constant.REQUEST_EMPTY_TYPE_STRING || id != constant.REQUEST_EMPTY_TYPE_STRING {
		for i, roomName := range redis.RoomList {
			if roomName == room {
				roomIndex = i
			}
		}
		if roomIndex != constant.RESPONSE_ERROR_TYPE_INT {
			for i, playerID := range playerList[roomIndex] {
				if playerID == id {
					playerIndex = i
				}
			}
		}
	}

	if roomIndex != constant.RESPONSE_ERROR_TYPE_INT && playerIndex != constant.RESPONSE_ERROR_TYPE_INT {
		tetris := tetrisNow[roomIndex][playerIndex].TetrisType
		for i, tetrisName := range tetrisList {
			if tetris == tetrisName {
				tetrisIndex = i
			}
		}
		if tetrisIndex == constant.RESPONSE_ERROR_TYPE_INT {
			return constant.RESPONSE_ERROR_TYPE_INT, constant.RESPONSE_ERROR_TYPE_INT, constant.RESPONSE_ERROR_TYPE_INT, errors.New("get tetris index failed!")
		}

		if tetrisIndex == constant.RESPONSE_ERROR_TYPE_INT {
			return constant.RESPONSE_ERROR_TYPE_INT, constant.RESPONSE_ERROR_TYPE_INT, tetrisIndex, nil
		}
	}

	return roomIndex, playerIndex, tetrisIndex, nil
}
