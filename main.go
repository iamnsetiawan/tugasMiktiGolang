package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

// Struktur Book
type Book struct {
	ID     int json:"id"
	Title  string json:"title"
	Author string json:"author"
	Year   int json:"year"
}

// Fungsi untuk menginisialisasi database dan membuat tabel jika belum ada
func initDB() *sql.DB {
	database, err := sql.Open("sqlite3", "C:/sqlite3/test.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		year INTEGER
	);
	`
	_, err = database.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	return database
}

// Handler untuk mendapatkan semua buku
func getBooks(c echo.Context) error {
	db := initDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, title, author, year FROM books")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error querying books"})
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error scanning book data"})
		}
		books = append(books, book)
	}

	return c.JSON(http.StatusOK, books)
}

// Handler untuk menambahkan buku baru
func createBook(c echo.Context) error {
	db := initDB()
	defer db.Close()

	book := new(Book)
	if err := c.Bind(book); err != nil {
		return err
	}

	query := "INSERT INTO books (title, author, year) VALUES (?, ?, ?)"
	_, err := db.Exec(query, book.Title, book.Author, book.Year)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error inserting new book"})
	}

	return c.JSON(http.StatusCreated, book)
}

// Handler untuk memperbarui buku berdasarkan ID
func updateBook(c echo.Context) error {
	db := initDB()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	var book Book
	if err := c.Bind(&book); err != nil {
		return err
	}

	query := "UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?"
	_, err := db.Exec(query, book.Title, book.Author, book.Year, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating book"})
	}

	return c.JSON(http.StatusOK, book)
}

// Handler untuk menghapus buku berdasarkan ID
func deleteBook(c echo.Context) error {
	db := initDB()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	query := "DELETE FROM books WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error deleting book"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Main function
func main() {
	e := echo.New()

	// Routes CRUD API
	e.GET("/books", getBooks)          // Mendapatkan semua buku
	e.POST("/books", createBook)       // Menambahkan buku baru
	e.PUT("/books/:id", updateBook)    // Memperbarui buku
	e.DELETE("/books/:id", deleteBook) // Menghapus buku

	// Jalankan server di localhost:8080
	fmt.Println("Server running at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}