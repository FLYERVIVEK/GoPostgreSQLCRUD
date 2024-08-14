package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/flyervivek/golangpostgree/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func Createconnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Connected to Database")
	return db

}

func Createstock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to Decode the Request Body %v", err)
	}

	insertedId := insertthestock(stock)
	res := response{
		ID:      insertedId,
		Message: "Stock Created Successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func insertthestock(stock models.Stock) int64 {
	db := Createconnection()
	defer db.Close()
	var id int64

	sqlstatement := `INSERT INTO stocks (name, price, company) VALUES ($1,$2,$3) RETURNING stockid`

	err := db.QueryRow(sqlstatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id

}

func Getallstocks(w http.ResponseWriter, r *http.Request) {

	stock, err := getallthestocks()
	if err != nil {
		log.Fatalf("Unable to get all stocks %v", err)
	}

	json.NewEncoder(w).Encode(stock)

}

func getallthestocks() ([]models.Stock, error) {
	db := Createconnection()
	defer db.Close()
	var stocks []models.Stock

	sqlstatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable to scan the rows %v", err)
		}

		stocks = append(stocks, stock)

	}

	return stocks, err

}

func Getstock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to get the id from request %v", err)
	}

	stock, err := getthestock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get the stock %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func getthestock(id int64) (models.Stock, error) {

	db := Createconnection()
	defer db.Close()

	var stock models.Stock

	sqlstatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlstatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return stock, err

}

func Updatestock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to get the id from request %v", err)
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body %v", err)
	}

	updatedrows := updatethestock(int64(id), stock)
	msg := fmt.Sprintf("Successfully Updated, Total rows affected are %v", updatedrows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func updatethestock(id int64, stock models.Stock) int64 {
	db := Createconnection()
	defer db.Close()
	sqlstatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	res, err := db.Exec(sqlstatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowaffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	return rowaffected

}

func Deletestock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to get the id %v", err)
	}

	rowsaffected := deletethestock(int64(id))

	msg := fmt.Sprintf("Successfully deleted the stock, rows affected %v", rowsaffected)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func deletethestock(id int64) int64 {
	db := Createconnection()
	defer db.Close()
	sqlstatement := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(sqlstatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowwsaffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowwsaffected

}
