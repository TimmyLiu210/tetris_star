package constant

const (
	PLAYERPREFIX = "player-"
	TETRISPREFIX = "tetris-"

	SESSIONPREFIX = "chat_id-"
	MESSAGEPREFIX = "message-"

	ROOMPREFIX = "room-"
)

const (
	PLACE     string = "place"
	PLACEHALL string = "hall"
)

// variables limit
const (
	MAXROOMCOUNT int = 5 //最大房間數量
	MINROOMCOUNT int = 1

	MAXACCOUNT  int = 8 //最大帳號字數
	MAXPASSWORD int = 6 //最大密碼字數
	MAXNICKNAME int = 8 //最大暱稱字數

	MAXPLAYERINFOLIST int64 = 6 //玩家列表總數
)

// player list
const (
	PLAYERID = iota
	PLAYERICON
	PLAYERACCOUNT
	PLAYERPASSWORD
	PLAYERNICKNAME
	PLAYERWIN

	NOWPLAYERPLACE
	PLAYERISROOMOWNER
)

// room list
const (
	ROOMDEX = iota
	ROOMOWNER
	ROOMPARTNER
)

// is room owner or not
const (
	OWNERSTATETRUE  string = "true"
	OWNERSTATEFALSE string = "false"
)

// room state
const (
	ROOMSTATEPLAYING string = "playing"
	ROOMSTATEWAITING string = "waiting"
)

// tetris const
const (
	TETRISWIDTH         = 10
	TETRISLENGTH        = 60
	TETRISTYPELENGTH    = 7
	TETRISCHANNELBUFFER = 6

	INITDEX    = 1
	INITTETRIS = 6

	PLAYERCOUNT = 2

	TETRIS_X           = 4
	TETRIS_Y           = 4
	TETRIS_ROTATE_TYPE = 4
	TETRIS_COUNT       = 4

	TETRIS_COORDINATE_COUNT = 2
)

// commond type
const (
	ROTATE = iota + 1
	FALL
	STORE
	MOVE
)

// tetris index
const (
	TETRIS_I = iota
	TETRIS_J
	TETRIS_L
	TETRIS_O
	TETRIS_S
	TETRIS_T
	TETRIS_Z
)

const (
	PLAYER_A = iota
	PLAYER_B
)

const (
	TETRIS_FALL_SPEED = 1
	TETRIS_MOVE_SPEED = 1

	TETRIS_FALL_WAITING = 1

	TETRIS_MAP_EMPTY = "empty"

	TETRIS_ROTATE_INITIAL_TYPE = 0
	TETRIS_CRASH_TYPE_ROTATE   = 0
	TETRIS_CRASH_TYPE_FALL     = 1

	TETRIS_COORDINATE_X = 0
	TETRIS_COORDINATE_Y = 1

	TETRISSTORECOUNT = 1

	TETRIS_WAITING_COUNT = 4

	PLAYER_ONE = 0
	PLAYER_TWO = 1
)

const (
	RESPONSE_ERROR_TYPE_STRING = ""
	RESPONSE_ERROR_TYPE_INT    = 0
)
