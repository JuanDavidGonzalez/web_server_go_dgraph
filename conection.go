package main

import (
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

type CancelFunc func()

func newClient() (*dgo.Dgraph, CancelFunc) {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	// defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

// func setup(dg *dgo.Dgraph){
// 	//Set the eschema into dgraph
// 	op := &api.Operation{}
// 	op.Schema = `
// 		name: string @index(exact) .
// 		age: int .
// 		type Buyer {
// 		name
// 		age
// 	`

// 	ctx := context.Background()
// 	err := dg.Alter(ctx, op)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
