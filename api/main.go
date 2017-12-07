package main

import (
    "goji.io"
    "gopkg.in/mgo.v2"
    "net/http"
    "github.com/rs/cors"
    "log"

    "./models/bap"
    "./models/chatgroupdiscussion"
    "./models/dosen"
    "./models/jadwalkuliah"
    "./models/mahasiswa"
    "./models/matakuliah"
    "./models/pengajar"
    "./models/pengelolaruangan"
    "./models/pesan"
    "./models/pesananruangan"
    "./models/pesertakuliah"
    "./models/ruangan"
    "./models/tatausaha"
    "./models/user"
    "./socket"
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
    
    bap.RoutesBAP(mux,session)
    chatgroupdiscussion.RoutesChatGroupDiscussion(mux,session)
    dosen.RoutesDosen(mux,session)
    jadwalkuliah.RoutesJadwalKuliah(mux,session)
    mahasiswa.RoutesMahasiswa(mux,session)
    matakuliah.RoutesMataKuliah(mux,session)
    pengajar.RoutesPengajar(mux,session)
    pengelolaruangan.RoutesPengelolaRuangan(mux,session)
    pesan.RoutesPesan(mux,session)
    pesananruangan.RoutesPesananRuangan(mux,session)
    pesertakuliah.RoutesPesertaKuliah(mux,session)
    ruangan.RoutesRuangan(mux,session)
    tatausaha.RoutesTataUsaha(mux,session)
    user.RoutesUser(mux,session)
    socket.RoutesSocket(mux,session)
    handler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, AllowCredentials:true,}).Handler(mux)
    log.Println("Starting Listen server....")
    http.ListenAndServe("localhost:8080", handler)
}

func ensureIndex(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()
    user.EnsureUser(session)
    matakuliah.EnsureMataKuliah(session)
}
