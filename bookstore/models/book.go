package models

import (
	"encoding/json"
	"fmt"
	"io"
)

//Book struct
type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32
}

func (b Book) String() string {
	return fmt.Sprintf("%s, %s, %s, £%.2f\n", b.Isbn, b.Title, b.Author, b.Price)
	//return fmt.Sprintf("Book: %s | %s | %s", b.Title, b.Author, b.Isbn)
}

//WriteJSON to demo interface usage,
// e.g io.Writer ifc use in  WriteJSON
// To write to bufferWriter and fileWriter
func (b *Book) WriteJSON(w io.Writer) error {
	js, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, err = w.Write(js)
	return err
}
