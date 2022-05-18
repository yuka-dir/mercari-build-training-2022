package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"mercari-build-training-2022/app/models"
)

const (
	ImgDir = "images"
)

type Response struct {
	Message string `json:"message"`
}

func sendError(c echo.Context, errMessage string) error {
	c.Logger().Errorf(errMessage)
	message := fmt.Sprintf("error: %s", errMessage)
	res := Response{Message: message}
	return c.JSON(http.StatusInternalServerError, res)
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func getItem(c echo.Context) error {
	items, err := models.GetItem("")
	if err != nil {
		return sendError(c, err.Error())
	}
	dbItems := models.Items{Items: items}
	if len(dbItems.Items) == 0 {
		res := Response{Message: "No Records Found"}
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, dbItems)
}

func getItemById(c echo.Context) error {
	id := c.Param("id")

	item, err := models.GetItemById(id)
	if err != nil {
		return sendError(c, err.Error())
	}
	return c.JSON(http.StatusOK, item)
}

func addItem(c echo.Context) error {
	name := c.FormValue("name")
	category := c.FormValue("category")

	// Read file
	image, err := c.FormFile("image")
	if err != nil {
		return sendError(c, err.Error())
	}
	if !strings.HasSuffix(image.Filename, ".jpg") {
		return sendError(c, "Image path does not end with .jpg")
	}
	src, err := image.Open()
	if err != nil {
		return sendError(c, err.Error())
	}
	defer src.Close()

	// Hashing
	extension := strings.Index(image.Filename, ".")
	hashed := sha256.Sum256([]byte(image.Filename[:extension]))
	hashedImgName := fmt.Sprintf("%x", hashed) + ".jpg"

	// Destination
	dst, err := os.Create(path.Join(ImgDir, hashedImgName))
	if err != nil {
		return sendError(c, err.Error())
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return sendError(c, err.Error())
	}

	item := models.Item{Name: name, Category: category, Image: hashedImgName}

	if err = models.AddItem(item); err != nil {
		return sendError(c, err.Error())
	}
	message := fmt.Sprintf("item received: %s", item.Name)
	res := Response{Message: message}
	return c.JSON(http.StatusOK, res)
}

func searchItem(c echo.Context) error {
	key := c.FormValue("keyword")
	items, err := models.SearchItem(key)
	if err != nil {
		return sendError(c, err.Error())
	}
	dbItems := models.Items{Items: items}
	if len(dbItems.Items) == 0 {
		res := Response{Message: "No Records Found"}
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, dbItems)
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

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
	// Database
	err := models.SetupDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		return
	}

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
	e.GET("/items/:id", getItemById)
	e.POST("/items", addItem)
	e.GET("/search", searchItem)
	e.GET("/image/:imageFilename", getImg)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
