package main

import (
	"cruzes.co/junedc/entain/sport/db"
	"cruzes.co/junedc/entain/sport/proto/sport"
	"cruzes.co/junedc/entain/sport/service"
	"database/sql"
	"flag"

	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9091", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9091")
	if err != nil {
		return err
	}

	eventDB, err := sql.Open("sqlite3", "./db/racing.db")
	if err != nil {
		return err
	}

	eventsRepo := db.NewEventsRepo(eventDB)
	if err := eventsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sport.RegisterSportServer(
		grpcServer,
		service.NewSportService(
			eventsRepo,
		),
	)

	log.Printf("gRPC server listening on: %s\n", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
