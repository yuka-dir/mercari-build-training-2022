package main

import (
	"fmt"
	// "encoding/json"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"go/app/models"
)

const (
	ImgDir = "image"
	JsonFile = "items.json"
)

type Response struct {
	Message string `json:"message"`
}

func sendError(c echo.Context, err_message string) error {
	c.Logger().Errorf(err_message)
	message := fmt.Sprintf("error: %s", err_message)
	res := Response{Message: message}
	return c.JSON(http.StatusInternalServerError, res)
}

func readJsonFile() ([]byte, error) {
	encoded_json, err := os.ReadFile(JsonFile)
	if err != nil {
		return encoded_json, err
	}
	return encoded_json, nil
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func getItem(c echo.Context) error {
	items, err := models.GetItems()
	if err != nil {
		return sendError(c, err.Error())
	}
	db_items := models.Items{items}
	if db_items.Items == nil {
		res := Response{Message: "No Records Found"}
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, db_items)
}

func addItem(c echo.Context) error {
	id, _ := strconv.Atoi(c.FormValue("id"))
	name := c.FormValue("name")
	category := c.FormValue("category")

	item := models.Item{Id: id, Name: name, Category: category}

	success, err := models.AddItem(item)

	if success {
		message := fmt.Sprintf("item received: %s", item.Name)
		res := Response{Message: message}
		return c.JSON(http.StatusOK, res)
	} else {
		return sendError(c, err.Error())
	}
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("itemImg"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	front_url := os.Getenv("FRONT_URL")
	if front_url == "" {
		front_url = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{front_url},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Connect Database
	models.ConnectDatabase()

	// Routes
	e.GET("/", root)
	e.GET("/items", getItem)
	e.POST("/items", addItem)
	e.GET("/image/:itemImg", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
