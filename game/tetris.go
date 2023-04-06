package game

import (
	"tetris/constant"
	"tetris/redis"

	"gopkg.in/olahol/melody.v1"
)

type TetrisCommond struct {
	Player  string
	Commond int
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
	tetrisNow [][][]int

	// 遊戲玩家列表
	playerList [][]string
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

	for tetrisType := range tetrisRotate {
		switch tetrisType {
		case constant.TETRIS_I:
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
		case constant.TETRIS_S:
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
