package me.biemans.atlasboards;

import me.biemans.atlasboards.api.Ticket;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping
public class MainApiController {
    private final AtlasBoardsClient atlasBoardsClient;

    public MainApiController(AtlasBoardsClient atlasBoardsClient) {
        this.atlasBoardsClient = atlasBoardsClient;
    }

    @GetMapping("")
    public String getHelloWorld() {
            Ticket ticket = atlasBoardsClient.createTicket("general", "world");
            if (ticket != null) {
                return "title=" + ticket.getTitle() + " content=" + ticket.getContent();
            }
        return "Hello!";
    }

    @PostMapping
    public ResponseEntity<?> createTicket(@RequestBody TicketClass ticketClass) {
        Ticket ticket = atlasBoardsClient.createTicket(ticketClass.getTitle(), ticketClass.getContent());
        if (ticket != null) {
            TicketClass ticketCreated = new TicketClass(ticket.getTitle(), ticket.getContent());
            return ResponseEntity.ok().body(ticketCreated);
        }

        return ResponseEntity.internalServerError().build();
    }

    @GetMapping("/list")
    public List<TicketClass> getTickets() {
        return atlasBoardsClient.listTickets();
    }
}
