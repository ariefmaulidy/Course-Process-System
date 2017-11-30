package jsonhandler

import (
    "fmt"
    "net/http"
)

type MessageJSON struct{
    Status  bool    `json:"status"`
    Message string    `json:"message"`
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}

func SendWithJSON(w http.ResponseWriter, message string, code int)  {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    fmt.Fprintf(w, message)
}


