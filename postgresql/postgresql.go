package postgresql

import (
	"context"
	"fmt"
	"net/url"
	"tetris/config"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	PoolWr PoolWrapper
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PoolWrapper struct {
	master *pgxpool.Pool
	slave  *pgxpool.Pool
}

func (p *PoolWrapper) Write() *pgxpool.Pool {
	return p.master
}

func (p *PoolWrapper) Read() *pgxpool.Pool {
	if p.slave == nil {
		return p.master
	}
	return p.slave
}

func Initialize() {
	masterUri := fmt.Sprintf("postgresql://%s", config.PostgresqlMaster)
	PoolWr.master = getPool(masterUri, "read-write")

}

func getPool(uri, mode string) *pgxpool.Pool {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Add("target_session_attrs", mode)
	q.Add("connect_timeout", "10")
	q.Add("pool_max_conns", "50")
	q.Add("pool_max_conn_lifetime", "180s")
	q.Add("pool_max_conn_idle_time", "180s")
	u.RawQuery = q.Encode()
	uri = u.String()

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		panic(err)
	}

	dbName := config.PostgresDbName
	cfg.ConnConfig.Database = dbName
	cfg.ConnConfig.RuntimeParams = map[string]string{
		"application_name":                    "tetris",
		"statement_timeout":                   "30000", // ms
		"idle_in_transaction_session_timeout": "30000", // ms
	}

	user := config.PostgresUser
	password := config.PostgresPassword
	if user != "" && password != "" {
		cfg.ConnConfig.User = user
		cfg.ConnConfig.Password = password
	}
	return setupPool(cfg)
}

func setupPool(cfg *pgxpool.Config) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	return pool
}