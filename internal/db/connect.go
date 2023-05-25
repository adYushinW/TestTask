package db

import (
	"database/sql"
	"fmt"

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
	return fmt.Sprintf(
		`host=%s port=%s user=%s password=%s database=%s sslmode=%s`,
		cdb.Host, cdb.Port, cdb.User, cdb.Password, cdb.Database, cdb.SSLmode,
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

func newConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", constDB().String())
	if err != nil {
		return nil, err
	}
	fmt.Println(constDB().String())
	if err := db.Ping(); err != nil {
		fmt.Println("Пинг не прошёл")
		return nil, err
	}

	return db, nil
}
