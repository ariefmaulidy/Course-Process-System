package dosen

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
    "../matakuliah"
)

type Dosen struct {
	Nama			string		`json:"nama"`
	IdUser			int	 		`json:"iduser"`
	Departemen		string		`json:"departemen"`
    Fakultas        string      `json:"fakultas"`
	NIDN			string		`json:"nidn"`
}

func RoutesDosen(mux *goji.Mux, session *mgo.Session) {
    mux.HandleFunc(pat.Get("/dosen"), AllDosen(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/adddosen"), AddDosen(session))
    mux.HandleFunc(pat.Get("/dosen/:iduser"), GetAttributeDosen(session))
}


func AllDosen(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("dosen")

        var dosen []Dosen
        err := c.Find(bson.M{}).All(&dosen)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all dosen: ", err)
            return
        }

        respBody, err := json.MarshalIndent(dosen, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddDosen(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var dosen Dosen
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&dosen)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("dosen")

        //untuk auto increment
        var lastDosen Dosen
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastDosen)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastDosen.IdDosen
        }
        currentId := lastId + 1
        dosen.IdDosen = currentId

        err = c.Insert(dosen)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert dosen: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+dosen.NIDN)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeDosen(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdUser := pat.Param(r, "iduser")

        c := session.DB("ccs").C("dosen")

        var dosen Dosen
        err := c.Find(bson.M{"iduser": IdUser}).One(&dosen)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find dosen: ", err)
            return
        }

        respBody, err := json.MarshalIndent(dosen, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}