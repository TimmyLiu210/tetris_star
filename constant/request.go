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
    "account": "a1094137",
    "password": "553525"
    }
}
*/