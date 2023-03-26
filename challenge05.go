package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const PORT = ":8080" //declare PORT

/*
json:"" tags to explicitly specify the JSON property names for each field, especially when the JSON property names need to differ from the Go struct field names.
*/

// Book Struct with spesify name
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

// Insert Initial Data Books
var books []Book

// ? Main
func main() {

	books = append(books, Book{ID: 1, Title: "Book 1", Author: "Author 1", Desc: "Desc 1"})
	books = append(books, Book{ID: 2, Title: "Book 2", Author: "Author 2", Desc: "Desc 2"})
	books = append(books, Book{ID: 3, Title: "Book 3", Author: "Author 3", Desc: "Desc 3"})
	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/books/", getBookByID)
	http.HandleFunc("/books/create", createBook)
	http.HandleFunc("/books/update/", updateBook)
	http.HandleFunc("/books/delete/", deleteBook)
	fmt.Println("Application is listening on port", PORT)
	//log.Fatal() used to log a fatal error and terminate the program.
	//when Port in use/cant connect -> terminate the program
	log.Fatal(http.ListenAndServe(PORT, nil))
}

// ? All Book
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	//If Method is GET DO
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(books)
		return
	}
	//Else
	http.Error(w, "Method is not allowed", http.StatusBadRequest)
}

// ? Show Book by ID
func getBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	//id = ..., err = bool (if id false)
	//strconv.Atoi is used to convert the string "id" to an integer. If the string is not a valid integer, strconv.Atoi() will return an error.
	//
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/books/"):])

		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		for _, book := range books {
			//when book not found
			if book.ID == id {
				json.NewEncoder(w).Encode(book)
				return
			}
		}
		http.Error(w, "Book not found", http.StatusNotFound)
	} else {
		http.Error(w, "Method is not allowed", http.StatusBadRequest)
	}
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		author := r.FormValue("author")
		desc := r.FormValue("desc")

		if len(title) == 0 || len(author) == 0 {
			http.Error(w, "Title and author are required fields", http.StatusBadRequest)
			return
		}

		if desc == "" {
			desc = "Book doesn't have description"
		}

		book := Book{
			ID:     len(books) + 1,
			Title:  title,
			Author: author,
			Desc:   desc,
		}

		books = append(books, book)
		// Respond with the new book in JSON format
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(book)
		return
	} else {
		http.Error(w, "Method is not allowed", http.StatusBadRequest)
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		w.Header().Add("Content-type", "application/json")

		id, err := strconv.Atoi(r.URL.Path[len("/books/update/"):])
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		var updatedBook Book

		updatedBook = books[id-1]

		// Update the book details, if provided
		if r.FormValue("title") != "" {
			updatedBook.Title = r.FormValue("title")
		}
		if r.FormValue("author") != "" {
			updatedBook.Author = r.FormValue("author")
		}
		if r.FormValue("desc") != "" {
			updatedBook.Desc = r.FormValue("desc")
		}

		updatedBook.ID = id
		for index, book := range books {
			if book.ID == id {
				books[index] = updatedBook
				json.NewEncoder(w).Encode(updatedBook)
				return
			}
		}
		http.Error(w, "Book not found", http.StatusNotFound)
	} else {
		http.Error(w, "Method is not allowed", http.StatusBadRequest)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	id, err := strconv.Atoi(r.URL.Path[len("/books/delete/"):])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	for index, book := range books {
		if book.ID == id {
			books = append(books[:index], books[index+1:]...)
			json.NewEncoder(w).Encode(books)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}
