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