package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

type Book struct {
	ID			string	`json:"id"`
	Title		string	`json:"title"`
	Author		string	`json:"author"`
	Quantity	int		`json:"quantity"`
}

var books = []Book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}


func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}


func createBook(c *gin.Context) {
	var newBook Book;
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookByID(id string) (*Book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("Book not Found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"message":"Book not found."},
		)
		return
		
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message":"Message id query parameter."},
		)
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"message":"Book not found."},
		)
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message":"Book not available."},
		)
		return
	}

	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message":"Message id query parameter."},
		)
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"message":"Book not found."},
		)
		return
	}

	book.Quantity++
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}