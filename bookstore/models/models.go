package models

import (
	"database/sql"
)

//Book struct
type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32
}

//BookModel Create a custom BookModel type which wraps the sql.DB connection pool.
type BookModel struct {
	DB *sql.DB
}

// All returns a slice of all books in the books table.
func (m BookModel) All() ([]Book, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := m.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bks []Book

	for rows.Next() {
		var bk Book

		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
		if err != nil {
			return nil, err
		}

		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bks, nil
}
