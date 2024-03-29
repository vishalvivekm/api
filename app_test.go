package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatalln("error occurred while initialising the database")
	}
	createTable()
	m.Run()

}
func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products(
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    quantity int,
    price float(10,7),
    PRIMARY KEY (id)
);`
	if _, err := a.DB.Exec(createTableQuery); err != nil {
		log.Fatalln(err)
	}

}
func clearTable() {
	_, err := a.DB.Exec("DELETE  from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
	if err != nil {
		log.Fatalln(err)
	}
}
func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name, quantity, price)  VALUES ('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("monitor", 100, 400.00)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}
func checkStatusCode(t *testing.T, expectedStatusCode, actualStatuCode int) {
	if expectedStatusCode != actualStatuCode {
		t.Errorf("expected status: %v, received: %v", expectedStatusCode, actualStatuCode)
	}
}
func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}
func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name": "chair", "quantity": 1, "price": 100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	//log.Println(m)

	if m["name"] != "chair" {
		t.Errorf("expected name: %v, got: %v", "chair", m["name"])
	}
	if m["quantity"] != 1.00 {
		t.Errorf("expected quantity: %v, got: %v", 1.00, m["quantity"])
	}

}
func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	var initialGetResponseValue map[string]interface{}

	err := json.Unmarshal(response.Body.Bytes(), &initialGetResponseValue)
	if err != nil {
		t.Errorf(err.Error())
	}

	putreqBody := []byte(`{"name": "pens", "quantity": 5, "price": 20.00}`)
	//reqMap := map[string]interface{}{"name": "pens", "quantity": 5, "price": 20.00}
	//putreqBody, _ := json.Marshal(reqMap)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(putreqBody))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	var getResponseValueAfterPutReq map[string]interface{}

	err = json.Unmarshal(response.Body.Bytes(), &getResponseValueAfterPutReq)
	if err != nil {
		t.Errorf("can not unmarshal, %v", err.Error())
	}
	log.Println(initialGetResponseValue)
	log.Println(getResponseValueAfterPutReq)
	/*
	   2024/03/29 23:54:02 map[id:1 name:connector price:10 quantity:10]
	   2024/03/29 23:54:02 map[id:1 name:pens price:20 quantity:5]

	*/
	if initialGetResponseValue["id"] != getResponseValueAfterPutReq["id"] {
		t.Errorf("expected id: %v, got: %v", initialGetResponseValue["id"], getResponseValueAfterPutReq["id"])
	}
	if initialGetResponseValue["name"] == getResponseValueAfterPutReq["name"] {
		t.Errorf("expected name: %v, got: %v", getResponseValueAfterPutReq["name"], initialGetResponseValue["name"])
	}
	if initialGetResponseValue["quantity"] == getResponseValueAfterPutReq["quantity"] {
		t.Errorf("expected quantity: %v, got: %v", getResponseValueAfterPutReq["quantity"], initialGetResponseValue["quantity"])
	}
	if initialGetResponseValue["price"] == getResponseValueAfterPutReq["price"] {
		t.Errorf("expected price: %v, got: %v", getResponseValueAfterPutReq["quantity"], initialGetResponseValue["price"])
	}
}

//func TestUpdateProduct(t *testing.T) {
//	clearTable()
//	addProduct("connector", 10, 10)
//
//	req, _ := http.NewRequest("GET", "/product/1", nil)
//	response := sendRequest(req)
//	checkStatusCode(t, http.StatusOK, response.Code)
//	//log.Printf("inside updatePRoductTEst: %v", string(response.Body.Bytes())) //  inside updatePRoductTEst: {"id":1,"name":"connector","quantity":10,"price":10}
//	q := map[string]interface{}{"name": "pens", "quantity": 5, "price": 20.00}
//	bf, err := json.Marshal(q)
//	if err != nil {
//		log.Fatalf("can't marshal q, %v", err)
//	}
//	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(bf))
//	response = sendRequest(req)
//	checkStatusCode(t, http.StatusOK, response.Code)
//	//log.Printf("inside updatePRoductTEst: %v", string(response.Body.Bytes())) // inside updatePRoductTEst: {"id":1,"name":"pens","quantity":5,"price":20}
//
//	req, _ = http.NewRequest("GET", "/product/1", nil)
//	response = sendRequest(req)
//	checkStatusCode(t, http.StatusOK, response.Code)
//
//}
