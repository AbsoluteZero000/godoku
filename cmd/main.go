package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var sudokuBoard = [][]int{
	{5, 3, 0, 0, 7, 0, 0, 0, 0},
	{6, 0, 0, 1, 9, 5, 0, 0, 0},
	{0, 9, 8, 0, 0, 0, 0, 6, 0},
	{8, 0, 0, 0, 6, 0, 0, 0, 3},
	{4, 0, 0, 8, 0, 3, 0, 0, 1},
	{7, 0, 0, 0, 2, 0, 0, 0, 6},
	{0, 6, 0, 0, 0, 0, 2, 8, 0},
	{0, 0, 0, 4, 1, 9, 0, 0, 5},
	{0, 0, 0, 0, 8, 0, 0, 7, 9},
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", sudokuBoard)
	})

	e.POST("/update", func(c echo.Context) error {
		row, _ := strconv.Atoi(c.FormValue("row"))
		col, _ := strconv.Atoi(c.FormValue("col"))
		value, _ := strconv.Atoi(c.FormValue("value"))
		if row >= 0 && row < 9 && col >= 0 && col < 9 {
			if value >= 0 && value <= 9 {
				sudokuBoard[row][col] = value
			}
		}
		return c.Render(http.StatusOK, "board.html", sudokuBoard)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
