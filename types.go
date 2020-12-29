package main

type Product struct {
	Uid   string   `json:"uid,omitempty"`
	Id    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Price string   `json:"price,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Buyer struct {
	Uid   string   `json:"uid,omitempty"`
	Id    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Age   int      `json:"age,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Transaction struct {
	Uid      string   `json:"uid,omitempty"`
	id       string   `json:"id,omitempty"`
	Id       string   `json:"id,omitempty"`
	Buyer    string   `json:"buyer,omitempty"`
	Ip       string   `json:"ip,omitempty"`
	Device   string   `json:"device,omitempty"`
	Products []string `json:"products,omitempty"`
	DType    []string `json:"dgraph.type,omitempty"`
}

type BuyerReponse []Buyer

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
