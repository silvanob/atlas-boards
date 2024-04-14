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
	server := newServer()
	fmt.Println(server.cards[0])
	pb.RegisterAtlasBoardsServer(grpcServer, server)

	grpcServer.Serve(lis)
	for {
		stringSplit := strings.Split(StringPrompt("What do you want to add to the list?"), " ")
		if stringSplit[0] == "add" {
			addToSlice(&server.cards, card{title: "CUSTOM1", content: strings.Join(stringSplit[1:], " ")})
			fmt.Println(server.cards)
		} else if stringSplit[0] == "list" {
			fmt.Println(server.cards)
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

type atlasBoardsServer struct {
	pb.UnimplementedAtlasBoardsServer
	cards []card
}

func (s *atlasBoardsServer) CreateTicket(ctx context.Context, ticket *pb.Ticket) (*pb.Ticket, error) {
	s.cards = append(s.cards, card{title: ticket.Title, content: ticket.Content})
	return ticket, nil
}

func (s *atlasBoardsServer) ListTickets(listTicket *pb.TicketRequest, listTicketsServer pb.AtlasBoards_ListTicketsServer) error {
	for _, carditem := range s.cards {
		if err := listTicketsServer.Send(&pb.Ticket{Title: carditem.title, Content: carditem.content}); err != nil {
			return err
		}
		fmt.Println(carditem)
	}
	return nil
}
func newServer() *atlasBoardsServer {
	s := &atlasBoardsServer{cards: []card{card{title: "Hello", content: "Hello"}}}
	return s
}
