package pesan

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type Pesan struct {
	IdPesan	        int			`json:"idpesan"`
	IsiPesan	    string		`json:"isipesan"`
}

func RoutesPesan(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/pesan"), AllPesan(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpesan"), AddPesan(session))
    mux.HandleFunc(pat.Get("/pesan/:idpesan"), GetAttributePesan(session))
}

func AllPesan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pesan")

        var pesan []Pesan
        err := c.Find(bson.M{}).All(&pesan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pesan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPesan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var pesan Pesan
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pesan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("pesan")

        //untuk auto increment
        var lastPesan Pesan
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastPesan)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastPesan.IdPesan)
        }
        currentId := lastId + 1
        pesan.IdPesan = strconv.Itoa(currentId)

        err = c.Insert(pesan)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pesan: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+pesan.IdPesan)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetPesan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdPesan := pat.Param(r, "idpesan")

        c := session.DB("ccs").C("pesan")

        var pesan Pesan
        err := c.Find(bson.M{"idpesan": IdPesan}).One(&pesan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pesan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

