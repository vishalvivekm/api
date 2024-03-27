package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json: "name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProductsFromDB(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var products []product
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price FROM products where id=%v", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *product) createProductInDB(db *sql.DB) error {
	query := fmt.Sprintf("insert into products(name, quantity, price) values('%v', %v, %v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id) //reflect.TypeOf(id) // int64
	return nil
}
func (p *product) updateProductInDB(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no row with given id exists")
	}
	return err
}
