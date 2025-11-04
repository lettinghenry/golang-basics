package main

import (
    "flag"
    "fmt"
    "net/http"
    "time"
)

var addr = flag.String("addr",":8080","HTTP network address")


func main() {
    testPrint()

    flag.Parse()

    http.HandleFunc("/", home)
    http.HandleFunc("/stream", streamHandler)
    fmt.Println("Server running on", *addr)
    err := http.ListenAndServe(*addr, nil)

   if err!= nil{
        fmt.Println("Error:", err)
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

    // Create a channel to signal when client disconnects
    clientGone := r.Context().Done()

    fmt.Println("Client connected to stream")

    // Send data every second
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    counter := 0
    for {
        select {
        case <-clientGone:
            fmt.Println("Client disconnected from stream")
            return
        case <-ticker.C:
            counter++
            // Send SSE formatted data
            fmt.Fprintf(w, "data: Message %d at %s\n\n", counter, time.Now().Format("15:04:05"))
            
            // Flush the data to client immediately
            if flusher, ok := w.(http.Flusher); ok {
                flusher.Flush()
            }
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