package models

import (
	"database/sql"

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

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", dbSource)
	if err != nil {
		return err
	}
	DB = db
	return nil
}