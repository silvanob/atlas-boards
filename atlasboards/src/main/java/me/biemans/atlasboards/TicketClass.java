package me.biemans.atlasboards;

public class TicketClass {
    private String title;
    private String content;

    public TicketClass(String title, String content) {
        this.title = title;
        this.content = content;
    }

    public TicketClass() {}

    public String getContent() {
        return content;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public void setContent(String content) {
        this.content = content;
    }
}
