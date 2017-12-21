package pesananruangan

import (
    "encoding/json"
    "log"
    "net/http"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../jsonhandler"
    "../../auth"
    "../ruangan"
    "../matakuliah"
)

type PesananRuangan struct {
	IdPesanan       int         `json:"idpesanan"`
    IdRuangan	    int			`json:"idruangan"`
    IdMataKuliah    int         `json:"idmatakuliah"`
    Tanggal         string      `json:"tanggal"`
    Waktu           string      `json:"waktu"`
    Status          string      `json:"status"`
    Alasan          string      `json:"alasan"`
    IdPemesan       int         `json:"idpemesan"`
    ConfirmedBy     int         `json:"confirmedby"`
}

type DataSend struct{
    DataPesanan []PesananRuangan        `json:"datapesanan"`
    DataRuangan []ruangan.Ruangan       `json:"dataruangan"`
    DataMatkul  []matakuliah.MataKuliah `json:"datamatkul"`
}

func RoutesPesananRuangan(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/pesananruangan"), AllPesananRuangan(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Get("/mypesananruangan"), auth.Validate(MyPesananRuangan(session)))
    mux.HandleFunc(pat.Post("/addpesananruangan"), auth.Validate(AddPesananRuangan(session)))
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
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
        }

        session := s.Copy()
        defer session.Close()

        var pesananruangan PesananRuangan
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pesananruangan)
        print(err)
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
        pesananruangan.IdPemesan = claims.IdUser

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

func MyPesananRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
        }
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pesananruangan")
        d := session.DB("ccs").C("ruangan")
        e := session.DB("ccs").C("matakuliah")

        var datasend DataSend
        err := c.Find(bson.M{"idpemesan":claims.IdUser}).All(&datasend.DataPesanan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pesananruangan: ", err)
            return
        }

        var matkul matakuliah.MataKuliah
        var ruangan ruangan.Ruangan
        for _,data := range datasend.DataPesanan {
            e.Find(bson.M{"idmatakuliah": data.IdMataKuliah}).One(&matkul)
            datasend.DataMatkul = append(datasend.DataMatkul, matkul)
            d.Find(bson.M{"idruangan": data.IdRuangan}).One(&ruangan)
            datasend.DataRuangan = append(datasend.DataRuangan, ruangan)            
        }

        respBody, err := json.MarshalIndent(datasend, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)

    }
}