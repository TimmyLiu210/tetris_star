package postgresql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
)

// player info
type PlayerInfo struct {
	PlayerID string `json:"player_id"`
	Account  string `json:"account"`
	Password string `json:"password"`
	NickName string `json:"nickname"`
	Icon     int    `json:"icon"`
	Win      int    `json:"win"`
}

// room info include playing state and player info in room
type RoomState struct {
	IsPlaying  bool     `json:"is_playing"`
	PlayerList []string `json:"player_list"`
}

// insert new player info to database
func (p *PlayerInfo) Register(ctx context.Context, db DB) error {
	cols := []string{
		"player_id",
		"account",
		"password",
		"nickname",
		"icon",
		"win",
	}

	setMap := map[string]interface{}{}
	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonData, &setMap); err != nil {
		return err
	}
	index := 1
	values := []string{}
	args := []interface{}{}
	for _, col := range cols {
		values = append(values, fmt.Sprintf("$%d", index))
		args = append(args, setMap[col])
		index++
	}

	query := fmt.Sprintf(`
		INSERT INTO players (%s)
		VALUES (%s)
	`,
		strings.Join(cols, ","),
		strings.Join(values, ","),
	)
	_, err = db.Exec(ctx, query, args...)

	return err
}

// update the player info in database
func (p *PlayerInfo) Upsert(ctx context.Context, db DB) error {
	cols := []string{
		"player_id",
		"account",
		"password",
		"nickname",
		"icon",
		"win",
	}

	exCols := []string{}
	for _, col := range cols {
		exCols = append(exCols, "EXCLUDED."+col)
	}

	setMap := map[string]interface{}{}
	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonData, &setMap); err != nil {
		return err
	}
	index := 1
	values := []string{}
	args := []interface{}{}
	for _, col := range cols {
		values = append(values, fmt.Sprintf("$%d", index))
		args = append(args, setMap[col])
		index++
	}

	query := fmt.Sprintf(`
		INSERT INTO players (%s)
		VALUES (%s)
		ON CONFLICT (account)
		DO UPDATE SET (%s) = (%s)
	`,
		strings.Join(cols, ","),
		strings.Join(values, ","),
		strings.Join(cols, ","),
		strings.Join(exCols, ","),
	)
	_, err = db.Exec(ctx, query, args...)

	return err
}

// get player info from database
func (p *PlayerInfo) GetPlayer(ctx context.Context, db DB) (*PlayerInfo, error) {
	query := `
		SELECT
		json_build_object(
			'player_id', player_id,
			'icon', icon,
			'account', account,
			'password', password,
			'nickname', nickname,
			'win', win
		) AS player
		FROM
			players
		WHERE
			account = $1

	`

	rows, err := db.Query(ctx, query, p.Account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	password := p.Password
	var dbPlayer *PlayerInfo
	if err := pgxscan.ScanOne(&dbPlayer, rows); err != nil {
		return nil, err
	}

	if password != p.Password {
		return nil, errors.New("Password incurrect!")
	}
	return dbPlayer, err
}
