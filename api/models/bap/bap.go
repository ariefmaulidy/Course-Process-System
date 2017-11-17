package bap

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type BAP struct {
	IdBAP	        int			`json:"idcgd"`
	Tanggal	        time.Time   `json:"tanggal"`
    IdMataKuliah    int         `json:"idmatakuliah"`
    TopikKuliah     string      `json:"topikkuliah"`
    JumlahPeserta   int         `json:"jumlahpeserta"`
    IdDosen         int         `json:"iddosen"`
    CreatedBy       int         `json:"createdby"`
    UpdatedBy       int         `json:"updatedby"`
}

func RoutesBAP(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/bap"), AllBAP(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addbap"), AddBAP(session))
    mux.HandleFunc(pat.Get("/bap/:idbap"), GetAttributeBAP(session))
}

func AllBAP(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("bap")

        var bap []BAP
        err := c.Find(bson.M{}).All(&bap)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all bap: ", err)
            return
        }

        respBody, err := json.MarshalIndent(bap, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddBAP(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var bap BAP
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&bap)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("bap")

        //untuk auto increment
        var lastBAP BAP
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastBAP)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastBAP.IdBAP)
        }
        currentId := lastId + 1
        bap.IdBAP = strconv.Itoa(currentId)

        err = c.Insert(bap)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert bap: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+bap.IdBAP)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeBAP(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdBAP := pat.Param(r, "idbap")

        c := session.DB("ccs").C("bap")

        var bap BAP
        err := c.Find(bson.M{"idbap": IdBAP}).One(&bap)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find bap: ", err)
            return
        }

        respBody, err := json.MarshalIndent(bap, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

