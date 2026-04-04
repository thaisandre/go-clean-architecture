package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/20-CleanArch/configs"
	"github.com/devfullcycle/20-CleanArch/internal/event/handler"
	"github.com/devfullcycle/20-CleanArch/internal/infra/graph"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/service"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web/webserver"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName)
	db := connectMySQL(dsn)
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs.AMQPUrl)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	// REST
	ws := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	ws.AddHandler("POST", "/order", webOrderHandler.Create)
	ws.AddHandler("GET", "/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go ws.Start()

	// gRPC
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// GraphQL
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func connectMySQL(dsn string) *sql.DB {
	var db *sql.DB
	var err error
	for i := range 30 {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Println("Connected to MySQL")
				return db
			}
		}
		log.Printf("Waiting for MySQL... attempt %d/30: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("Could not connect to MySQL after 30 attempts: %v", err)
	return nil
}


func getRabbitMQChannel(amqpUrl string) *amqp.Channel {
	var conn *amqp.Connection
	var err error
	for i := range 30 {
		conn, err = amqp.Dial(amqpUrl)
		if err == nil {
			break
		}
		log.Printf("Waiting for RabbitMQ... attempt %d/30: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ after 30 attempts: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
