package main

import (
    "flag"
    "fmt"
    "net/http"
)

var addr = flag.String("addr",":8080","HTTP network address")


func main() {
    testPrint()

    flag.Parse()

    http.HandleFunc("/", home)
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

func testPrint(){
    fmt.Println("Hello, World!")

    messages := make(chan string)

    go func ()  {messages <- "ping"}()

    msg := <- messages

    fmt.Println(msg)

}