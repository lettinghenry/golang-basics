package main

import (
    "flag"
    "fmt"
    "net/http"
    "sync"
    "time"
)

var addr = flag.String("addr",":8080","HTTP network address")

// Global broadcast system
type Broadcaster struct {
    clients map[chan string]bool
    mutex   sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
    return &Broadcaster{
        clients: make(map[chan string]bool),
    }
}

func (b *Broadcaster) AddClient(client chan string) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.clients[client] = true
}

func (b *Broadcaster) RemoveClient(client chan string) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    delete(b.clients, client)
    close(client)
}

func (b *Broadcaster) Broadcast(message string) {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    for client := range b.clients {
        select {
        case client <- message:
        default:
            // Client channel is full, skip
        }
    }
}

var broadcaster = NewBroadcaster()


func main() {
    testPrint()

    flag.Parse()

    // Start the global message broadcaster
    go startGlobalBroadcast()

    http.HandleFunc("/", home)
    http.HandleFunc("/stream", streamHandler)
    fmt.Println("Server running on", *addr)
    err := http.ListenAndServe(*addr, nil)

   if err!= nil{
        fmt.Println("Error:", err)
   }

}

func startGlobalBroadcast() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    counter := 0
    for {
        select {
        case <-ticker.C:
            counter++
            message := fmt.Sprintf("Message %d at %s", counter, time.Now().Format("15:04:05"))
            broadcaster.Broadcast(message)
        }
    }
}

func home(w http.ResponseWriter, r *http.Request){
        fmt.Fprintln(w,"Hello World!")
        fmt.Println("called /home")
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
    // Set headers for Server-Sent Events
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering
    
    // Ensure we can flush
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }

    // Create client channel and register with broadcaster
    clientChan := make(chan string, 10) // Buffer to prevent blocking
    broadcaster.AddClient(clientChan)
    defer broadcaster.RemoveClient(clientChan)

    // Send initial connection message
    fmt.Fprintf(w, "data: Connected to synchronized stream\n\n")
    flusher.Flush()

    // Create a channel to signal when client disconnects
    clientGone := r.Context().Done()

    fmt.Println("Client connected to synchronized stream")

    for {
        select {
        case <-clientGone:
            fmt.Println("Client disconnected from stream")
            return
        case message := <-clientChan:
            // Send SSE formatted data
            fmt.Fprintf(w, "data: %s\n\n", message)
            flusher.Flush()
        }
    }
}

func testPrint(){
    fmt.Println("Hello, World!")

    messages := make(chan string)

    go func ()  {messages <- "ping"}()

    msg := <- messages

    fmt.Println(msg)

}