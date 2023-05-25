package db

import (
	"database/sql"
	"fmt"

	"github.com/adYushinW/TestTask/internal/model"
)

type Database interface {
	GetCats(attribute string, order string, limit string, offset string) ([]*model.Cats, error)
	AddCat(name string, color string, tail_length uint8, whiskers_length uint8) ([]*model.Cats, error)
	CatColor() ([]*model.Cat_colors_info, error)
}

type database struct {
	conn *sql.DB
}

func New() (Database, error) {
	conn, err := newConnect()
	if err != nil {
		return nil, err
	}

	return &database{
		conn: conn,
	}, nil
}

func (db *database) GetCats(attribute string, order string, limit string, offset string) ([]*model.Cats, error) {
	query := "SELECT name, color, tail_length, whiskers_length FROM cats"

	if attribute != "" {
		query = fmt.Sprintf("%s ORDER BY %s", query, attribute)
	}

	if order != "" {
		query = fmt.Sprintf("%s %s", query, order)
	}

	if limit != "" {
		query = fmt.Sprintf("%s LIMIT %s", query, limit)
	}

	if offset != "" {
		query = fmt.Sprintf("%s OFFSET %s", query, offset)
	}

	rows, err := db.conn.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Cats, 0)

	for rows.Next() {
		cats := new(model.Cats)

		if err := rows.Scan(&cats.Name, &cats.Color, &cats.Tail_length, &cats.Whiskers_length); err != nil {
			continue
		}

		result = append(result, cats)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil

}

func (db *database) AddCat(name string, color string, tail_length uint8, whiskers_length uint8) ([]*model.Cats, error) {

	query := "INSERT INTO cats (name, color, tail_length, whiskers_length) VALUES ($1, $2, $3, $4) RETURNING name, color, tail_length, whiskers_length"

	row := db.conn.QueryRow(query, name, color, tail_length, whiskers_length)

	result := make([]*model.Cats, 0)

	cat := new(model.Cats)

	err := row.Scan(&cat.Name, &cat.Color, &cat.Tail_length, &cat.Whiskers_length)

	if err != nil {
		return result, err
	}

	result = append(result, cat)

	return result, err
}

func (db *database) CatColor() ([]*model.Cat_colors_info, error) {
	query := "SELECT color, Count(color) FROM cats GROUP BY color "

	rows, err := db.conn.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Cat_colors_info, 0)
	fmt.Println("я тут")
	for rows.Next() {
		cats := new(model.Cat_colors_info)
		if err := rows.Scan(&cats.Color, &cats.Count); err != nil {
			continue
		}

		query := "INSERT INTO cat_colors_info (color, count) VALUES ($1, $2) RETURNING color, count"
		row := db.conn.QueryRow(query, &cats.Color, &cats.Count)

		cat := new(model.Cat_colors_info)
		if err := row.Scan(&cat.Color, &cat.Count); err != nil {
			continue
		}

		result = append(result, cat)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
