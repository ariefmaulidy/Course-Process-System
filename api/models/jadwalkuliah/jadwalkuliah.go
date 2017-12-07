package jadwalkuliah

import (
    "encoding/json"
    "log"
    "time"
    "net/http"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../pesertakuliah"
    "../ruangan"
    "../../auth"
    "../../jsonhandler"
)

type JadwalKuliah struct {
	IdJadwalKuliah	int			`json:"idjadwalkuliah"`
    IdPJ            int         `json:"idpj"`
	IdMataKuliah	int	      	`json:"idmatakuliah"`
    IdRuangan       int         `json:"idruangan"`
    WaktuMulai		time.Time   `json:"waktumulai"`
    WaktuSelesai    time.Time   `json:"waktuselesai"`
    Hari            string      `json:"hari"`
}

func RoutesJadwalKuliah(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/jadwalkuliah"), AllJadwalKuliah(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addjadwalkuliah"), AddJadwalKuliah(session))
    mux.HandleFunc(pat.Get("/jadwalkuliah/:idjadwalkuliah"), GetAttributeJadwalKuliah(session))
    mux.HandleFunc(pat.Get("/detailjadwalkuliah/:idjadwalkuliah"), auth.Validate(GetDetailJadwalKuliah(session))) //jadwal yang dilihat di TU terdapat list mahasiswa
}


func AllJadwalKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("jadwalkuliah")

        var jadwalkuliah []JadwalKuliah
        err := c.Find(bson.M{}).All(&jadwalkuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all jadwalkuliah: ", err)
            return
        }

        respBody, err := json.MarshalIndent(jadwalkuliah, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddJadwalKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var jadwalkuliah JadwalKuliah
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&jadwalkuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("jadwalkuliah")

        //untuk auto increment
        var lastJadwalKuliah JadwalKuliah
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastJadwalKuliah)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastJadwalKuliah.IdJadwalKuliah
        }
        currentId := lastId + 1
        jadwalkuliah.IdJadwalKuliah = currentId

        err = c.Insert(jadwalkuliah)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert jadwalkuliah: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeJadwalKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdJadwalKuliah := pat.Param(r, "idjadwalkuliah")

        c := session.DB("ccs").C("jadwalkuliah")

        var jadwalkuliah JadwalKuliah
        err := c.Find(bson.M{"idjadwalkuliah": IdJadwalKuliah}).One(&jadwalkuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find jadwalkuliah: ", err)
            return
        }

        respBody, err := json.MarshalIndent(jadwalkuliah, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func GetDetailJadwalKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdJadwalKuliah := pat.Param(r, "idjadwalkuliah")

        c := session.DB("ccs").C("jadwalkuliah")
        d := session.DB("ccs").C("ruangan")
        e := session.DB("ccs").C("pesertakuliah")

        type DataSend struct {
            DataJadwalKuliah        JadwalKuliah                     `json:"datajadwalkuliah"`
            DataRuangan             ruangan.Ruangan                  `json:"dataruangan"`
            DataPesertaKuliah       []pesertakuliah.PesertaKuliah    `json:"datapesertakuliah"`
        }

        var datasend DataSend

        err := c.Find(bson.M{"idjadwalkuliah": IdJadwalKuliah}).One(&datasend.DataJadwalKuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find jadwalkuliah: ", err)
            return
        }

        err = d.Find(bson.M{"idruangan": datasend.DataJadwalKuliah.IdRuangan}).One(&datasend.DataRuangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find dataruangan: ", err)
            return
        }

        err = e.Find(bson.M{"idjadwalkuliah": IdJadwalKuliah}).All(&datasend.DataPesertaKuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find datapesertakuliah: ", err)
            return
        }

        

        respBody, err := json.MarshalIndent(datasend, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}
