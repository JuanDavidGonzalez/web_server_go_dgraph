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
	Uid        string   `json:"uid,omitempty"`
	Id         string   `json:"id,omitempty"`
	BuyerId    string   `json:"buyerid,omitempty"`
	Ip         string   `json:"ip,omitempty"`
	Device     string   `json:"device,omitempty"`
	Productids []string `json:"productids,omitempty"`
	DType      []string `json:"dgraph.type,omitempty"`
}

// type BuyerList struct {
// 	Data []Buyer `json:"data"`
// }

type BuyerReponse []Buyer

// type ProductResponse []Product

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// func (p *Product) ToJson() ([]byte, error) {
// 	return json.Marshal(p)
// }
