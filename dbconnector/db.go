/*
Package db create a postgres database connection pool.
*/
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool is the pointer to the db connection pool.
var DBPool *pgxpool.Pool

// InitDBConn returns a postgresdb connection pool
func InitDBConn(ctx context.Context, host string, port string, user string, password string, database string) error {

	DBPool = nil

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return fmt.Errorf("error creating db pool: %w", err)
	}

	err = PingDB(ctx, dbpool)
	if err != nil {
		return fmt.Errorf("error creating db pool: %w", err)
	}

	DBPool = dbpool

	return nil
}

// PingDB checks the db connection.
func PingDB(ctx context.Context, dbpool *pgxpool.Pool) error {
	err := dbpool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}

	return nil
}
