package router

import (
	"github.com/flyervivek/golangpostgree/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/createstock", middleware.Createstock).Methods("POST")
	r.HandleFunc("/getallstocks", middleware.Getallstocks).Methods("GET")
	r.HandleFunc("/getstock/{id}", middleware.Getstock).Methods("GET")
	r.HandleFunc("/updatestock/{id}", middleware.Updatestock).Methods("PUT")
	r.HandleFunc("/deletestock/{id}", middleware.Deletestock).Methods("DELETE")

	return r

}
