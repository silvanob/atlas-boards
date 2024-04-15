package main

import (
	"context"
	"fmt"

	pb "github.com/silvanob/atlas-boards/api"
)

type AtlasBoardsServer struct {
	pb.UnimplementedAtlasBoardsServer
	cardStorage StorageDriver
}

func (s *AtlasBoardsServer) CreateTicket(ctx context.Context, ticket *pb.Ticket) (*pb.Ticket, error) {
	s.cardStorage.Add(card{title: ticket.Title, content: ticket.Content})
	return ticket, nil
}

func (s *AtlasBoardsServer) ListTickets(listTicket *pb.TicketRequest, listTicketsServer pb.AtlasBoards_ListTicketsServer) error {
	for _, carditem := range s.cardStorage.List() {
		if err := listTicketsServer.Send(&pb.Ticket{Title: carditem.title, Content: carditem.content}); err != nil {
			return err
		}
		fmt.Println(carditem)
	}
	return nil
}
func NewServer(cardStorage StorageDriver) *AtlasBoardsServer {
	if cardStorage == nil {
		cardStorage = NewCardStorageSlice()
	}
	s := &AtlasBoardsServer{cardStorage: cardStorage}
	return s
}
