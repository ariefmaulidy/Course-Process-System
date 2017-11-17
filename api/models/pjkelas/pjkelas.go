package pjkelas

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type PJKelas struct {
	IdPJ			int			`json:"idpj"`
	Nama			string		`json:"nama"`
	IdUser			int	 		`json:"iduser"`
	Departemen		string		`json:"departemen"`
	NIM				string		`json:"nim"`
}

func RoutesPJKelas(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/pjkelas"), AllPJKelas(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpjkelas"), AddPJKelas(session))
    mux.HandleFunc(pat.Get("/pjkelas/:iduser"), GetAttributePJKelas(session))
}

func EnsurePJKelas(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("ccs").C("pjkelas")

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

func AllPJKelas(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pjkelas")

        var pjkelas []PJKelas
        err := c.Find(bson.M{}).All(&pjkelas)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pjkelas: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pjkelas, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPJKelas(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var pjkelas PJKelas
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pjkelas)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("pjkelas")

        //untuk auto increment
        var lastPJKelas PJKelas
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastPJKelas)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastPJKelas.IdPJ)
        }
        currentId := lastId + 1
        pjkelas.IdPJ = strconv.Itoa(currentId)

        err = c.Insert(pjkelas)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pjkelas: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+pjkelas.NIM)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributePJKelas(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdUser := pat.Param(r, "iduser")

        c := session.DB("ccs").C("pjkelas")

        var pjkelas PJKelas
        err := c.Find(bson.M{"iduser": IdUser}).One(&pjkelas)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pjkelas: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pjkelas, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

