package pengajar

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
    "../matakuliah"
    "../ruangan"
    "../jadwalkuliah"
    "../dosen"

)

type Pengajar struct {
	IdPengajar			int		`json:"idpengajar"`
	IdMataKuliah	    int		`json:"idmatakuliah"`
	IdUser	            int	 	`json:"iduser"`
}

func RoutesPengajar(mux *goji.Mux, session *mgo.Session) {
    mux.HandleFunc(pat.Get("/pengajar"), AllPengajar(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpengajar"), AddPengajar(session))
	mux.HandleFunc(pat.Get("/pengajar/:idpengajar"), auth.Validate(GetAttributePengajar(session)))
	mux.HandleFunc(pat.Get("/jadwalmengajar"), auth.Validate(JadwalMengajar(session)))
    mux.HandleFunc(pat.Get("/getpengajar/:idmatakuliah"), auth.Validate(GetPengajar(session)))
}

func AllPengajar(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pengajar")

        var pengajar []Pengajar
        err := c.Find(bson.M{}).All(&pengajar)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pengajar: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pengajar, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPengajar(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var pengajar Pengajar
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pengajar)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("pengajar")

        //untuk auto increment
        var lastpengajar Pengajar
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastpengajar)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastpengajar.IdPengajar
        }
        currentId := lastId + 1
        pengajar.IdPengajar = currentId

        err = c.Insert(pengajar)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pengajar: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributePengajar(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdPengajar := pat.Param(r, "idpengajar")

        c := session.DB("ccs").C("pengajar")

        var pengajar Pengajar
        err := c.Find(bson.M{"idpengajar": IdPengajar}).One(&pengajar)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pengajar: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pengajar, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
	}

}

func JadwalMengajar(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
		}
		
		if claims.Class == "Dosen"{
			session := s.Copy()
			defer session.Close()
	
			IdUser := claims.IdUser

			c := session.DB("css").C("pengajar")

			var pengajar []Pengajar

			type jadwalfix struct{
				MataKuliah 	 string      `json:"matakuliah"`
				Waktu	     string      `json:"waktu"`
				Ruangan		 string      `json:"ruangan"`
            }
            realjadwal := []jadwalfix{}
            err := c.Find(bson.M{"iduser": IdUser}).One(&pengajar)
            
            if err != nil{
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed find pengajar: ", err)
                return
            }

            for _,element := range pengajar{
                a := session.DB("css").C("matakuliah")
                var matkul matakuliah.MataKuliah
                err := a.Find(bson.M{"idmatakuliah":element.IdMataKuliah}).One(&matkul) 

                d := session.DB("css").C("jadwalkuliah")
                var jadwal []jadwalkuliah.JadwalKuliah
                err = d.Find(bson.M{"idmatakuliah":element.IdMataKuliah}).All(&jadwal)
                
                if err != nil{
                    jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                    log.Println("Failed find pengajar: ", err)
                    return
                }

                for _,element2 := range jadwal{
                    e := session.DB("css").C("ruangan")
                    var ruangan ruangan.Ruangan
                    e.Find(bson.M{"idruangan":element2.IdRuangan}).One(&ruangan)
                    j := jadwalfix{MataKuliah:matkul.NamaMataKuliah, Waktu:element2.Waktu, Ruangan:ruangan.NamaRuangan }
                    realjadwal = append(realjadwal,j)
                }
            }

            respBody, err := json.MarshalIndent(realjadwal, "", "  ")
            if err != nil {
                log.Fatal(err)
            }
            jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
                
        }else{
            jsonhandler.SendWithJSON(w, "you dont have permission", http.StatusNotFound)
            return
        }
	}
}

func GetPengajar(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pengajar")
        d := session.DB("ccs").C("dosen")

        var DataPengajar []Pengajar
        var DataDosen []dosen.Dosen 

        IdMataKuliah, _ := strconv.Atoi(pat.Param(r, "idmatakuliah"))


        err := c.Find(bson.M{"idmatakuliah":IdMataKuliah}).All(&DataPengajar)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pengajar: ", err)
            return
        }

        var tempdosen dosen.Dosen

        for _,data := range DataPengajar {
            d.Find(bson.M{"iduser": data.IdUser}).One(&tempdosen)
            DataDosen = append(DataDosen, tempdosen)
        }

        respBody, err := json.MarshalIndent(DataDosen, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}