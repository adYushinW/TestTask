package gin

import (
	"net/http"
	"strconv"

	"github.com/adYushinW/TestTask/internal/app"
	"github.com/gin-gonic/gin"
)

func Service(app *app.App) error {
	r := gin.Default()

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

		attribute := c.Query("attribute")
		order := c.Query("order")
		limit := c.Query("limit")
		offset := c.Query("offset")

		catsInfo, err := app.GetCats(attribute, order, limit, offset)
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
