package tatausaha

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../mahasiswa"
    "../../auth"
    "../../jsonhandler"
)

type Tatausaha struct {
	IdTU		int    	`json:"idtu"`
	Nama		string 	`json:"nama"`
	IdUser      int     `json:"iduser"`
	Departemen	string	`json:"departemen"`
	Fakultas	string	`json:"fakultas"`
}

func RoutesTataUsaha(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/tatausaha"), AllTataUsaha(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addtatausaha"), AddTataUsaha(session))
    mux.HandleFunc(pat.Get("/tatausaha/:iduser"), GetAttributeTataUsaha(session))
    //untuk memasukkan BAP cek di fungsi AddBAP di model bap
    mux.HandleFunc(pat.Put("/assignpjkelas"), auth.Validate(AssignPJKelas(session)))
    //untuk book ruangan cek di fungsi menambahPesananRuangan  di model pesanan ruangan
}

func AllTataUsaha(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("tatausaha")

        var tatausaha []Tatausaha
        err := c.Find(bson.M{}).All(&tatausaha)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all tatausaha: ", err)
            return
        }

        respBody, err := json.MarshalIndent(tatausaha, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddTataUsaha(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var tatausaha Tatausaha
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&tatausaha)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("tatausaha")

        //untuk auto increment
        var lastTataUsaha Tatausaha
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastTataUsaha)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastTataUsaha.IdTU)
        }
        currentId := lastId + 1
        tatausaha.IdTU = strconv.Itoa(currentId)


        err = c.Insert(tatausaha)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert tatausaha: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+tatausaha.Nama)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeTataUsaha(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdUser := pat.Param(r, "iduser")

        c := session.DB("ccs").C("tatausaha")

        var tatausaha Tatausaha
        err := c.Find(bson.M{"iduser": IdUser}).One(&tatausaha)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find tatausaha: ", err)
            return
        }

        respBody, err := json.MarshalIndent(tatausaha, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}


func AssignPJKelas (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        }		

        IdJadwalKuliah := pat.Param(r, "idjadwalkuliah")

        session := s.Copy()
		defer session.Close()

		r.ParseMultipartForm(500000)
        nim := r.FormValue("nim")

        c := session.DB("ccs").C("mahasiswa")
        d := session.DB("ccs").C("jadwalkuliah")

        var mahasiswa Mahasiswa
        err := c.Find(bson.M{"nim": nim}).One(&mahasiswa)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find mahasiswa: ", err)
            return
        }

        c.Update(bson.M{"nim" : nim}, bson.M{"$set": bson.M{"status": "pjkelas"}})
        d.Update(bson.M{"idjadwalkuliah": IdJadwalKuliah}, bson.M{"$set": bson.M{"idpj": mahasiswa.IdPJ}})

        w.WriteHeader(http.StatusNoContent)
	}
}


