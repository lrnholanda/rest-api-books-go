package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// form: "title" to map the JSON field name to the struct
// binding: "required" to enforce the value is required

// Create the request and response structure for the books

type BookCreateUpdateRequest struct {
	Title  string `form:"title" binding:"required"`
	Author string `form:"author"`
}

// json: "id" to map the struct Name to its Json field name

type BookResponse struct {
	Id        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	Author    string             `json:"author"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type BookListResponse struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"` // bson to map mongo _id to id
	Title string             `json:"title"`
}

// Create the Database model structure:

type Book struct {
	field.DefaultField `bson:"inline"`
	Title              string `bson:"title" validate:"required"`
	Author             string `bson:"author"`
}

func CreateBook(ctx *gin.Context) {
	var newBook BookCreateUpdateRequest

	// to bind the received JSON to BookRequest to strip the unnecessary fields.
	if err := ctx.ShouldBind(&newBook); err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	// setting data to book model struct
	book := Book{
		Title:  newBook.Title,
		Author: newBook.Author,
	}
	_, err := collection.InsertOne(ctx, &book) //Inserting the Book Data to database

	// to send error response if any error occurs
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Something went wrong, Try again after sometime")
		return
	}

	// to send sucess response on completion
	ctx.JSON(http.StatusCreated, GetBooksResponse(book))
}

func GetBooksResponse(book Book) (bookResponse BookResponse) {
	// setting response for book
	bookResponse = BookResponse{
		Id:        book.DefaultField.Id,
		Title:     book.Title,
		Author:    book.Author,
		CreatedAt: book.CreateAt,
		UpdatedAt: book.UpdateAt,
	}
	return
}

func GetBook(ctx *gin.Context) {

	// to get and convert the received path variable to desired type
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	//Getting the Book Data from database
	var book Book
	err = collection.Find(ctx, bson.M{"_id": bookId}).One(&book)

	// to send error response if any error occurs
	if err != nil {
		ctx.JSON(http.StatusNotFound, "Book Not Found")
		return
	}

	// to send succes response on completion
	ctx.JSON(http.StatusOK, GetBooksResponse(book))
}

// Update a book API

func UpdateBook(ctx *gin.Context) {
	// to get and covert the received path variable to desired type
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Book ID")
		return
	}

	var newBook BookCreateUpdateRequest

	// to bind the received JSON to BookRequest to strip the unnecessary fields.
	if err := ctx.ShouldBind(&newBook); err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	// Getting the Book Data from database
	var book Book
	err = collection.Find(ctx, bson.M{"_id": bookId}).One(&book)

	// to send error response if any error occurs
	if err != nil {
		ctx.JSON(http.StatusNotFound, "Book Not Found")
		return
	}

	// set the updated value in the book
	book.Author = newBook.Author
	book.Title = newBook.Title

	// update in database
	err = collection.ReplaceOne(ctx, bson.M{"_id": bookId}, &book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Something went wrong, Try again after sometime")
		return
	}

	// to send success response on completion
	ctx.JSON(http.StatusOK, GetBooksResponse(book))
}

func DeleteBook(ctx *gin.Context) {

	// to get and convert the received path variable to desired type
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	// Getting the Book Data from database
	var book Book
	err = collection.Find(ctx, bson.M{"_id": bookId}).One(&book)

	// to send error response if any error occurs
	if err != nil {
		ctx.JSON(http.StatusNotFound, "Book not found")
		return
	}

	// Deleting the book
	err = collection.RemoveId(ctx, bookId)

	// to send error response if any error occurs
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Something went wrong, Try again after sometime")
		return
	}
	// to send success response on completion
	ctx.JSON(http.StatusOK, true)
}

func GetBooks(ctx *gin.Context) {

	// Getting the Book Data to database
	var books []BookListResponse
	err := collection.Find(ctx, bson.M{}).All(&books)

	// to send error response if any error occurs

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, "Something went wrong, Try again after sometime")
		return
	}
	// to send success response on completion
	ctx.JSON(http.StatusOK, books)
}
