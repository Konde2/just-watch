package main

import (
	"log"
	"net"
	"os"

	"metadata-service/handlers"
	"metadata-service/pb"
	"metadata-service/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(nil)

	repo := repository.NewVideoRepository(client)
	service := handlers.NewMetadataService(repo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterMetadataServiceServer(s, service)
	reflection.Register(s) // для grpcurl

	log.Println("Metadata Service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
