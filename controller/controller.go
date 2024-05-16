package controller

import (
	"context"
	"crud/helper"
	"crud/models"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = helper.ConnectDB()

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []models.Book
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var book models.Book
		err := cur.Decode(&book)
		if err != nil {
			helper.GetError(err, w)
			return
		}
		books = append(books, book)
	}

	if err := cur.Err(); err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	var existingBook models.Book
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingBook)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func UpdateBook(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	var existingBook models.Book
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingBook)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	update := bson.M{
		"$set": book,
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	book.ID = objID
	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	var existingBook models.Book
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingBook)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"result": "success", "message": "Book Deleted Successfully!"})
}
