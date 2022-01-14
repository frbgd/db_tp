package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

var PgEngine *PostgresDbEngine

type PostgresDbEngine struct {
	connStr string
	CP      *pgxpool.Pool
}

func NewPostgresDbEngine(host string, port string, database string, user string, password string) (*PostgresDbEngine, error) {
	p := new(PostgresDbEngine)
	p.connStr = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable pool_max_conns=30",
		host,
		port,
		database,
		user,
		password,
	)
	config, err := pgxpool.ParseConfig(p.connStr)
	if err != nil {
		return nil, err
	}
	p.CP, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (engine *PostgresDbEngine) Close() {
	engine.CP.Close()
}
