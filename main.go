package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()
	r.Get("/", HandleRoot)
	r.Get("/products", ListProducts)
	r.Get("/buyers", ListBuyers)
	r.Get("/load_data", LoadData)
	r.Get("/buyer", BuyerDetail)
	r.Get("/product", GetProduct)
	r.Get("/other_buyer", GetOtherBuyer)

	http.ListenAndServe(":3000", r)
}
