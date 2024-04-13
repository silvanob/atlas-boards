package me.biemans.atlasboards;

import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import io.grpc.StatusRuntimeException;
import me.biemans.atlasboards.api.AtlasBoardsGrpc;
import me.biemans.atlasboards.api.Ticket;
import me.biemans.atlasboards.api.TicketRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

@Component
public class AtlasBoardsClient {
    private static final Logger logger = LoggerFactory.getLogger(AtlasBoardsClient.class);

    private final AtlasBoardsGrpc.AtlasBoardsBlockingStub blockingStub;

    public AtlasBoardsClient() {
        String target = "localhost:50051";
        ManagedChannel channel = Grpc.newChannelBuilder(target, InsecureChannelCredentials.create())
                .build();
        blockingStub = AtlasBoardsGrpc.newBlockingStub(channel);
    }

    public Ticket createTicket(String title, String content) {
        logger.info("*** Create a ticket: title={}, content={}", title, content);
        Ticket ticket = Ticket.newBuilder().setTitle(title).setContent(content).build();
        Ticket ticketReceived;
        try {
            ticketReceived = blockingStub.createTicket(ticket);
            return ticketReceived;
        } catch (StatusRuntimeException e) {
            logger.warn("RPC failed: {}", e.getStatus());
        }
        return null;
    }

    public List<TicketClass> listTickets() {
        Iterator<Ticket> tickets;
        List<TicketClass> ticketList = new ArrayList<>();
        try {
            tickets = blockingStub.listTickets(TicketRequest.newBuilder().build());
            while (tickets.hasNext()) {
                Ticket ticket = tickets.next();
                ticketList.add(new TicketClass(ticket.getTitle(), ticket.getContent()));
            }
            return ticketList;
        } catch (StatusRuntimeException e) {
            logger.warn("RPC failed: {}", e.getStatus());
        }
        return null;
    }
}
