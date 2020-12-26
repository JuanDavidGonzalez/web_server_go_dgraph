package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/dgo/protos/api"
)

const URL string = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/"

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func ListProducts(w http.ResponseWriter, r *http.Request) {

	dg, cancel := newClient()

	q := `query Data{
		data(func: type(Product)) {
			id
			name
			price
		}
	}`

	ctx := context.Background()
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	type Root struct {
		Data []Product `json:"data"`
	}
	var root Root
	err = json.Unmarshal(resp.Json, &root)
	if err != nil {
		log.Fatal(err)
	}
	cancel()
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func ListBuyers(w http.ResponseWriter, r *http.Request) {

	dg, cancel := newClient()

	q := `query Data{
		data(func: type(Buyer)) {
			id
			name
			age
		}
	}`

	ctx := context.Background()
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	type Root struct {
		Data []Buyer `json:"data"`
	}
	var root Root
	err = json.Unmarshal(resp.Json, &root)
	if err != nil {
		log.Fatal(err)
	}
	cancel()
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func LoadData(w http.ResponseWriter, r *http.Request) {

	dates, exist := r.URL.Query()["date"]
	if !exist || len(dates[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "date invalid")
		return
	}
	date := dates[0]

	types, exist := r.URL.Query()["type"]
	if !exist || len(types[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "type invalid")
		return
	}
	typeData := types[0]

	response, err := http.Get(URL + typeData + "?date=" + date)
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	var result = Response{}
	if typeData == "buyers" {
		result = loadBuyers(response)
	} else if typeData == "products" {
		result = loadProducts(response)
	} else {
		result = Response{Status: "400", Message: "invalid type"}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func loadBuyers(r *http.Response) Response {

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	dg, cancel := newClient()

	var resp = Response{}

	op := &api.Operation{}
	op.Schema = `
		name: string @index(exact) .
		age: int .
		type Buyer {
			name
			age
		}
	`
	ctx := context.Background()
	err1 := dg.Alter(ctx, op)
	if err1 != nil {
		log.Fatal(err1)
	}
	var responseObject BuyerReponse
	json.Unmarshal(responseData, &responseObject)

	for index, item := range responseObject {
		if index > 10 {
			break
		}

		p := Buyer{
			Id:    item.Id,
			Name:  item.Name,
			Age:   item.Age,
			DType: []string{"Buyer"},
		}

		mu := &api.Mutation{
			CommitNow: true,
		}
		pb, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}

		mu.SetJson = pb
		_, err2 := dg.NewTxn().Mutate(ctx, mu)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	cancel()

	resp = Response{Status: "200", Message: "buyers loaded"}
	return resp
}

func loadProducts(r *http.Response) Response {

	reader := csv.NewReader(r.Body)
	reader.LazyQuotes = true

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	dg, cancel := newClient()

	var resp = Response{}

	op := &api.Operation{}
	op.Schema = `
		name: string @index(exact) .
		price: string .
		type Product {
			name
			price
		}
	`
	ctx := context.Background()
	err1 := dg.Alter(ctx, op)
	if err1 != nil {
		log.Fatal(err1)
	}

	for index, row := range data {
		if index == 6 {
			break
		}

		s := strings.Split(row[0], "'")
		p := Product{
			Id:    string(s[0]),
			Name:  string(s[1]),
			Price: string(s[2]),
			DType: []string{"Product"},
		}

		mu := &api.Mutation{
			CommitNow: true,
		}
		pb, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}

		mu.SetJson = pb
		_, err2 := dg.NewTxn().Mutate(ctx, mu)
		if err2 != nil {
			log.Fatal(err2)
		}

	}
	cancel()

	resp = Response{Status: "200", Message: "products loaded"}
	return resp
}
