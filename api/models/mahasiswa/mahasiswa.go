package mahasiswa

import (
    "encoding/json"
    "log"
    "net/http"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../jsonhandler"
)

type Mahasiswa struct {
	IdMahasiswa	    int			`json:"idmahasiswa"`
	Nama			string		`json:"nama"`
	IdUser			int	 		`json:"iduser"`
	Departemen		string		`json:"departemen"`
	NIM				string		`json:"nim"`
    Status          string      `json:"status"`
}

func RoutesMahasiswa(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/mahasiswa"), AllMahasiswa(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addmahasiswa"), AddMahasiswa(session))
    mux.HandleFunc(pat.Get("/mahasiswa/:iduser"), GetAttributeMahasiswa(session))
}

func EnsureMahasiswa(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("ccs").C("mahasiswa")

    index := mgo.Index{
        Key:        []string{"nim"},
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

func AllMahasiswa(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("mahasiswa")

        var mahasiswa []Mahasiswa
        err := c.Find(bson.M{}).All(&mahasiswa)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all mahasiswa: ", err)
            return
        }

        respBody, err := json.MarshalIndent(mahasiswa, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddMahasiswa(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var mahasiswa Mahasiswa
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&mahasiswa)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("mahasiswa")

        //untuk auto increment
        var lastMahasiswa Mahasiswa
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastMahasiswa)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastMahasiswa.IdMahasiswa
        }
        currentId := lastId + 1
        mahasiswa.IdMahasiswa = currentId

        err = c.Insert(mahasiswa)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert mahasiswa: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+mahasiswa.NIM)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeMahasiswa(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdUser := pat.Param(r, "iduser")

        c := session.DB("ccs").C("mahasiswa")

        var mahasiswa Mahasiswa
        err := c.Find(bson.M{"iduser": IdUser}).One(&mahasiswa)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find mahasiswa: ", err)
            return
        }

        respBody, err := json.MarshalIndent(mahasiswa, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

