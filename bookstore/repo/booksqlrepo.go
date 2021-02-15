package repo

import (
	"database/sql"

	"../models"
	"github.com/jmoiron/sqlx"
)

//BookSqlRepo Create a custom BookModel type which wraps the sql.DB connection pool.
type BookSqlRepo struct {
	DB *sqlx.DB
}

//CountBooks num of books
func (repo BookSqlRepo) CountBooks() (int, error) {
	var count int
	err := repo.DB.QueryRow("SELECT count(*) FROM BOOKS").Scan(&count)
	return count, err
}

//CountAuthors num of distinct authors
func (repo BookSqlRepo) CountAuthors() (int, error) {
	var count int
	err := repo.DB.QueryRow("SELECT count(distinct author) FROM BOOKS").Scan(&count)
	return count, err
}

//ById find Book by ISBN
func (repo BookSqlRepo) ById(isbn string) (models.Book, error) {
	row := repo.DB.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn)
	var book models.Book
	err := row.Scan(&book.Isbn, &book.Title, &book.Author, &book.Price)
	if err == sql.ErrNoRows {
		return book, nil //No rows is not an error
	} else if err != nil {
		return book, err //Other errors valid
	}
	return book, err
}

//CalcBooksPerAuthor - num of books written by auther
func (repo BookSqlRepo) CalcBooksPerAuthor() (float64, error) {
	bookCount, err1 := repo.CountBooks()
	if err1 != nil {
		return 0, err1
	}
	if bookCount == 0 { // avoid divide by zero
		return 0, nil
	}

	authorCount, err2 := repo.CountAuthors()
	if err2 != nil {
		return 0, err2
	}
	return float64(bookCount) / float64(authorCount), nil
}

// All returns a slice of all books in the books table.
func (repo BookSqlRepo) All() ([]models.Book, error) {
	rows, err := repo.DB.Query("SELECT * FROM BOOKS")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // important

	var books []models.Book

	for rows.Next() {
		var book models.Book

		err := rows.Scan(&book.Isbn, &book.Title, &book.Author, &book.Price)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	rows.Close() //explicitly calling at end of resultset use is good idea.
	//You should always check for an error at the end of the for rows.Next() loop
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (repo BookSqlRepo) InitWithTestData() error {
	sql_create_table_books := `CREATE TABLE IF NOT EXISTS books (
		isbn char(14) PRIMARY KEY NOT NULL,
		title varchar(255) NOT NULL,
		author varchar(255) NOT NULL,
		price decimal(5,2) NOT NULL
	);`

	sql_delete_books := `DELETE FROM books;`

	sql_insert_books := `INSERT INTO books (isbn, title, author, price) VALUES
	('978-1503261969', 'Emma', 'Jayne Austen', 9.44),
	('978-1505255607', 'The Time Machine', 'H. G. Wells', 5.99),
	('978-1503379640', 'The Prince', 'Niccolò Machiavelli', 6.99);
	-- commit;
	`

	_, err1 := repo.DB.Exec(sql_create_table_books)
	if err1 != nil {
		return err1
	}
	_, err1 = repo.DB.Exec(sql_delete_books)
	if err1 != nil {
		return err1
	}
	_, err1 = repo.DB.Exec(sql_insert_books)
	return err1
}
