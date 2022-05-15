package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbSchema = "../db/items.db"
	dbSource = "../db/mercari.sqlite3"
)

var DB *sql.DB

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Image string `json:"image"`
}

func SetupDatabase() error {
	// Connect database
	db, err := sql.Open("sqlite3", dbSource)
	if err != nil {
		return err
	}
	DB = db

	// Create items table
	f, err := os.Open(dbSchema)
	if err != nil {
		return err
	}

	defer f.Close()

	schema, err := os.ReadFile(dbSchema)
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return err
	}
	return nil
}

func GetItem(query string) ([]Item, error) {
	if query == "" {
		query = "SELECT name, category FROM items"
	}

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Name, &item.Category); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

func GetItemById(id string) (Item, error) {
	item := Item{}

	stmt, err := DB.Prepare("SELECT name, category, image FROM items WHERE id = ?")
	if err != nil {
		return item, err
	}

	sqlErr := stmt.QueryRow(id).Scan(&item.Name, &item.Category, &item.Image)
	switch {
	case sqlErr == sql.ErrNoRows:
		return item, fmt.Errorf("No item with id %s", id)
	case sqlErr != nil:
		return item, sqlErr
	default:
		return item, nil
	}
}

func AddItem(newItem Item) (bool, error) {
	// Check image extension
	if !strings.HasSuffix(newItem.Image, ".jpg") {
		return false, fmt.Errorf("Image path does not end with .jpg")
	}

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO items(name, category, image) VALUES(?, ?, ?)")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	// Hashing
	begin := strings.LastIndex(newItem.Image, "/") + 1
	end := strings.Index(newItem.Image, ".")
	imgName := newItem.Image[begin:end]

	hashed := sha256.Sum256([]byte(imgName))

	newItem.Image = fmt.Sprintf("%x", hashed) + ".jpg"

	_, err = stmt.Exec(newItem.Name, newItem.Category, newItem.Image)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func SearchItem(key string) ([]Item, error) {
	q := fmt.Sprintf("SELECT name, category FROM items WHERE name='%s' or category='%s'", key, key)

	items, err := GetItem(q)

	return items, err
}
