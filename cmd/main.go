package main

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initializeBoard(sudokuBoard *[9][9]int, K int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			sudokuBoard[i][j] = 0
		}
	}

	fillDiagonals(sudokuBoard)

	fillRemaining(sudokuBoard, 0, 3)

	removeKDigits(sudokuBoard, K)
}

func fillDiagonals(sudokuBoard *[9][9]int) {
	for i := 0; i < 3; i += 1 {
		fillBox(sudokuBoard, i*3, i*3)
	}
}

func fillBox(sudokuBoard *[9][9]int, row int, col int) {
	var num int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for {
				num = generateNumber(9)
				if unUsedInBox(sudokuBoard, row, col, num) {
					break
				}
			}
			sudokuBoard[i+row][j+col] = num
		}
	}
}

func generateNumber(n int) int {
	return 1 + rand.Intn(n)
}

func unUsedInBox(sudokuBoard *[9][9]int, row int, col int, num int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if sudokuBoard[row+i][col+j] == num {
				return false
			}
		}
	}
	return true
}

func unUsedInRow(sudokuBoard *[9][9]int, row int, num int) bool {
	for i := 0; i < 9; i++ {
		if sudokuBoard[row][i] == num {
			return false
		}
	}
	return true
}

func unUsedInCol(sudokuBoard *[9][9]int, col int, num int) bool {
	for i := 0; i < 9; i++ {
		if sudokuBoard[i][col] == num {
			return false
		}
	}
	return true
}

func checkIfSafe(sudokuBoard *[9][9]int, i, j, num int) bool {
	return unUsedInRow(sudokuBoard, i, num) && unUsedInCol(sudokuBoard, j, num) && unUsedInBox(sudokuBoard, i-i%3, j-j%3, num)
}

func fillRemaining(sudokuBoard *[9][9]int, i, j int) bool {
	if j >= 9 && i < 9-1 {
		i = i + 1
		j = 0
	}
	if i >= 9 && j >= 9 {
		return true
	}
	if i < 3 {
		if j < 3 {
			j = 3
		}
	} else if i < 6 {
		if j == int(i/3)*3 {
			j = j + 3
		}
	} else {
		if j == 6 {
			i = i + 1
			j = 0
			if i >= 9 {
				return true
			}
		}
	}

	for num := 1; num <= 9; num++ {
		if checkIfSafe(sudokuBoard, i, j, num) {
			sudokuBoard[i][j] = num
			if fillRemaining(sudokuBoard, i, j+1) {
				return true
			}
			sudokuBoard[i][j] = 0
		}
	}
	return false
}

func removeKDigits(sudokuBoard *[9][9]int, K int) {
	for i := 0; i < K; i++ {
		for {
			row := rand.Intn(9)
			col := rand.Intn(9)
			if sudokuBoard[row][col] != 0 {
				sudokuBoard[row][col] = 0
				break
			}
		}
	}
}

func checkIfWon(sudokuBoard *[9][9]int) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudokuBoard[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func printBoard(sudokuBoard *[9][9]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%d ", sudokuBoard[i][j])
		}
		fmt.Printf("\n")
	}
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

type FormData struct {
	SudokuBoard [9][9]int
	Error       string
}

func newFormData(k int) FormData {
	data := FormData{
		SudokuBoard: [9][9]int{},
		Error:       "",
	}
	initializeBoard(&data.SudokuBoard, k)
	return data
}

func main() {
	formData := newFormData(30)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", formData)
	})

	e.POST("/update", func(c echo.Context) error {
		row, _ := strconv.Atoi(c.FormValue("row"))
		col, _ := strconv.Atoi(c.FormValue("col"))
		value, _ := strconv.Atoi(c.FormValue("value"))
		formData.Error = ""

		if row >= 0 && row < 9 && col >= 0 && col < 9 {
			if value >= 0 && value <= 9 {
				if value == 0 || checkIfSafe(&formData.SudokuBoard, row, col, value) {
					formData.SudokuBoard[row][col] = value
					if checkIfWon(&formData.SudokuBoard) {
						return c.Render(http.StatusOK, "won.html", formData)
					}
					return c.Render(http.StatusOK, "game.html", formData)
				} else {
					formData.Error = "Invalid input"
					return c.Render(http.StatusBadRequest, "game.html", formData)
				}
			}
		}
		formData.Error = "Invalid input"
		return c.Render(http.StatusBadRequest, "game.html", formData)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
