package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"context"
	"flag"
	"log"
	"net"
	"sync"

	pb "github.com/silvanob/atlas-boards/routeguide"
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
	channel := make(chan card)
	slices := []card{card{title: "Hello", content: "Hello"}}
	fmt.Println("This will be a board for managing tasks")
	fmt.Println(slices[0])

	grpcServer := grpc.NewServer()
	server := newServer(&channel)
	pb.RegisterRouteGuideServer(grpcServer, server)
	go grpcServer.Serve(lis)
	go func() {
		for x := range channel {
			addToSlice(&slices, x)
			fmt.Println(slices)
		}
	}()
	for {
		stringSplit := strings.Split(StringPrompt("What do you want to add to the list?"), " ")
		if stringSplit[0] == "add" {
			server.channel <- card{title: "CUSTOM1", content: strings.Join(stringSplit[1:], " ")}
		} else {
			fmt.Println("Bad command!")
		}
	}

}

func addToSlice(slice *[]card, stringu card) {
	*slice = append(*slice, stringu)
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

type card struct {
	title   string
	content string
}

type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature // read-only after initialized
	channel       chan card
	mu            sync.Mutex // protects routeNotes
	routeNotes    map[string][]*pb.RouteNote
}

func (s *routeGuideServer) CreateTicket(ctx context.Context, ticket *pb.Ticket) (*pb.Ticket, error) {
	s.channel <- card{title: ticket.Title, content: ticket.Content}

	return ticket, nil
}
func newServer(channel *chan card) *routeGuideServer {
	s := &routeGuideServer{channel: *channel, routeNotes: make(map[string][]*pb.RouteNote)}
	return s
}
