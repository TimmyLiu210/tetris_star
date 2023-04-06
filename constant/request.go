package constant
// player event type
const (
	SIGN_UP = iota + 100
	SIGN_IN
	SIGN_OUT

	IN_ROOM
	OUT_ROOM

	START_GAME
	END_GAME

	GAME_COMMOND
)

// broadcast event type
const (
	ROOM_PLAYER_CHANGE = iota + 300
	ROOM_WAITING
	ROOM_START_GAME
)
/*

{
    "event_type": 101,
    "data": {
    "account": "123456",
    "password": "123456"
    }
}


{
	"event_type": 103,
	"Room": "room-1"
}

*/

const (
	REQUEST_EMPTY_TYPE_STRING = ""
	REQUEST_EMPTY_TYPE_INT = 0
)