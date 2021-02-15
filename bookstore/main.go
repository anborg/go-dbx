package main

import (
	"fmt"
	"log"
	"net/http"

	"./drivermgr"
	"./models"
	"./repo"
)

func main() {
	console_msg := `Try urls :
	http://localhost:3000/books
	http://localhost:3000/books/stats
	http://localhost:3000/books/show?isbn=978-1503261969
	`

	var err error
	ds := getDSSqlte()                     // read ds from config
	db, err := drivermgr.GetConnection(ds) // get connection
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Defer db cloes at main() level ok?

	withEnv := &WithEnv{repo: repo.BookSqlRepo{DB: db}} // wire concrete SQL-repo for BookRepo interface
	// err1 := withEnv.repo.InitWithTestData()
	// if err1 != nil {
	// 	log.Fatal(err1)
	// }

	http.HandleFunc("/books", withEnv.booksIndex)      // wire Handler
	http.HandleFunc("/books/stats", withEnv.bookStats) // wire Handler
	http.HandleFunc("/books/show", withEnv.booksShow)
	log.Println(console_msg)

	http.ListenAndServe(":3000", nil)
}

//WithEnv Create a custom Env struct which holds a connection pool.
type WithEnv struct {
	repo BookRepo //Use interface in
}

//BookRepo generic interface implmented by SQL/other repo.
type BookRepo interface {
	InitWithTestData() error
	All() ([]models.Book, error)
	CountBooks() (int, error)
	CountAuthors() (int, error)
	CalcBooksPerAuthor() (float64, error)
	ById(string) (models.Book, error)
}

// booksIndex sends a HTTP response listing all books.
func (env *WithEnv) booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bks, err := env.repo.All()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.Isbn, bk.Title, bk.Author, bk.Price)
		//fmt.Fprintf(w, "%s \n", bk.String())
		//bk.WriteJSON(w)
	}
}

func (env *WithEnv) bookStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bookCount, _ := env.repo.CountBooks() //TODo convert to struct
	authorCount, _ := env.repo.CountAuthors()
	booksPerAuthor, _ := env.repo.CalcBooksPerAuthor()
	fmt.Fprintf(w, "BookCount: %d  AuthorCount: %d   BookPerAuthor: %.2f", bookCount, authorCount, booksPerAuthor)
}

func (env *WithEnv) booksShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var bk models.Book
	bk, err := env.repo.ById(isbn)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	//if bk != nil {
	fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.Isbn, bk.Title, bk.Author, bk.Price)
	//}
	return
}

// --- utility functions-----
func getDSPostgres() drivermgr.Datasource {
	ds := drivermgr.Datasource{
		DbType: "postgres",
		URL:    "postgres://postgres:postgres@localhost:5432/demodb?sslmode=disable&currentSchema=demo", //In PG Use currentSchema=demo
	}
	return ds
}
func getDSOracle() drivermgr.Datasource {
	ds := drivermgr.Datasource{
		DbType: "oracle",
		URL:    "hr/hr@//localhost:1521/xepdb1", //In oracle username maps to schema
	}
	return ds
}
func getDSSqlte() drivermgr.Datasource {
	ds := drivermgr.Datasource{
		DbType: "sqlite",
		// URL:    ":memory:",
		URL: "./db/gotest.sqlite",
	}
	return ds
}
func getDSMssqlserver() drivermgr.Datasource {
	ds := drivermgr.Datasource{
		//sqlserver://localhost:1433;databaseName=tempdb;user=sa;password=admin1234!
		URL:    "sqlserver://sa:admin1234!@localhost:1433/?databaseName=tempdb&param2=value", // MSSQL do not have a way mentioning schema prefix in DS level. So need to add prefix to query level?
		DbType: "sqlserver",
	}
	return ds
}

// Demo of using request context - dont use this for passing long lasting objects like Connectionpool
// Create some middleware which swaps out the existing request context
// with new context.Context value containing the connection pool.
// func injectDB(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), "db", db)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }
