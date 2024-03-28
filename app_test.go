package main

import (
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
