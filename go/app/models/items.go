package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
const dbSchema = "../db/items.db"
const dbSource = "../db/mercari.sqlite3"

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func GetItems() ([]Item, error) {
	rows, err := DB.Query("SELECT name, category FROM items")
	if (err != nil) {
		return nil, err
	}

	defer rows.Close()

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

func AddItem(newItem Item) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO items(name, category) VALUES(?, ?)")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newItem.Name, newItem.Category)
	if err != nil {
		return false, err
	}

	tx.Commit() // TODO: Should send err and check ?

	return true, nil
}

func SearchItem(key string) ([]Item, error) {
	q := fmt.Sprintf("SELECT name, category FROM items WHERE name='%s' or category='%s'", key, key)
	rows, err := DB.Query(q)
	if (err != nil) {
		return nil, err
	}

	defer rows.Close()

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
		return err;
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return err
	}
	return nil
}