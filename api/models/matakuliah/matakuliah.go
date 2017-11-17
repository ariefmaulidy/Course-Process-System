package matakuliah

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type MataKuliah struct {
	IdMataKuliah	int			`json:"idmatakuliah"`
	NamaMataKuliah	string		`json:"namamatakuliah"`
    KodeMataKuliah  string      `json:"kodematakuliah"`
	IdKordinator	int	 		`json:"idkordinator"`
	IdDosen		    int		    `json:"iddosen"`
    IdPJ            int         `json:"idpj"`
	Semester		int		    `json:"semester"`
}

func RoutesMataKuliah(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/matakuliah"), AllMataKuliah(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addmatakuliah"), AddMataKuliah(session))
    mux.HandleFunc(pat.Get("/matakuliah/:idmatakuliah"), GetAttributeMataKuliah(session))
}

func EnsureMataKuliah(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("ccs").C("matakuliah")

    index := mgo.Index{
        Key:        []string{"kodematakuliah"},
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


func AllMataKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("matakuliah")

        var matakuliah []MataKuliah
        err := c.Find(bson.M{}).All(&matakuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all matakuliah: ", err)
            return
        }

        respBody, err := json.MarshalIndent(matakuliah, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddMataKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var matakuliah MataKuliah
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&matakuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("matakuliah")

        //untuk auto increment
        var lastMataKuliah MataKuliah
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastMataKuliah)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastMataKuliah.IdMataKuliah)
        }
        currentId := lastId + 1
        matakuliah.IdMataKuliah = strconv.Itoa(currentId)

        err = c.Insert(matakuliah)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert matakuliah: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+matakuliah.KodeMataKuliah)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeMataKuliah(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdMataKuliah := pat.Param(r, "idmatakuliah")

        c := session.DB("ccs").C("matakuliah")

        var matakuliah MataKuliah
        err := c.Find(bson.M{"idmatakuliah": IdMataKuliah}).One(&matakuliah)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find matakuliah: ", err)
            return
        }

        respBody, err := json.MarshalIndent(matakuliah, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

