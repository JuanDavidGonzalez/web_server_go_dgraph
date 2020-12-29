package main

import (
	"bytes"
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

	(w).Header().Set("Access-Control-Allow-Origin", "*")

	dg, cancel := newClient()

	q := `query Data{
		data(func: type(Product)) {
			uid
			id
			name
			price
		}
	}`

	ctx := context.Background()
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}
	type Root struct {
		Data []Product `json:"data"`
	}
	var root Root
	err = json.Unmarshal(resp.Json, &root)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}
	cancel()
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func ListBuyers(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")

	dg, cancel := newClient()

	q := `query Data{
		data(func: type(Buyer)) {
			uid
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
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}
	cancel()
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func BuyerDetail(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")

	ids, exist := r.URL.Query()["id"]
	if !exist || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id invalid")
		return
	}
	id := ids[0]

	dg, cancel := newClient()

	bid := make(map[string]string)
	bid["$id"] = id

	q := `query Data($id: string){
		data(func: type(Transaction))@filter(eq(buyer,$id)){
			buyer
			device
			ip
			products
		}
	}`

	ctx := context.Background()
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, bid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}
	type Root struct {
		Data []Transaction `json:"data"`
	}
	var root Root
	err = json.Unmarshal(resp.Json, &root)
	if err != nil {
		log.Fatal(err)
	}
	cancel()
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")

	ids, exist := r.URL.Query()["id"]
	if !exist || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id invalid")
		return
	}
	id := ids[0]

	dg, cancel := newClient()

	bid := make(map[string]string)
	bid["$id"] = id

	q := `query Data($id: string){
		data(func: type(Product))@filter(eq(id,$id)){
			id
			name
			price
		}
	}`

	ctx := context.Background()
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, bid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}
	type Root struct {
		Data []Transaction `json:"data"`
	}
	var root Root
	err = json.Unmarshal(resp.Json, &root)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}
	cancel()
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func GetOtherBuyer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	ips, exist := r.URL.Query()["ip"]
	if !exist || len(ips[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ip invalid")
		return
	}
	ip := ips[0]

	dg, cancel := newClient()

	bid := make(map[string]string)
	bid["$ip"] = ip

	q := `query Data($ip: string){
			var(func:eq(ip,$ip)){
			buyerid as buyer
		  }
		data(func: type(Buyer))@filter(eq(id,val(buyerid))){
			name
			age
			id
		}
	}`

	ctx := context.Background()
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, bid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}

	type Root struct {
		Data []Buyer `json:"data"`
	}
	var root Root
	err = json.Unmarshal(resp.Json, &root)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error")
		return
	}

	cancel()
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
}

func LoadData(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")

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
	} else if typeData == "transactions" {
		result = loadTransactions(response)
	} else {
		result = Response{Status: "400", Message: "invalid type"}
	}

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
		id: string @index(exact) .
		age: int .
		type Buyer {
			name
			id
			age
		}
	`
	ctx := context.Background()
	err = dg.Alter(ctx, op)
	if err != nil {
		fmt.Println(err)
		resp = Response{Status: "500", Message: "internal error"}
		return resp
	}
	var responseObject BuyerReponse
	json.Unmarshal(responseData, &responseObject)

	for _, item := range responseObject {
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
			fmt.Println(err)
			resp = Response{Status: "500", Message: "internal error"}
			return resp
		}

		mu.SetJson = pb
		_, err = dg.NewTxn().Mutate(ctx, mu)
		if err != nil {
			log.Fatal(err)
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
		id: string @index(exact) .
		price: string .
		type Product {
			id
			name
			price
		}
	`
	ctx := context.Background()
	err = dg.Alter(ctx, op)
	if err != nil {
		fmt.Println(err)
		resp = Response{Status: "500", Message: "internal error"}
		return resp
	}

	for _, row := range data {
		quoted := false
		s := strings.FieldsFunc(row[0], func(r rune) bool {
			if r == '"' {
				quoted = !quoted
			}
			return !quoted && r == '\''
		})
		// s := strings.Split(row[0], "'")
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
			log.Fatal("aqui...")
			log.Fatal(err)
		}

		mu.SetJson = pb
		_, err = dg.NewTxn().Mutate(ctx, mu)
		if err != nil {
			fmt.Println(err)
			resp = Response{Status: "500", Message: "internal error"}
			return resp
		}

	}
	cancel()

	resp = Response{Status: "200", Message: "products loaded"}
	return resp
}

func loadTransactions(r *http.Response) Response {

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := bytes.Split(responseData, []byte{0, 0})

	var resp = Response{}

	dg, cancel := newClient()

	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		ip: string @index(exact).
		device: string .
		buyer: string .
		products: [string] .
		type Transaction {
			id
			buyer
			device
			ip
			products
		}
	`
	ctx := context.Background()
	err = dg.Alter(ctx, op)
	if err != nil {
		fmt.Println(err)
		resp = Response{Status: "500", Message: "internal error"}
		return resp
	}
	if len(data) > 0 {
		for index, item := range data {
			if index > 100 {
				break
			}
			element := bytes.Split(item, []byte{0})

			prod_str := string(element[4])
			prod_trim := strings.Trim(prod_str, "()")
			products := strings.Split(prod_trim, ",")

			t := Transaction{
				Id:       string(element[0]),
				Buyer:    string(element[1]), //Buyer{Uid: string(element[1])},
				Ip:       string(element[2]),
				Device:   string(element[3]),
				Products: products,
				DType:    []string{"Transaction"},
			}

			mu := &api.Mutation{
				CommitNow: true,
			}
			pb, err := json.Marshal(t)
			if err != nil {
				log.Fatal(err)
			}

			mu.SetJson = pb
			_, err = dg.NewTxn().Mutate(ctx, mu)
			if err != nil {
				fmt.Println(err)
				resp = Response{Status: "500", Message: "internal error"}
				return resp
			}
		}
		resp = Response{Status: "200", Message: "transactions loaded"}
	} else {
		resp = Response{Status: "200", Message: "transaction data not found"}
	}
	cancel()

	return resp
}
