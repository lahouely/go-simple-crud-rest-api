package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "test_db")
	if err != nil {
		log.Fatal("Error occured while initialising the database.")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10,7),
		PRIMARY KEY(id)
		);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products;")
	a.DB.Exec("ALTER TABLE products AUTO_INCREMENT=1;")
}

func addProduct(name string, quantity int, price float64) (int64, error) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES ('%v', %v, %v)", name, quantity, price)
	result, err := a.DB.Exec(query)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func TestGetProduct(t *testing.T) {
	clearTable()
	id, err := addProduct("mouse", 9, 36.00)
	if err != nil {
		t.Errorf("Error while adding a Product: %v.", err)
	}
	request, _ := http.NewRequest("GET", "/product/"+strconv.Itoa(int(id)), nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{
			"name":"chair",
			"quantity":8,
			"price":35
		}`)
	request, _ := http.NewRequest("POST", "/product/", bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")
	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, got: %v, type:%T", "chair", m["name"], m["name"])
	}
	if m["quantity"] != float64(8) {
		t.Errorf("Expected name: %v, got: %v, type:%T", 8, m["quantity"], m["quantity"])
	}
	if m["price"] != float64(35) {
		t.Errorf("Expected name: %v, got: %v, type:%T", 35, m["price"], m["price"])
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	id, err := addProduct("screen", 9, 36.00)
	if err != nil {
		t.Errorf("Error while adding a Product: %v.", err)
	}

	request, _ := http.NewRequest("GET", "/product/"+strconv.Itoa(int(id)), nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/product/"+strconv.Itoa(int(id)), nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/product/"+strconv.Itoa(int(id)), nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	id, err := addProduct("harddrive", 7, 346.00)
	if err != nil {
		t.Errorf("Error while adding a Product: %v.", err)
	}

	request, _ := http.NewRequest("GET", "/product/"+strconv.Itoa(int(id)), nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
	fmt.Println(">>>", response.Code)

	var oldm map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldm)

	product := []byte(`{
		"name":"ssd",
		"quantity":4,
		"price":70
	}`)
	request, _ = http.NewRequest("PUT", "/product/"+strconv.Itoa(int(id)), bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")
	response = sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "ssd" {
		t.Errorf("Expected name: %v, got: %v, type:%T", "ssd", m["name"], m["name"])
	}
	if m["quantity"] != float64(4) {
		t.Errorf("Expected name: %v, got: %v, type:%T", 4, m["quantity"], m["quantity"])
	}
	if m["price"] != float64(70) {
		t.Errorf("Expected name: %v, got: %v, type:%T", 70, m["price"], m["price"])
	}

}

func checkStatusCode(t *testing.T, expectedStatusCode, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, recieved: %v.", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}
