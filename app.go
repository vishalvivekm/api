package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}

func (app *App) getProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := getProductsFromDB(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(writer, http.StatusOK, products)
}

func (app *App) getProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, fmt.Sprintf("invalid product id: %v", vars["id"]))
		return
	}
	p := product{ID: key}
	err = p.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(writer, http.StatusNotFound, fmt.Sprintf("product with id=%v not found", p.ID))
		default:
			sendError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(writer, http.StatusOK, p)
}
func (app *App) createProduct(writer http.ResponseWriter, request *http.Request) {
	var p product

	err := json.NewDecoder(request.Body).Decode(&p)
	//log.Println(p)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid request payload")
		return
	}
	err = p.createProductInDB(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(writer, http.StatusOK, p)

}
func (app *App) updateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	var p product
	err = json.NewDecoder(request.Body).Decode(&p)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	p.ID = key
	err = p.updateProductInDB(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, fmt.Sprintf("error while updating the product in db, err: %v", err.Error()))
		return
	}
	sendResponse(writer, http.StatusOK, p)
}

func (app *App) deleteProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	p := product{ID: key}
	err = p.deleteProductFromDB(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, fmt.Sprintf("could not delete product from db, err: %v", err.Error()))
		return
	}
	sendResponse(writer, http.StatusOK, map[string]string{"result": "deleting successful"})
}
func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(response)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("could not write response, err: %v", err.Error()))
		return
	}
}
func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMsg := map[string]string{"error: ": err}
	sendResponse(w, statusCode, errorMsg)
}

/*//type jsonResponse struct {
//	Msg           string `json:"msg"`
//	UpdatedProduct string `json:"update product"`
//}

func (app *App) updateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}

	var p product
	err = json.NewDecoder(request.Body).Decode(&p)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	p.ID = key
	err = p.updateProductInDB(app.DB)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	//msg := jsonResponse{
	//	Msg:            "success",
	//	UpdatedProduct: fmt.Sprintf("updatedProduct: %v", p),
	//}

	//sendResponse(writer, http.StatusOK, msg)
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("could not marshal response, err: %v", err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("could not write response, err: %v", err.Error()))
	}
}*/
