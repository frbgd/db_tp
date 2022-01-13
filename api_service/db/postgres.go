package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

var PgEngine *PostgresDbEngine

type PostgresDbEngine struct {
	connStr        string
	connectionPool *pgxpool.Pool
}

func NewPostgresDbEngine(user string, password string, database string) (*PostgresDbEngine, error) {
	p := new(PostgresDbEngine)
	p.connStr = fmt.Sprintf("host=/var/run/postgresql dbname=%s user=%s password=%s sslmode=disable pool_max_conns=5000",
		database,
		user,
		password,
	)
	config, err := pgxpool.ParseConfig(p.connStr)
	if err != nil {
		return nil, err
	}
	p.connectionPool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (engine *PostgresDbEngine) Close() {
	engine.connectionPool.Close()
}
