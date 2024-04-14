package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"flag"
	"log"
	"net"

	pb "github.com/silvanob/atlas-boards/api"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("This will be a board for managing tasks")
	grpcServer := grpc.NewServer()
	cardStorage := NewCardStorage()
	server := NewServer(cardStorage)
	pb.RegisterAtlasBoardsServer(grpcServer, server)

	go grpcServer.Serve(lis)
	for {
		stringSplit := strings.Split(StringPrompt("What do you want to add to the list?"), " ")
		if stringSplit[0] == "add" {
			cardStorage.AddCard("CUSTOM1", strings.Join(stringSplit[1:], " "))
		} else if stringSplit[0] == "list" {
			fmt.Println(cardStorage.ListCards())
		} else if stringSplit[0] == "remove" {
			cardStorage.RemoveCardByTitle(stringSplit[1])
		} else {
			fmt.Println("Bad command!")
		}
	}
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
