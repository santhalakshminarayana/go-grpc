package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-grpc/serv-a/config"
	"github.com/go-grpc/serv-a/server"
)

var (
	host        = flag.String("host", "0.0.0.0", "Host Address")
	port        = flag.Int("port", 4010, "Host listening Port")
	envFilepath = flag.String("envpath", ".env", "ENV file path")
)

func main() {
	flag.Parse()

	// Init config
	config.InitConfig(*envFilepath)

	config := config.GetConfig()

	log.Printf("Starting server %v ...", config.ServiceName)

	server.StartServer(*host, *port)
	os.Exit(0)
}
