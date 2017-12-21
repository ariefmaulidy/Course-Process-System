 package bap

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
    "../jadwalkuliah"
    "../pesertakuliah"
    "../dosen"
    "../matakuliah"
    "../ruangan"
)

type BAP struct {
	IdBAP	        int			`json:"idbap"`
	Tanggal	        string      `json:"tanggal"`
    IdJadwalKuliah  int         `json:"idjadwalkuliah"`
    TopikKuliah     string      `json:"topikkuliah"`
    JumlahPeserta   int         `json:"jumlahpeserta"`
    IdUser          int         `json:"iduser"`
    Waktu           string      `json:"waktu"`
    Catatan         string      `json:"catatan"`
    CreatedBy       int         `json:"createdby"`
    UpdatedBy       int         `json:"updatedby"`
}

type DataSend struct{
    DataBAP         []BAP                           `json:"databap"`
    DataPengajar    []dosen.Dosen                   `json:"datapengajar"`
    DataJadwal      jadwalkuliah.JadwalKuliah       `json:"datajadwal"`
    DataMatkul      matakuliah.MataKuliah           `json:"datamatkul"`
    DataPeserta     []pesertakuliah.PesertaKuliah   `json:"datapeserta"`
    NamaRuangan     string                           `json:"ruangan"`
    PJDosen         string                           `json:"datadosen"`
}

func RoutesBAP(mux *goji.Mux, session *mgo.Session) {
    mux.HandleFunc(pat.Get("/bap"), AllBAP(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Get("/jadwalbap/:idjadwal"), GetJadwalBAP(session))
    mux.HandleFunc(pat.Post("/addbap"), auth.Validate(AddBAP(session)))
    mux.HandleFunc(pat.Get("/bap/:idbap"), GetAttributeBAP(session))
    mux.HandleFunc(pat.Put("/editbap/:idbap"), auth.Validate(EditBAP(session)))
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

func GetJadwalBAP(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdJadwalKuliah, _ := strconv.Atoi(pat.Param(r, "idjadwal"))

        c := session.DB("ccs").C("bap")
        d := session.DB("ccs").C("jadwalkuliah")
        e := session.DB("ccs").C("matakuliah")
        f := session.DB("ccs").C("pesertakuliah")
        g := session.DB("ccs").C("ruangan")
        h := session.DB("ccs").C("dosen")
        

        var datasend DataSend
        var kelas ruangan.Ruangan
        var dosenpj dosen.Dosen
        var pengajar dosen.Dosen

        err := c.Find(bson.M{"idjadwalkuliah":IdJadwalKuliah}).All(&datasend.DataBAP)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get bap: ", err)
            return
        }

        for _,data := range datasend.DataBAP{
            h.Find(bson.M{"iduser":data.IdUser}).One(&pengajar)
            datasend.DataPengajar = append(datasend.DataPengajar,pengajar)
        }
        err = d.Find(bson.M{"idjadwalkuliah":IdJadwalKuliah}).One(&datasend.DataJadwal)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get jadwal: ", err)
            return
        }

        err = e.Find(bson.M{"idmatakuliah":datasend.DataJadwal.IdMataKuliah}).One(&datasend.DataMatkul)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get matakuliah: ", err)
            return
        }

        err = f.Find(bson.M{"idjadwalkuliah":IdJadwalKuliah}).All(&datasend.DataPeserta)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all peserta: ", err)
            return
        }

        err = g.Find(bson.M{"idruangan":datasend.DataJadwal.IdRuangan}).One(&kelas)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all ruangan: ", err)
            return
        }
        datasend.NamaRuangan = kelas.NamaRuangan

        err = h.Find(bson.M{"iduser":datasend.DataMatkul.IdKordinator}).One(&dosenpj)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all bap: ", err)
            return
        }
        datasend.PJDosen = dosenpj.Nama

        respBody, err := json.MarshalIndent(datasend, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddBAP(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
		}
        
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
            lastId = lastBAP.IdBAP
        }
        currentId := lastId + 1
        bap.IdBAP = currentId
        bap.CreatedBy = claims.IdUser

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

func EditBAP(s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
    return func(w http.ResponseWriter, r *http.Request){
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
        }
        
        session := s.Copy()
        defer session.Close()

        var varmap map[string]interface{}
        in := []byte(`{}`)
        json.Unmarshal(in, &varmap)

        IdBAP := pat.Param(r, "idbap")

        c := session.DB("ccs").C("bap")

        r.ParseMultipartForm(500000)

        if r.FormValue("tanggal") != ""{
            varmap["tanggal"] = r.FormValue("tanggal")
        }

        if r.FormValue("topikkuliah") != ""{
            varmap["topikkuliah"] = r.FormValue("topikkuliah")
        }

        if r.FormValue("jumlahpeserta") != ""{
            varmap["jumlahpeserta"] = r.FormValue("jumlahpeserta")
        }

        if r.FormValue("iduser") != ""{
            varmap["iduser"] = r.FormValue("iduser")
        }

        varmap["updatedBy"] = claims.IdUser

        err := c.Update(bson.M{"idbap": IdBAP},bson.M{"$set": varmap})
        if err != nil {
            switch err {
            default:
                jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed update BAP: ", err)
                jsonhandler.SendWithJSON(w, "Gagal mengupdate BAP", http.StatusOK)
                return
            case mgo.ErrNotFound:
                jsonhandler.SendWithJSON(w, "BAP not found", http.StatusNotFound)
                return
            }
        }
    }
}