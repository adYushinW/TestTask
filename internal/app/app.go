package app

import (
	"github.com/adYushinW/TestTask/internal/db"
	"github.com/adYushinW/TestTask/internal/model"
)

type App struct {
	db db.Database
}

func New(db db.Database) *App {
	return &App{
		db: db,
	}
}

func (a *App) GetCats(attribute string, order string, limit string, offset string) ([]*model.Cats, error) {
	return a.db.GetCats(attribute, order, limit, offset)
}

func (a *App) AddCat(name string, color string, tail_length uint8, whiskers_length uint8) ([]*model.Cats, error) {
	return a.db.AddCat(name, color, tail_length, whiskers_length)
}

func (a *App) CatColor() ([]*model.Cat_colors_info, error) {
	return a.db.CatColor()
}

func (a *App) CatsInfo() ([]*model.Cats_stat, error) {
	return a.db.CatsInfo()
}
