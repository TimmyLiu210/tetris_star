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
}

var (
	// tetris的地圖
	tetrisMap [][][]string

	// tetris 運轉規則
	tetrisRule [constant.TETRISTYPELENGTH][constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int

	// 接收玩家指令的chan
	playerACh []chan int
	playerBCh []chan int

	// 正在移動的tetris 位置
	tetrisNow [][]TetrisSite

	// 遊戲玩家列表
	playerList [][]string

	tetrisIndex [][][]int

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

	go func() {
		GameServer()
	}()

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

func GameServer() error {

	return nil
}

func Commond(m *melody.Melody, s *melody.Session, msg []byte) error {
	return nil
}

func StartGame() error {
	return nil
}

func EndGame() error {
	return nil
}

func TetrisFall(roomIndex int) error {

	return nil
}

func TetrisRotate(roomIndex int, id string) error {
	player, err := GetPlayerIndex(roomIndex, id)
	if err != nil {
		return err
	}

	tetris, err := getTetrisIndex(tetrisNow[roomIndex][player].TetrisType)
	if err != nil {
		return err
	}

	for coordinate, rotateRule := range tetrisRule[tetris][tetrisNow[roomIndex][player].TetrisRotateType] {
		switch coordinate {
		case constant.TETRIS_COORDINATE_X:
			for index, rule := range rotateRule {
				tetrisNow[roomIndex][player].X[index] += rule
			}
		case constant.TETRIS_COORDINATE_Y:
			{
				for index, rule := range rotateRule {
					tetrisNow[roomIndex][player].Y[index] += rule
				}
			}
		}
	}

	tetrisNow[roomIndex][player].TetrisRotateType = tetrisNow[roomIndex][player].TetrisRotateType % 4

	return nil
}

func CreateNewTetris(roomIndex int, id string) error {

	var (
		tetrisN string
	)
	newTetrixIndex := rand.Intn(6)

	player, err := GetPlayerIndex(roomIndex, id)
	if err != nil {
		return err
	}

	tetrisIndex[roomIndex][player][newTetrixIndex] += 1
	tetrisN = strconv.Itoa(tetrisIndex[roomIndex][player][newTetrixIndex])

	tetrisNow[roomIndex][player] = TetrisSite{
		TetrisName:       tetrisList[newTetrixIndex] + tetrisN,
		TetrisType:       tetrisList[newTetrixIndex],
		TetrisRotateType: constant.TETRIS_ROTATE_INITIAL_TYPE,
		X:                tetrisStartSite[newTetrixIndex][constant.TETRIS_COORDINATE_X],
		Y:                tetrisStartSite[newTetrixIndex][constant.TETRIS_COORDINATE_Y],
	}

	return nil
}

func GetPlayerIndex(roomIndex int, id string) (int, error) {
	for playerIndex, playerID := range playerList[roomIndex] {
		if playerID == id {
			return playerIndex, nil
		}
	}
	return 0, errors.New("get player index failed!")
}

func getTetrisIndex(name string) (int, error) {
	for tetrisIndex, tetrisName := range tetrisList {
		if tetrisName == name {
			return tetrisIndex, nil
		}
	}
	return 0, errors.New("get tetris name index failed!")
}
