package tempcgd

import (
    "encoding/json"
    "log"
    "net/http"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../jsonhandler"
    "../../auth"
)

type TempCGD struct {
    IdCGD	        int			`json:"idcgd"`
    JumlahPesan     int         `json:"jumlahpesan"`
    IdUser          int         `json:"iduser"`
}

func RoutesTempCGD(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/tempcgd"), AllTempCGD(session)) //untuk retrieve smua yang di db
    //mux.HandleFunc(pat.Post("/addtempcgd"), AddTempCGD(session))
    mux.HandleFunc(pat.Get("/tempcgd/:idtempcgd"), auth.Validate(GetAttributeTempCGD(session)))
}

func AllTempCGD(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("tempcgd")

        var tempcgd []TempCGD
        err := c.Find(bson.M{}).All(&tempcgd)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all tempcgd: ", err)
            return
        }

        respBody, err := json.MarshalIndent(tempcgd, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

/*func AddTempCGD(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var tempcgd TempCGD
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&tempcgd)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("tempcgd")

        //untuk auto increment
        var lastTempCGD TempCGD
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastTempCGD)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastTempCGD.IdTempCGD
        }
        currentId := lastId + 1
        tempcgd.IdTempCGD = currentId

        err = c.Insert(tempcgd)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert tempcgd: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}*/

func GetAttributeTempCGD(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        }

        session := s.Copy()
        defer session.Close()

        IdTempCGD := pat.Param(r, "idtempcgd")

        c := session.DB("ccs").C("tempcgd")

        var tempcgd TempCGD
        err := c.Find(bson.M{"$and": []bson.M{bson.M{"idtempcgd": IdTempCGD},bson.M{"iduser": claims.IdUser}}}).One(&tempcgd)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find tempcgd: ", err)
            return
        }

        respBody, err := json.MarshalIndent(tempcgd, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

