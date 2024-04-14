package main

import (
	"context"
	"fmt"

	pb "github.com/silvanob/atlas-boards/api"
)

type atlasBoardsServer struct {
	pb.UnimplementedAtlasBoardsServer
	cardStorage *CardStorage
}

func (s *atlasBoardsServer) CreateTicket(ctx context.Context, ticket *pb.Ticket) (*pb.Ticket, error) {
	s.cardStorage.AddCard(ticket.Title, ticket.Content)
	return ticket, nil
}

func (s *atlasBoardsServer) ListTickets(listTicket *pb.TicketRequest, listTicketsServer pb.AtlasBoards_ListTicketsServer) error {
	for _, carditem := range s.cardStorage.ListCards() {
		if err := listTicketsServer.Send(&pb.Ticket{Title: carditem.title, Content: carditem.content}); err != nil {
			return err
		}
		fmt.Println(carditem)
	}
	return nil
}
func NewServer(cardStorage *CardStorage) *atlasBoardsServer {
	if cardStorage == nil {
		cardStorage = NewCardStorage()
	}
	s := &atlasBoardsServer{cardStorage: cardStorage}
	return s
}
