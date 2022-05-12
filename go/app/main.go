package main

import (
	"fmt"
	"encoding/json"
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
	JsonFile = "items.json"
)

type Response struct {
	Message string `json:"message"`
}

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
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
	encoded_json, err := readJsonFile()
	if err != nil {
		return sendError(c, err.Error())
	}
	return c.JSONBlob(http.StatusOK, encoded_json)
}

func addItem(c echo.Context) error {
	// Get form data
	name := c.FormValue("name")
	category := c.FormValue("category")
	c.Logger().Infof("Receive item: %s %s", name, category)

	message := fmt.Sprintf("item received: %s", name)
	res := Response{Message: message}

	// Create json file
	f, err := os.OpenFile("items.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return sendError(c, err.Error())
	}

	// Setting close json file
	defer f.Close()

	// Read data to json file
	items := []Item{}
	save_items := Items{items}

	encoded_json, err := readJsonFile()
	if err != nil {
		return sendError(c, err.Error())
	}

	// Parse encoded_json
	if len(encoded_json) != 0 {
		err = json.Unmarshal(encoded_json, &save_items)
		if err != nil {
			return sendError(c, err.Error())
		}
	}

	// Add data to decode_data
	append_item := Item{Name: name, Category: category}
	save_items.Items = append(save_items.Items, append_item)

	// Set indent and encoding as JSON
	encode_items, err := json.MarshalIndent(save_items, "", " ")
	if err != nil {
		return sendError(c, err.Error())
	}

	// Write decode_data to json file
	err = os.WriteFile("items.json", encode_items, 0644)
	if err != nil {
		return sendError(c, err.Error())
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
	e.GET("/items", getItem)
	e.POST("/items", addItem)
	e.GET("/image/:itemImg", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
