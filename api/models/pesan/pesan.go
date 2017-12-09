package pesan

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
    "../tempcgd"
    "../chatgroupdiscussion"
)

type Pesan struct {
	IdPesan	        int			`json:"idpesan"`
	IsiPesan	    string		`json:"isipesan"`
    IdCGD           int         `json:"idcgd"`
    IdPengirim      int         `json:"idpengirim"`
    ClassPengirim   string      `json:"classpengirim"`
}

func RoutesPesan(mux *goji.Mux, session *mgo.Session) {
    mux.HandleFunc(pat.Get("/pesan"), AllPesan(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addpesan/:idcgd"), auth.Validate(AddPesan(session)))
    mux.HandleFunc(pat.Get("/pesan/:idpesan"), GetAttributePesan(session))
}

func AllPesan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("pesan")

        var pesan []Pesan
        err := c.Find(bson.M{}).All(&pesan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all pesan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddPesan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        }   

        session := s.Copy()
        defer session.Close()

        var pesan Pesan
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&pesan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        var IdCGD int

        IdCGD,err = strconv.Atoi(pat.Param(r, "idcgd"))

        c := session.DB("ccs").C("pesan")
        d := session.DB("ccs").C("tempcgd")
        e := session.DB("ccs").C("chatgroupdiscussion")

        var tempcgd tempcgd.TempCGD

        err = d.Find(bson.M{"$and": []bson.M{bson.M{"idcgd": IdCGD}, bson.M{"iduser": claims.IdUser}}}).One(&tempcgd)
        if err != nil {
            tempcgd.IdCGD = IdCGD
            tempcgd.JumlahPesan = 1
            tempcgd.IdUser = claims.IdUser
            d.Insert(tempcgd)
        } else {
            tempcgd.JumlahPesan = tempcgd.JumlahPesan + 1
            d.Update(bson.M{"$and": []bson.M{bson.M{"idcgd": IdCGD}, bson.M{"iduser": claims.IdUser}}}, bson.M{"$set": bson.M{"jumlahpesan": tempcgd.JumlahPesan}})
        }

        var cgd chatgroupdiscussion.ChatGroupDiscussion

        cgd.JumlahPesan = cgd.JumlahPesan + 1
        e.Update(bson.M{"idcgd": IdCGD}, bson.M{"$set": bson.M{"jumlahpesan": cgd.JumlahPesan}})
        
        //untuk auto increment
        var lastPesan Pesan
        var lastId int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastPesan)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastPesan.IdPesan
        }
        currentId := lastId + 1
        pesan.IdPesan = currentId

        pesan.IdCGD = IdCGD
        pesan.IdPengirim = claims.IdUser
        pesan.ClassPengirim = claims.Class

        err = c.Insert(pesan)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert pesan: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributePesan(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdPesan := pat.Param(r, "idpesan")

        c := session.DB("ccs").C("pesan")

        var pesan Pesan
        err := c.Find(bson.M{"idpesan": IdPesan}).One(&pesan)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find pesan: ", err)
            return
        }

        respBody, err := json.MarshalIndent(pesan, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

