package main

import (
    "goji.io"
    "gopkg.in/mgo.v2"
    "net/http"
    "github.com/rs/cors"
    "log"

    "./models/users"
    "./models/lapak"
    "./models/order"
    "./models/recentbid"
    "./socket"
    "./searchengine"
    "./models/admin"
)

func main() {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    ensureIndex(session)

    mux := goji.NewMux()
    
    users.RoutesUser(mux,session)
    lapak.RoutesLapak(mux,session)
    order.RoutesOrder(mux,session)
    recentbid.RoutesRecentBid(mux,session)
    searchengine.RoutesSearchEngine(mux,session)
    socket.RoutesSocket(mux,session)
    admin.RoutesAdmin(mux,session)
    handler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, AllowCredentials:true,}).Handler(mux)
    log.Println("Starting Listen server....")
    http.ListenAndServe("localhost:8080", handler)
}

func ensureIndex(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    users.EnsureUser(session)
    lapak.EnsureLapak(session)
    order.EnsureOrder(session)
    recentbid.EnsureRecentBid(session)
}
