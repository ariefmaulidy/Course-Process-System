package pengelolaruangan

import (
	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
    "../pesananruangan"
)

type PengelolaRuangan struct {
	Nama	            string	`json:"nama"`
	IdUser	            int	 	`json:"iduser"`
}

func RoutesPengelolaRuangan(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/pengelolaruangan"), AllPengelolaRuangan(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpengelolaruangan"), AddPengelolaRuangan(session))
    mux.HandleFunc(pat.Get("/pengelolaruangan/:idpengelolaruangan"), auth.Validate(GetAttributePengelolaRuangan(session)))
    mux.HandleFunc(pat.Put("/persetujuanpesanan/:idpesanan"), auth.Validate(PersetujuanPesanan(session)))
    mux.HandleFunc(pat.Put("/penolakanpesanan/:idpesanan"), auth.Validate(PenolakanPesanan(session)))
}

func AllPengelolaRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pengelolaruangan")

        var pengelolaruangan []PengelolaRuangan
        err := c.Find(bson.M{}).All(&pengelolaruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pengelolaruangan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pengelolaruangan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPengelolaRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var pengelolaruangan PengelolaRuangan
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pengelolaruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("pengelolaruangan")

        //untuk auto increment
        var lastPengelolaRuangan PengelolaRuangan
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastPengelolaRuangan)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastPengelolaRuangan.IdPengelolaRuangan)
        }
        currentId := lastId + 1
        pengelolaruangan.IdPengelolaRuangan = strconv.Itoa(currentId)

        err = c.Insert(pengelolaruangan)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pengelolaruangan: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+pengelolaruangan.IdPengelolaRuangan)
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributePengelolaRuangan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdPengelolaRuangan := pat.Param(r, "idpengelolaruangan")

        c := session.DB("ccs").C("pengelolaruangan")

        var pengelolaruangan PengelolaRuangan
        err := c.Find(bson.M{"idpengelolaruangan": IdPengelolaRuangan}).One(&pengelolaruangan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pengelolaruangan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pengelolaruangan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func PersetujuanPesanan(s *mgo.Session) func(w http.ResponseWriter,r *http.Request){
    return func(w http.ResponseWriter,r *http.Request){
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
        }

        session := s.Copy()
        defer session.Close()

        IdPesanan = pat.Param(r,"idpesanan")

        c := session.DB("ccs").C("pesananruangan")
        c.Update(bson.M{"idpesanan": IdPesanan}, bson.M{"$set": bson.M{"status": "Disetujui","confirmedby":claims.IdUser}})
        w.WriteHeader(http.StatusNoContent)
    }
}

func PenolakanPesanan(s *mgo.Session) func(w http.ResponseWriter,r *http.Request){
    return func(w http.ResponseWriter,r *http.Request){
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
        }
        
        session := s.Copy()
        defer session.Close()

        IdPesanan = pat.Param(r,"idpesanan")

        c := session.DB("ccs").C("pesananruangan")
        c.Update(bson.M{"idpesanan": IdPesanan}, bson.M{"$set": bson.M{"status": "Ditolak"}})
        w.WriteHeader(http.StatusNoContent)
    }
}