package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description,omitempty`
}

func (b Book) ToJSON() []byte {
	ToJSON, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return ToJSON
}

func FromJSON(data []byte) Book {
	book := Book{}
	err := json.Unmarshal(data, &book)
	if err != nil {
		panic(err)
	}
	return book
}

var Books = map[string]Book{
	"078945613": Book{Title: "Book 1", Author: "Author 1", ISBN: "078945613"},
	"054623897": Book{Title: "Book 2", Author: "Author 2", ISBN: "054623897"},
}

func BookHandeFunc(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Path[len("/api/books/"):]

	if len(isbn) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch method := r.Method; method {
	case http.MethodGet:
		book, founded := GetBook(isbn)
		if founded {
			writeJSON(w, book)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		exists := UpdateBook(isbn, book)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		DeleteBook(isbn)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Unsupperted request method.`))
	}
}

func BooksHandeFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		books := AllBooks()
		writeJSON(w, books)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		isbn, created := CreateBook(book)
		if created {
			w.Header().Add("Location", "/api/books/"+isbn)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Unsupperted request method.`))
	}
}

func AllBooks() []Book {
	values := make([]Book, len(Books))
	idx := 0
	for _, book := range Books {
		values[idx] = book
		idx++
	}
	return values
}

func UpdateBook(isbn string, book Book) bool {
	_, founded := Books[isbn]
	if founded {
		Books[isbn] = book
	}
	return founded
}

func DeleteBook(isbn string) {
	delete(Books, isbn)
}

func GetBook(isbn string) (Book, bool) {
	book, founded := Books[isbn]
	return book, founded
}

func CreateBook(book Book) (string, bool) {
	_, exist := Books[book.ISBN]
	if exist {
		return "", false
	}
	Books[book.ISBN] = book
	return book.ISBN, true
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
