package gin

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/adYushinW/TestTask/internal/app"
	"github.com/gin-gonic/gin"
)

var count uint64

func Service(app *app.App) error {
	go func() {
		t := time.NewTicker(time.Minute)

		for range t.C {
			atomic.StoreUint64(&count, 0)
			t.Reset(time.Minute)
		}
	}()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if count > 2 {
			c.JSON(http.StatusTooManyRequests, "Too Many Requests")
			return
		}

		atomic.AddUint64(&count, 1)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Cats Service. Version 0.1")
	})

	r.GET("/count", func(c *gin.Context) {
		c.JSON(http.StatusOK, count)
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
		var l, o uint64

		attribute := c.Query("attribute")
		order := c.Query("order")
		limit := c.Query("limit")
		offset := c.Query("offset")

		if limit != "" {
			l, err = strconv.ParseUint(limit, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, "Bad Request")
				return
			}
		}

		if limit != "" {
			o, err = strconv.ParseUint(offset, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, "Bad Request")
				return
			}
		}

		catsInfo, err := app.GetCats(attribute, order, l, o)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		c.JSON(http.StatusOK, catsInfo)
	})

	r.POST("/cat", func(c *gin.Context) {

		name := c.Query("name")
		color := c.Query("color")

		tail_length, err := strconv.ParseUint(c.Query("tail_length"), 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, "tail_lenth only numeric > 0")
			return
		}
		tl8 := uint8(tail_length)

		whiskers_length, err := strconv.ParseUint(c.Query("whiskers_length"), 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, "tail_lenth only numeric > 0")
			return
		}
		wl8 := uint8(whiskers_length)

		newCat, err := app.AddCat(name, color, tl8, wl8)
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
