package models

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
const dbSource = "../db/mercari.sqlite3"

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Id		 int	`json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

func GetItems() ([]Item, error) {
	rows, err := DB.Query("SELECT * FROM items")
	if (err != nil) {
		return nil, err
	}

	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Id, &item.Name, &item.Category); err != nil {
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

	stmt, err := tx.Prepare("INSERT INTO items(id, name, category) VALUES(?, ?, ?)")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newItem.Id, newItem.Name, newItem.Category)
	if err != nil {
		return false, err
	}

	tx.Commit() // TODO: Should send err and check ?

	return true, nil
}

func SearchItem(key string) ([]Item, error) {
	q := fmt.Sprintf("SELECT id, name, category FROM items WHERE name='%s' or category='%s'", key, key)
	rows, err := DB.Query(q)
	if (err != nil) {
		return nil, err
	}

	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Id, &item.Name, &item.Category); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", dbSource)
	if err != nil {
		return err
	}
	DB = db
	return nil
}