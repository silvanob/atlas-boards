package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"flag"
	"log"
	"net"

	pb "github.com/silvanob/atlas-boards/api"
	"golang.org/x/term"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	if !term.IsTerminal(int(syscall.Stdin)) {
		fmt.Println("Terminal is not interactive, exiting!")
		return
	}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("This will be a board for managing tasks")
	grpcServer := grpc.NewServer()
	cardStorage := NewCardStorageMap()
	server := NewServer(cardStorage)
	pb.RegisterAtlasBoardsServer(grpcServer, server)

	go grpcServer.Serve(lis)
	for {
		command, args := parseCommandArgs(StringPrompt("What do you want to add to the list?"))
		switch command {
		case "add":
			cardStorage.Add(card{title: "CUSTOM1", content: strings.Join(args, " ")})
		case "list":
			fmt.Println(cardStorage.List())
		case "remove":
			cardStorage.Remove(args[0])
		default:
			fmt.Println("Bad command!")
		}
	}
}

func parseCommandArgs(commandString string) (string, []string) {
	stringSplit := strings.Split(commandString, " ")
	if stringLength := len(stringSplit); stringLength == 1 {
		return stringSplit[0], nil
	} else if stringLength >= 2 {
		return stringSplit[0], stringSplit[1:]
	}
	return "", nil
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
