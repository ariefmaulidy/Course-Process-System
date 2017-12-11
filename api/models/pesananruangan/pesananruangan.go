package pesananruangan

import (
    "encoding/json"
    "log"
    "time"
    "net/http"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../jsonhandler"
)

type PesananRuangan struct {
	IdPesanan       int         `json:"idpesanan"`
    IdRuangan	    int			`json:"idruangan"`
    IdMataKuliah    int         `json:"idmatakuliah"`
    Tanggal         time.Time   `json:"tanggal"`
    Waktu           string      `json:"waktu"`
    Status          string      `json:"status"`
    Alasan          string      `json:"alasan"`
    IdPemesan       string      `json:"idpemesan"`
    ConfirmedBy     int         `json:"confirmedby"`
}

func RoutesPesananRuangan(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/pesananruangan"), AllPesananRuangan(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpesananruangan"), AddPesananRuangan(session))
    mux.HandleFunc(pat.Get("/pesananruangan/:idpesanan"), GetAttributePesananRuangan(session))
}

func AllPesananRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pesananruangan")

        var pesananruangan []PesananRuangan
        err := c.Find(bson.M{}).All(&pesananruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pesananruangan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesananruangan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPesananRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var pesananruangan PesananRuangan
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pesananruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("pesananruangan")

        //untuk auto increment
        var lastPesananRuangan PesananRuangan
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastPesananRuangan)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastPesananRuangan.IdPesanan
        }
        currentId := lastId + 1
        pesananruangan.IdPesanan = currentId

        err = c.Insert(pesananruangan)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pesananruangan: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributePesananRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdPesanan := pat.Param(r, "idpesanan")

        c := session.DB("ccs").C("pesananruangan")

        var pesananruangan PesananRuangan
        err := c.Find(bson.M{"idpesanan": IdPesanan}).One(&pesananruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pesananruangan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesananruangan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

