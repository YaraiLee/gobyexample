package main

import (
    "net/http"
    "time"
)

func helloTimeout(res http.ResponseWriter, req *http.Request) {
    time.Sleep(5*time.Second)
    res.Header().Set("Content-Type", "application/json")
    res.Write([]byte("Hello web\n"))
}

func main() {
    http.HandleFunc("/", http.TimeoutHandler(http.HandlerFunc(helloTimeout), 3*time.Second, "timeout").ServeHTTP)
    http.ListenAndServe(":5000", nil)
}
