package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func TestGetProducts(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Replace the global db variable with our mock database
	db = mockDB

	// Set expectations
	rows := sqlmock.NewRows([]string{"id", "name", "price", "description"}).
		AddRow(1, "Test Product", 9.99, "This is a test product")
	mock.ExpectQuery("SELECT id, name, price, description FROM products").WillReturnRows(rows)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProducts)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := []Product{{ID: 1, Name: "Test Product", Price: 9.99, Description: "This is a test product"}}
	var got []Product
	err = json.Unmarshal(rr.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProduct(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Replace the global db variable with our mock database
	db = mockDB

	// Set expectations
	rows := sqlmock.NewRows([]string{"id", "name", "price", "description"}).
		AddRow(1, "Test Product", 9.99, "This is a test product")
	mock.ExpectQuery("SELECT id, name, price, description FROM products WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", getProduct)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := Product{ID: 1, Name: "Test Product", Price: 9.99, Description: "This is a test product"}
	var got Product
	err = json.Unmarshal(rr.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateProduct(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Replace the global db variable with our mock database
	db = mockDB

	// Set expectations
	mock.ExpectExec("INSERT INTO products").
		WithArgs("New Product", 19.99, "This is a new product").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new product
	newProduct := Product{
		Name:        "New Product",
		Price:       19.99,
		Description: "This is a new product",
	}

	payload, _ := json.Marshal(newProduct)

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createProduct)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdProduct Product
	err = json.Unmarshal(rr.Body.Bytes(), &createdProduct)
	if err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}

	if createdProduct.Name != newProduct.Name {
		t.Errorf("Expected product name %v, got %v", newProduct.Name, createdProduct.Name)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateProduct(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Replace the global db variable with our mock database
	db = mockDB

	// Set expectations
	mock.ExpectExec("UPDATE products").
		WithArgs("Updated Product", 29.99, "This is an updated product", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create an updated product
	updatedProduct := Product{
		Name:        "Updated Product",
		Price:       29.99,
		Description: "This is an updated product",
	}

	payload, _ := json.Marshal(updatedProduct)

	// Create a request to pass to our handler
	req, err := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", updateProduct)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedProduct Product
	err = json.Unmarshal(rr.Body.Bytes(), &returnedProduct)
	if err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}

	if returnedProduct.Name != updatedProduct.Name {
		t.Errorf("Expected product name %v, got %v", updatedProduct.Name, returnedProduct.Name)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteProduct(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Replace the global db variable with our mock database
	db = mockDB

	// Set expectations
	mock.ExpectExec("DELETE FROM products WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a request to pass to our handler
	req, err := http.NewRequest("DELETE", "/products/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", deleteProduct)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// We make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
