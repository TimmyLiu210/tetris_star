package game

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"tetris/constant"
	"tetris/redis"
	"time"
)

type TetrisCommond struct {
	Player  string
	Commond int
}

type TetrisSite struct {
	TetrisName string

	TetrisType       string
	TetrisRotateType int

	Coordinate [constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int
}

type Server struct {
	Start chan bool
	End   chan bool

	playerCommond chan TetrisCommond

	tetrisMap [constant.PLAYERCOUNT][constant.TETRISLENGTH][constant.TETRISWIDTH]string

	tetrisNow [constant.PLAYERCOUNT]TetrisSite

	tetrisWaiting [constant.PLAYERCOUNT][constant.TETRIS_WAITING_COUNT]TetrisSite

	tetrisStore [constant.PLAYERCOUNT][constant.TETRISSTORECOUNT]TetrisSite

	playerList [constant.PLAYERCOUNT]string

	playerMap map[string]int

	tetrisIndexList [constant.PLAYERCOUNT][constant.TETRISTYPELENGTH]int
}

var (
	// tetris 運轉規則
	tetrisRule [constant.TETRISTYPELENGTH][constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int

	tetrisList = []string{"I", "J", "L", "O", "S", "T", "Z"}

	tetrisTypeMap = make(map[string]int)

	tetrisServer = make(map[string]*Server)

	// 每個tetris 初始位置  [tetrisIndex][X,Y][中心點順時針順序]
	tetrisStartSite [constant.TETRISTYPELENGTH][constant.TETRIS_COORDINATE_COUNT][constant.TETRIS_COUNT]int

	tetrisMidInitialPointX = constant.TETRISWIDTH / 2
	tetrisMidInitialPointY = constant.TETRISLENGTH
)

func Initialize() {
	InitializeTetrisRule(&tetrisRule)

	for _, roomName := range redis.RoomList {
		newServer := InitialServer()

		tetrisServer[roomName] = &newServer

		go runServer(roomName, tetrisServer[roomName])
	}

}

func InitialServer() Server {
	var newServer = Server{
		Start: make(chan bool),
		End:   make(chan bool),

		tetrisMap: [constant.PLAYERCOUNT][constant.TETRISLENGTH][constant.TETRISWIDTH]string{},

		playerCommond: make(chan TetrisCommond),

		playerMap: make(map[string]int),
	}

	ResetServer(&newServer)

	return newServer
}

func StartGame(room string) {
	tetrisServer[room].Start <- true
}

func EndGame(room string) {
	tetrisServer[room].End <- true
}

func ResetServer(tetrisS *Server) {
	for i := range tetrisS.tetrisMap {
		for j := range tetrisS.tetrisMap[i] {
			for k := range tetrisS.tetrisMap[i][j] {
				tetrisS.tetrisMap[i][j][k] = constant.TETRIS_MAP_EMPTY
				tetrisS.tetrisMap[i][j][k] = constant.TETRIS_MAP_EMPTY
			}
		}
	}

	for i := range tetrisS.tetrisIndexList {
		for j := range tetrisS.tetrisIndexList[i] {
			tetrisS.tetrisIndexList[i][j] = 0
		}
	}

	for i := range tetrisS.tetrisNow {
		tetrisS.tetrisNow[i] = CreateNewTetris(&tetrisS.tetrisIndexList[i])
	}

	for i := range tetrisS.tetrisWaiting {
		for j := range tetrisS.tetrisWaiting[i] {
			tetrisS.tetrisWaiting[i][j] = CreateNewTetris(&tetrisS.tetrisIndexList[i])
		}
	}

	for i := range tetrisS.playerList {
		tetrisS.playerMap[tetrisS.playerList[i]] = i
	}
	log.Println("[Server Msg] Reset Server Success...")
	return
}

func runServer(roomName string, tetrisS *Server) {
	log.Println("[Server Msg]", roomName, "Server is runing...")
	var endMoving = make(chan bool)
	for {
		select {
		case <-tetrisS.Start:
			// count down

			// initial the stage
			log.Printf("[Server Msg] Start reset Server[%s]", roomName)
			ResetServer(tetrisS)
			go func() {
				for {
					select {
					case <-endMoving:
						return
					default:
						// tetris moving + waiting
						TetrisMoving(&tetrisS.tetrisMap[constant.PLAYER_A], &tetrisS.tetrisNow[constant.PLAYER_A])
						TetrisMoving(&tetrisS.tetrisMap[constant.PLAYER_B], &tetrisS.tetrisNow[constant.PLAYER_B])
						time.Sleep(constant.TETRIS_FALL_SPEED * time.Second)
					}
				}

			}()
			go func() {
				for {
					select {
					// player commond
					case c := <-tetrisS.playerCommond:
						log.Println(c)
					case <-tetrisS.End:
						// game end
						endMoving <- true

						// broadcast the end and set the new win for winner

						return
					default:
					}
				}
			}()

		default:
			time.Sleep(100 * time.Second)
		}
	}
}

// initialize the tetris rotate rule
func InitializeTetrisRule(tetrisRotate *[constant.TETRISTYPELENGTH][constant.TETRIS_ROTATE_TYPE][constant.TETRIS_X][constant.TETRIS_Y]int) {

	for index := range tetrisList {
		tetrisTypeMap[tetrisList[index]] = index
	}

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

func TetrisMoving(tetrisMap *[constant.TETRISLENGTH][constant.TETRISWIDTH]string, t *TetrisSite) (bool, error) {
	crash, err := CrashCheck(*tetrisMap, *t, constant.TETRIS_CRASH_TYPE_FALL)
	if err != nil {
		return false, err
	}
	if crash {
		for i := range t.Coordinate[constant.TETRIS_COORDINATE_Y] {
			tetrisMap[constant.TETRIS_COORDINATE_Y][i] = constant.TETRIS_MAP_EMPTY
			t.Coordinate[constant.TETRIS_COORDINATE_Y][i] -= constant.TETRIS_MOVE_SPEED
			tetrisMap[constant.TETRIS_COORDINATE_Y][i] = t.TetrisName
		}
	} else {
		return false, nil
	}

	return true, nil
}

func TetrisRotate(tetrisMap *[constant.TETRISLENGTH][constant.TETRISWIDTH]string, t *TetrisSite) (bool, error) {
	crash, err := CrashCheck(*tetrisMap, *t, constant.TETRIS_CRASH_TYPE_ROTATE)
	if err != nil {
		return false, err
	}
	if crash {
		for coordinate, rotateRule := range tetrisRule[tetrisTypeMap[t.TetrisType]][t.TetrisRotateType] {
			for index, rule := range rotateRule {
				t.Coordinate[coordinate][index] += rule
			}
		}
		t.TetrisRotateType = (t.TetrisRotateType + 1) % 4

	} else {
		return false, nil
	}

	return true, nil
}

func CreateNewTetris(tetrisIndexList *[constant.TETRISTYPELENGTH]int) TetrisSite {
	var (
		tetrisDex      string
		newTetrixIndex = rand.Intn(constant.TETRISTYPELENGTH - 1)
	)

	tetrisIndexList[newTetrixIndex] += 1

	tetrisDex = strconv.Itoa(tetrisIndexList[newTetrixIndex])

	newTetris := TetrisSite{
		TetrisName:       tetrisList[newTetrixIndex] + tetrisDex,
		TetrisType:       tetrisList[newTetrixIndex],
		TetrisRotateType: constant.TETRIS_ROTATE_INITIAL_TYPE,
		Coordinate:       tetrisStartSite[newTetrixIndex],
	}

	return newTetris
}

// 碰撞檢查
func CrashCheck(tetrisM [constant.TETRISLENGTH][constant.TETRISWIDTH]string, t TetrisSite, crashType int) (bool, error) {

	switch crashType {
	case constant.TETRIS_CRASH_TYPE_ROTATE:
		for coordinate, rotateRule := range tetrisRule[tetrisTypeMap[t.TetrisType]][t.TetrisRotateType] {
			for index, rule := range rotateRule {
				t.Coordinate[coordinate][index] += rule
			}
		}

		for i := 0; i < constant.TETRIS_COUNT; i++ {
			if tetrisM[t.Coordinate[constant.TETRIS_COORDINATE_X][i]][t.Coordinate[constant.TETRIS_COORDINATE_Y][i]] != constant.TETRIS_MAP_EMPTY {
				return false, nil
			}
		}
	case constant.TETRIS_CRASH_TYPE_FALL:
		for index := range t.Coordinate[constant.TETRIS_COORDINATE_Y] {
			t.Coordinate[constant.TETRIS_COORDINATE_Y][index] -= constant.TETRIS_FALL_SPEED
		}

		for i := range t.Coordinate[constant.TETRIS_COORDINATE_Y] {
			if tetrisM[t.Coordinate[constant.TETRIS_COORDINATE_X][i]][t.Coordinate[constant.TETRIS_COORDINATE_Y][i]] != constant.TETRIS_MAP_EMPTY {
				return false, nil
			}
		}

	default:
		return false, errors.New("get crash check failed!")
	}
	return true, nil
}

// 暫存方塊
func TetrisStore(tS *TetrisSite, t *TetrisSite, tetrisIndexList *[constant.TETRISTYPELENGTH]int) (bool, error) {
	tS = t
	tS.TetrisRotateType = constant.TETRIS_ROTATE_INITIAL_TYPE

	*t = CreateNewTetris(tetrisIndexList)

	return true, nil
}
