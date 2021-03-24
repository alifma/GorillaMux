package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16.2)"`
}

type Result struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Message string `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:password@/GorillaMux?charset=utf8&parseTime=true")
	if err != nil {
		log.Println("Koneksi Gagal", err)
	} else {
		log.Println("Koneksi Berhasil")
	}
	db.AutoMigrate(&Product{})
	handleRequest()
}

func handleRequest() {
	log.Println("GOAPI is Running on http://127.0.0.1:8000")
	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.HandleFunc("/", homePage).Methods("GET")
	apiRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	apiRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	apiRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	apiRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	apiRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", apiRouter))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ini Halaman Home")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads,_ := ioutil.ReadAll(r.Body)
	var product Product
	json.Unmarshal(payloads, &product)
	db.Create(&product)
	res := Result{Code: 200, Data: product, Message: "Add product Success"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Display product Success"}
	results, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	var product Product
	db.First(&product, productId)
	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: product, Message: "Display Detail Product Sukses"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	payloads,_ := ioutil.ReadAll(r.Body)

	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	var product Product
	db.First(&product, productId)
	db.Model(&product).Updates(productUpdate)

	res := Result{Code: 200, Data: product, Message: "Update product Success"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}


func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	var product Product
	db.First(&product, productId)
	db.Delete(&product)

	res := Result{Code: 200, Message: "Delete product Success"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}