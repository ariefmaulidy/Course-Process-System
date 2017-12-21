package tatausaha

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "strconv"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../mahasiswa"
    "../jadwalkuliah"
    "../pesananruangan"
    "../matakuliah"
    "../ruangan"
    "../../auth"
    "../../jsonhandler"
)

type Tatausaha struct {
	Nama		string 	`json:"nama"`
	IdUser      int     `json:"iduser"`
	Departemen	string	`json:"departemen"`
	Fakultas	string	`json:"fakultas"`
}

type DataSend struct{
    Ruangan             []ruangan.Ruangan               `json:"ruangan"`
    Jadwal              []jadwalkuliah.JadwalKuliah     `json:"jadwalkuliah"`
    Pesanan             []pesananruangan.PesananRuangan `json:"pesananruangan"`
    DataJadwalMatkul    []matakuliah.MataKuliah         `json:"datajadwalmatkul"`
    DataPesananMatkul   []matakuliah.MataKuliah         `json:"datapesananmatkul"`
}

func RoutesTataUsaha(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/tatausaha"), AllTataUsaha(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addtatausaha"), AddTataUsaha(session))
    mux.HandleFunc(pat.Get("/tatausaha/:iduser"), GetAttributeTataUsaha(session))
    //untuk memasukkan BAP cek di fungsi AddBAP di model bap
    mux.HandleFunc(pat.Put("/assignpjkelas/:idjadwalkuliah"), auth.Validate(AssignPJKelas(session)))
    mux.HandleFunc(pat.Put("/editpjkelas/:idjadwalkuliah"), auth.Validate(EditPJKelas(session)))
    //untuk book ruangan cek di fungsi menambahPesananRuangan  di model pesanan ruangan
    mux.HandleFunc(pat.Get("/listbookruangan/:tanggal"), auth.Validate(GetListRuanganBook(session)))
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
            lastId = lastTataUsaha.IdUser
        }
        currentId := lastId + 1
        tatausaha.IdUser = currentId


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


func AssignPJKelas(s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        }		

        if claims.Class == "TataUsaha"{
            IdJadwalKuliah, _ := strconv.Atoi(pat.Param(r, "idjadwalkuliah"))

            session := s.Copy()
            defer session.Close()

            r.ParseMultipartForm(500000)
            IdUser, _ := strconv.Atoi(r.FormValue("iduser"))

            c := session.DB("ccs").C("mahasiswa")
            d := session.DB("ccs").C("jadwalkuliah")
            e := session.DB("ccs").C("pesertakuliah")

            var mahasiswa mahasiswa.Mahasiswa
            err := c.Find(bson.M{"iduser": IdUser}).One(&mahasiswa)
            if err != nil {
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed find mahasiswa: ", err)
                return
            }

            e.Update(bson.M{"iduser" : IdUser}, bson.M{"$set": bson.M{"statuspj": "Penanggung Jawab"}})
            d.Update(bson.M{"idjadwalkuliah": IdJadwalKuliah}, bson.M{"$set": bson.M{"idpj": mahasiswa.IdUser}})

            w.WriteHeader(http.StatusNoContent)
        }
	}
}

func EditPJKelas(s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
    return func(w http.ResponseWriter, r *http.Request){
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        }       

        if claims.Class == "TataUsaha"{
            IdJadwalKuliah, _ := strconv.Atoi(pat.Param(r, "idjadwalkuliah"))

            session := s.Copy()
            defer session.Close()

            r.ParseMultipartForm(500000)
            IdUser, _ := strconv.Atoi(r.FormValue("iduser"))

            c := session.DB("ccs").C("mahasiswa")
            d := session.DB("ccs").C("jadwalkuliah")
            e := session.DB("ccs").C("pesertakuliah")

            var mahasiswa mahasiswa.Mahasiswa
            err := c.Find(bson.M{"iduser": IdUser}).One(&mahasiswa)
            if err != nil {
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed find mahasiswa: ", err)
                return
            }

            e.Update(bson.M{"iduser" : IdUser}, bson.M{"$set": bson.M{"statuspj": ""}})
            d.Update(bson.M{"idjadwalkuliah": IdJadwalKuliah}, bson.M{"$set": bson.M{"idpj": 0}})

            w.WriteHeader(http.StatusNoContent)
        }
    }
}

func GetListRuanganBook(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        }
        if claims.Class == "TataUsaha"{
            session := s.Copy()
            defer session.Close()
            
            Tanggal := pat.Param(r, "tanggal")
            
            var hari string

            var datasend DataSend

            var matkul matakuliah.MataKuliah
    
            layout := "2006-01-02T15:04:05.000Z"
            Tgl,_ := time.Parse(layout,Tanggal)
            hari = Tgl.Weekday().String()
            switch hari{
                case "Monday":
                    hari = "Senin"
                case "Tuesday":
                    hari = "Selasa"
                case "Wednesday":
                    hari = "Rabu"
                case "Thursday":
                    hari = "Kamis"
                case "Friday":
                    hari = "Jumat"
                case "Saturday":
                    hari = "Sabtu"
                case "Sunday":
                    hari = "Minggu"
            }
            
            log.Println(hari,Tanggal)
            c := session.DB("ccs").C("jadwalkuliah")
            d := session.DB("css").C("pesananruangan")
            e := session.DB("ccs").C("matakuliah")
            f := session.DB("ccs").C("ruangan")
    
            err := c.Find(bson.M{"hari": hari}).Sort("idruangan","waktu").All(&datasend.Jadwal)
            if err != nil {
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed find jadwal: ", err)
                return
            }
    
            err = d.Find(bson.M{"tanggal": Tanggal}).Sort("idruangan","waktu").All(&datasend.Pesanan)
            if err != nil {
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed find pesanan: ", err)
                return
            }

            err = f.Find(bson.M{}).All(&datasend.Ruangan)
            if err != nil {
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed find ruangan: ", err)
                return
            }            

            for _,data := range datasend.Jadwal {
                e.Find(bson.M{"idmatakuliah": data.IdMataKuliah}).One(&matkul)
                datasend.DataJadwalMatkul = append(datasend.DataJadwalMatkul, matkul)
            }   

            for _,data := range datasend.Pesanan {
                e.Find(bson.M{"idmatakuliah": data.IdMataKuliah}).One(&matkul)
                datasend.DataPesananMatkul = append(datasend.DataPesananMatkul, matkul)
            }

            respBody, err := json.MarshalIndent(datasend, "", "  ")
            if err != nil {
              log.Fatal(err)
            }
            jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
            return
        }else{
            jsonhandler.SendWithJSON(w, "you dont have permission", http.StatusNotFound)
            return
        }
        
    }
}




