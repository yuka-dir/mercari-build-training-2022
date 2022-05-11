package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	ImgDir = "image"
)

type Response struct {
	Message string `json:"message"`
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func addItem(c echo.Context) error {
	// Get form data
	name := c.FormValue("name")
	c.Logger().Infof("Receive item: %s", name)

	message := fmt.Sprintf("item received: %s", name)
	res := Response{Message: message}

	// Save data to items.json
	// // Create json file
	f, err := os.OpenFile("items.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		c.Logger().Infof(err.Error())
	}

	// // Read data to json file
	raw_data, err := os.ReadFile("items.json")
	if err != nil {
		c.Logger().Infof(err.Error())
	}

	var decode_data map[string]interface{}
	json.Unmarshal([]byte(raw_data), &decode_data)

	c.Logger().Infof("Mapping!!!: %s", fmt.Sprintf("%s", decode_data["test"])) // TODO: 後で消す
	items := decode_data["items"].([]interface{})
	c.Logger().Infof("\"name\": %s", fmt.Sprintf("%s", items[0].(map[string]interface{})["name"])) // TODO: 後で消す
	c.Logger().Infof("\"category\": %s", fmt.Sprintf("%s", items[0].(map[string]interface{})["category"])) // TODO: 後で消す

	// // Add data to json file

	// // Close json file
	if err := f.Close(); err != nil {
		c.Logger().Infof(err.Error())
	}

	return c.JSON(http.StatusOK, res)
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

	// Routes
	e.GET("/", root)
	e.POST("/items", addItem)
	e.GET("/image/:itemImg", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
