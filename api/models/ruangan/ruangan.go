package ruangan

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type Ruangan struct {
	IdRuangan	    int			`json:"idruangan"`
	NamaRuangan	    string		`json:"namaruangan"`
    Kelengkapan     string      `json:"kelengkapan"`
}

func RoutesRuangan(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/ruangan"), AllRuangan(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addruangan"), AddRuangan(session))
    mux.HandleFunc(pat.Get("/ruangan/:idruangan"), GetAttributeRuangan(session))
}

func EnsureRuangan(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("ccs").C("ruangan")

    index := mgo.Index{
        Key:        []string{"idruangan"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }

    err := c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }
}


func AllRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("ruangan")

        var ruangan []Ruangan
        err := c.Find(bson.M{}).All(&ruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all ruangan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(ruangan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var ruangan Ruangan
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&ruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("ruangan")

        //untuk auto increment
        var lastRuangan Ruangan
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastRuangan)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastRuangan.IdRuangan)
        }
        currentId := lastId + 1
        ruangan.IdRuangan = strconv.Itoa(currentId)

        err = c.Insert(ruangan)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert ruangan: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+ruangan.NamaRuangan)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdRuangan := pat.Param(r, "idruangan")

        c := session.DB("ccs").C("ruangan")

        var ruangan Ruangan
        err := c.Find(bson.M{"idruangan": IdRuangan}).One(&ruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find ruangan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(ruangan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}