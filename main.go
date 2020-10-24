package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	product "github.com/anatollupacescu/atlant/internal"
	"github.com/anatollupacescu/atlant/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var (
	products *mongo.Collection
	mongoCtx context.Context
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mongoCtx = context.Background()

	mongoURL, exists := os.LookupEnv("MONGODB_URL")

	if !exists { // default for local run
		mongoURL = "mongodb://localhost:27017"
	}

	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return fmt.Errorf("connect to mongo: %v", err)
	}

	log.Printf("ping mongo at URL: %v\n", mongoURL)

	if err := db.Ping(mongoCtx, nil); err != nil {
		return fmt.Errorf("ping mongo: %v", err)
	}

	products = db.Database("price").Collection("product")

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("unable to listen on port :50051: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterProductServiceServer(s, &ProductServiceServer{
		App: product.App{ProductDB: products},
	})

	go func() {
		log.Println("starting server on port :50051...")

		if err := s.Serve(listener); err != nil {
			log.Fatalf("serve rpc: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	log.Println("\nstopping the server...")

	s.Stop()

	if err := listener.Close(); err != nil {
		return err
	}

	log.Println("closing db connection")

	err = db.Disconnect(mongoCtx)

	log.Println("done")

	return err
}
