package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnector() *sql.DB {
	createPool()
	var db *sql.DB

	db, err := sql.Open("pgx", "host=localhost user=admin password=admin dbname=game_hangar sslmode=disable")
	if err != nil {
		panic("failed to connect database " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("failed to ping database " + err.Error())
	}
	return db
}

func createPool() {
	myConf := &pgxpool.Config{
		ConnConfig: &pgx.ConnConfig{
			Config: pgconn.Config{
				Host:     "localhost",
				Port:     5432,
				Database: "game_hangar",
				User:     "admin",
				Password: "admin",
			},
		},
	}
	print(myConf)

	// TODO : Check sslmode
	dns := "postgres://admin:admin@localhost:5432/game_hangar?sslmode=disable"
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		panic("error while parsing db config" + err.Error())
	}

	// 2. Production-ready configurations
	// Max connections: Use (CPU cores * 2) + 1 as a baseline
	config.MaxConns = 25
	// Min connections: Keeps warm connections ready for traffic spikes
	config.MinConns = 5
	// MaxConnLifetime: Prevents issues with stale connections or memory leaks
	config.MaxConnLifetime = 1 * time.Hour
	// MaxConnIdleTime: Closes connections that haven't been used recently
	config.MaxConnIdleTime = 30 * time.Minute
	// HealthCheckPeriod: How often the pool checks if connections are still alive
	config.HealthCheckPeriod = 1 * time.Minute
	// ConnectTimeout: Time limit for establishing the initial physical connection
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	// 3. Create the connection pool
	// Note: NewWithConfig does not immediately connect to the DB
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Sprintf("could not create connection pool: %w", err))
	}
	// 4. Ping the database to ensure the connection is actually working
	// We use a context with a timeout specifically for the ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("could not ping database: %w", err))
	}
	fmt.Printf("Connection poool created: %v\n", pool)
}
