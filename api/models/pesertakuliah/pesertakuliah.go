package pesertakuliah

import (


    "goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type PesertaKuliah struct {
    IdPesertaKuliah int         `json:"idpesertakuliah"`
    IdJadwalKuliah  int         `json:"idjadwalkuliah"`
    IdUser          int         `json:"iduser"`
}

func RoutesPesertaKuliah(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/pesertakuliah"), AllPesertaKuliah(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpesertakuliah"), AddPesertaKuliah(session))
    mux.HandleFunc(pat.Get("/pesertakuliah/:idpesertakuliah"), GetAttributePesertaKuliah(session))
}


func AllPesertaKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pesertakuliah")

        var pesertakuliah []PesertaKuliah
        err := c.Find(bson.M{}).All(&pesertakuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pesertakuliah: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesertakuliah, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPesertaKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var pesertakuliah PesertaKuliah
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pesertakuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("pesertakuliah")

        //untuk auto increment
        var lastPesertaKuliah PesertaKuliah
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastPesertaKuliah)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastPesertaKuliah.IdPesertaKuliah)
        }
        currentId := lastId + 1
        pesertakuliah.IdPesertaKuliah = strconv.Itoa(currentId)

        err = c.Insert(pesertakuliah)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pesertakuliah: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+pesertakuliah.NIM)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributePesertaKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdPesertaKuliah := pat.Param(r, "idpesertakuliah")

        c := session.DB("ccs").C("pesertakuliah")

        var pesertakuliah PesertaKuliah
        err := c.Find(bson.M{"idpesertakuliah": IdPesertaKuliah}).One(&pesertakuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pesertakuliah: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesertakuliah, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}



