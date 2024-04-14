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
		command, args := parseCommandArgs(StringPrompt("What do you want to add to the list?"))
		switch command {
		case "add":
			cardStorage.AddCard("CUSTOM1", strings.Join(args, " "))
		case "list":
			fmt.Println(cardStorage.ListCards())
		case "remove":
			cardStorage.RemoveCardByTitle(args[1])
		default:
			fmt.Println("Bad command!")
		}
	}
}

func parseCommandArgs(commandString string) (string, []string) {
	stringSplit := strings.Split(commandString, " ")
	if len(stringSplit) == 1 {
		return stringSplit[0], nil
	} else if len(stringSplit) >= 2 {
		return stringSplit[0], stringSplit[1:]
	} else {
		return "", nil
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
