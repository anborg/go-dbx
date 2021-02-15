package main

import (
	"fmt"
	"log"
	"net/http"

	"./drivermgr"
	"./models"
)

// Create a custom Env struct which holds a connection pool.
type WithEnv struct {
	//books models.BookModel
	// Replace the reference to models.BookModel with an interface
	// describing its methods instead. All the other code remains exactly
	// the same.
	books interface {
		All() ([]models.Book, error)
	}
}

func main() {
	var err error

	// ds := drivermgr.DataSource{
	// 	DbType: "oracle",
	// 	URL: "hr/hr@//localhost:1521/xepdb1",
	// }
	ds := drivermgr.DataSource{
		DbType: "postgres",
		URL:    "postgres://postgres:postgres@localhost:5432/demodb?sslmode=disable",
	}

	db, err := drivermgr.GetConnection(ds)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Create an instance of Env containing the connection pool.
	withEnv := &WithEnv{books: models.BookModel{DB: db}}

	http.HandleFunc("/books", withEnv.booksIndex)
	http.ListenAndServe(":3000", nil)
}

// type BookRepo struct{
//     db *sql.DB
// }

// booksIndex sends a HTTP response listing all books.
func (env *WithEnv) booksIndex(w http.ResponseWriter, r *http.Request) {
	bks, err := env.books.All()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.Isbn, bk.Title, bk.Author, bk.Price)
	}
}
