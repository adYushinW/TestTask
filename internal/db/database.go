package db

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/adYushinW/TestTask/internal/model"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	GetCats(attribute string, order string, limit uint64, offset uint64) ([]*model.Cats, error)
	AddCat(name string, color string, tail_length uint64, whiskers_length uint64) ([]*model.Cats, error)
	CatColor() ([]*model.Cat_colors_info, error)
	CatsInfo() ([]*model.Cats_stat, error)
}

type database struct {
	conn *pgx.Conn
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

func (db *database) GetCats(attribute string, order string, limit uint64, offset uint64) ([]*model.Cats, error) {
	qb := sq.Select("name", "color", "tail_length", "whiskers_length").
		From("cats").
		PlaceholderFormat(sq.Dollar)

	if attribute != "" {
		if strings.ToLower(order) == "desc" {
			qb = qb.OrderByClause("? DESC", attribute)
		} else {
			qb = qb.OrderByClause("?", attribute)
		}
	}

	if limit > 0 {
		qb = qb.Limit(limit)
	}

	if offset > 0 {
		qb = qb.Offset(offset)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.conn.Query(context.Background(), sql, args...)
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

func (db *database) AddCat(name string, color string, tail_length uint64, whiskers_length uint64) ([]*model.Cats, error) {

	//query := "INSERT INTO cats (name, color, tail_length, whiskers_length) VALUES ($1, $2, $3, $4) RETURNING name, color, tail_length, whiskers_length"

	qb := sq.Insert("cats").
		Columns("name", "color", "tail_length", "whiskers_length").
		Values(name, color, tail_length, whiskers_length).
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (name) DO NOTHING").
		Suffix("RETURNING name, color, tail_length, whiskers_length")

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	row, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	defer row.Close()
	result := make([]*model.Cats, 0)

	for row.Next() {
		cat := new(model.Cats)

		if err := row.Scan(&cat.Name, &cat.Color, &cat.Tail_length, &cat.Whiskers_length); err != nil {
			continue
		}

		result = append(result, cat)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *database) CatColor() ([]*model.Cat_colors_info, error) {
	sb := sq.Select("color", "COUNT(1)").From("cats").GroupBy("color")

	ib := sq.Insert("cat_colors_info").
		Columns("color", "count").
		Select(sb).
		Suffix("ON CONFLICT (color) DO UPDATE SET count = EXCLUDED.count ").
		Suffix("RETURNING color, count")

	sql, args, err := ib.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.conn.Query(context.TODO(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Cat_colors_info, 0)

	for rows.Next() {
		cat := new(model.Cat_colors_info)
		if err := rows.Scan(&cat.Color, &cat.Count); err != nil {
			continue
		}

		result = append(result, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *database) CatsInfo() ([]*model.Cats_stat, error) {

	query := `INSERT INTO cats_stat (tail_length_mean, tail_length_median, tail_length_mode, 
				whiskers_length_mean, whiskers_length_median,  whiskers_length_mode)
				SELECT
				  AVG(cats. tail_length) AS tail_length_mean,
				  PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY cats.tail_length) AS tail_length_median,
				  (
					SELECT array_agg(tail_length)
					FROM (
					 SELECT tail_length FROM cats GROUP BY tail_length HAVING tail_length = max(tail_length) ORDER BY tail_length ASC
					) as t1
				   ) as tail_length_mode,
				  AVG(cats.whiskers_length) AS whiskers_length_mean,
				  PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY cats.whiskers_length) AS whiskers_length_median,
				  (
				   SELECT array_agg(whiskers_length)
				   FROM (
				    SELECT whiskers_length FROM cats GROUP BY whiskers_length HAVING whiskers_length = max(whiskers_length) ORDER BY whiskers_length ASC
				   ) as t1
				  ) as whiskers_length_mode
				FROM cats
				RETURNING tail_length_mean, tail_length_median, tail_length_mode, whiskers_length_mean, whiskers_length_median, whiskers_length_mode`

	row, err := db.conn.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}
	defer row.Close()

	result := make([]*model.Cats_stat, 0)

	for row.Next() {
		cat := new(model.Cats_stat)
		if err := row.Scan(&cat.Tail_length_mean, &cat.Tail_length_median, &cat.Tail_length_mode,
			&cat.Whiskers_length_mean, &cat.Whiskers_length_median, &cat.Whiskers_length_mode); err != nil {
			continue
		}

		result = append(result, cat)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
