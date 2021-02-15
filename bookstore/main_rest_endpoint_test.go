package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"./models"
)

type mockBookRepo struct{}

func (m *mockBookRepo) All() ([]models.Book, error) {
	var bks []models.Book

	bks = append(bks, models.Book{"978-1503261969", "Emma", "Jayne Austen", 9.44})
	bks = append(bks, models.Book{"978-1505255607", "The Time Machine", "H. G. Wells", 5.99})

	return bks, nil
}

func (m *mockBookRepo) CountBooks() (int, error) {
	return 10, nil
}
func (m *mockBookRepo) CountAuthors() (int, error) {
	return 5, nil
}
func (m *mockBookRepo) CalcBooksPerAuthor() (float64, error) {
	return float64(5), nil
}
func TestBooksIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)

	env := WithEnv{repo: &mockBookRepo{}}

	http.HandlerFunc(env.booksIndex).ServeHTTP(rec, req)

	expected := "978-1503261969, Emma, Jayne Austen, £9.44\n978-1505255607, The Time Machine, H. G. Wells, £5.99\n"
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}

func TestStats(t *testing.T) { //TODO make a testmethod for model/books_test.go to test the CalcBookPerAuthor()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/stats", nil)
	env := WithEnv{repo: &mockBookRepo{}}
	http.HandlerFunc(env.bookStats).ServeHTTP(rec, req)
	expected := "BookCount: 10  AuthorCount: 5   BookPerAuthor: 5.00"
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}
