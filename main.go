package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Book struct {
	Id       int    `json:"Id"`
	Name     string `json:"Name"`
	Category string `json:"Category"`
}

var db *sql.DB

func main() {
	// Koneksi ke database
	var err error
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/library")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Membuat router
	r := mux.NewRouter()

	// Mengatur endpoint untuk method GET, POST, PUT, DELETE
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/add-books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Menjalankan server
	fmt.Println("Starting:")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Handler untuk mendapatkan semua buku
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Name, &book.Category); err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(books)
}

// Handler untuk membuat buku baru
func createBook(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    name := query.Get("name_books")
    category := query.Get("kategori_books")


    // Membuat struct Book dari data yang diterima dari query string
    book := Book{Name: name, Category: category}

    // Menyimpan buku ke database
    _, err := db.Exec("INSERT INTO books(name_books, kategori_books) VALUES(?, ?)", book.Name, book.Category)
    if err != nil {
        log.Fatal(err)
    }
    
    json.NewEncoder(w).Encode(book)
}


// Handler untuk memperbarui buku
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	_, err := db.Exec("UPDATE books SET name_books = ?, kategori_books = ? WHERE id_books = ?", book.Name, book.Category, id)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(book)
}

// Handler untuk menghapus buku
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM books WHERE id_books = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(`{"message": "Buku berhasil dihapus"}`))
}