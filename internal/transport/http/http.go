package http

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/adYushinW/TestTask/internal/app"
	"github.com/adYushinW/TestTask/internal/model"
	"github.com/gin-gonic/gin"
)

var count uint64

var mu sync.Mutex

func Service(app *app.App) error {

	go func() {
		t := time.NewTicker(time.Minute)
		for range t.C {
			mu.Lock()
			count = 0
			mu.Unlock()
			t.Reset(time.Minute)
		}
	}()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		func() {
			mu.Lock()
			defer mu.Unlock()

			count++

			fmt.Println(count)
			if count > 2 {

				c.JSON(http.StatusTooManyRequests, "Too Many Requests")
				c.Abort()
				return
			}

		}()
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Cats Service. Version 0.1")
	})

	r.GET("/cat_color", func(c *gin.Context) {
		catColor, err := app.CatColor()
		if err != nil {
			c.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		c.JSON(http.StatusOK, catColor)
	})

	r.GET("/cats", func(c *gin.Context) {
		var err error
		var lim, off uint64

		attribute := c.Query("attribute")
		order := c.Query("order")
		limit := c.Query("limit")
		offset := c.Query("offset")

		if limit != "" {
			lim, err = strconv.ParseUint(limit, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, "Wrong Limit")
				return
			}
		}

		if offset != "" {
			off, err = strconv.ParseUint(offset, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, "Wrong Offset")
				return
			}
		}

		catsInfo, err := app.GetCats(attribute, order, lim, off)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		c.JSON(http.StatusOK, catsInfo)
	})

	r.POST("/cat", func(c *gin.Context) {
		var err error

		cat := model.Cats{}

		if err := c.BindJSON(&cat); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		newCat, err := app.AddCat(cat.Name, cat.Color, cat.Tail_length, cat.Whiskers_length)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		c.JSON(http.StatusOK, newCat)
	})

	r.GET("/cats_stat", func(c *gin.Context) {

		catsInfo, err := app.CatsInfo()
		if err != nil {
			c.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		c.JSON(http.StatusOK, catsInfo)
	})

	return r.Run(":8080")
}
