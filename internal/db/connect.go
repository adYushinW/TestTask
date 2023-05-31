package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type ConfigDatabase struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLmode  string
}

func (cdb ConfigDatabase) String() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cdb.User, cdb.Password, cdb.Host, cdb.Port, cdb.Database, cdb.SSLmode,
	)
}

func constDB() ConfigDatabase {
	return ConfigDatabase{
		Host:     "172.21.0.2",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		SSLmode:  "disable",
	}
}

func newConnect() (*pgx.Conn, error) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, constDB().String())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		fmt.Println("Пинг не прошёл")
		return nil, err
	}

	return db, nil
}
