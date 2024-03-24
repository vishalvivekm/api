package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialise() error {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DBName)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}
func (app *App) Run(address string) {
	log.Fatalln(http.ListenAndServe(address, app.Router))
}
func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
}

func (app *App) getProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := getProductsFromDB(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
	}
	sendResponse(writer, http.StatusOK, products)
}
func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMsg := map[string]string{"error: ": err}
	sendResponse(w, statusCode, errorMsg)
}
