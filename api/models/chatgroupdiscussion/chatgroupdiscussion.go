package chatgroupdiscussion

import (
     "encoding/json"
    "log"
    "net/http"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../jsonhandler"
    "../jadwalkuliah"
    "../matakuliah"
    "../../auth"
    "../tempcgd"
)

type ChatGroupDiscussion struct {
	IdCGD	        int			`json:"idcgd"`
	IdJadwalKuliah	int		    `json:"idjadwalkuliah"`
    JumlahPesan     int         `json:"jumlahpesan"`
}

func RoutesChatGroupDiscussion(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/cgd"), AllChatGroupDiscussion(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/addcgd"), AddChatGroupDiscussion(session))
    mux.HandleFunc(pat.Get("/cgd/:idcgd"), GetAttributeChatGroupDiscussion(session))
    mux.HandleFunc(pat.Get("/roomcgd"), auth.Validate(AllRoomChatGroupDiscussion(session)))
}

func AllChatGroupDiscussion(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("chatgroupdiscussion")

        var chatgroupdiscussion []ChatGroupDiscussion
        err := c.Find(bson.M{}).All(&chatgroupdiscussion)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all chatgroupdiscussion: ", err)
            return
        }

        respBody, err := json.MarshalIndent(chatgroupdiscussion, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AddChatGroupDiscussion(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var chatgroupdiscussion ChatGroupDiscussion
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&chatgroupdiscussion)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("chatgroupdiscussion")

        //untuk auto increment
        var lastChatGroupDiscussion ChatGroupDiscussion
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastChatGroupDiscussion)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastChatGroupDiscussion.IdCGD
        }
        currentId := lastId + 1
        chatgroupdiscussion.IdCGD = currentId

        err = c.Insert(chatgroupdiscussion)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert chatgroupdiscussion: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}

func GetAttributeChatGroupDiscussion(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        IdCGD := pat.Param(r, "idcgd")

        c := session.DB("ccs").C("chatgroupdiscussion")

        var chatgroupdiscussion ChatGroupDiscussion
        err := c.Find(bson.M{"idcgd": IdCGD}).One(&chatgroupdiscussion)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed find chatgroupdiscussion: ", err)
            return
        }

        respBody, err := json.MarshalIndent(chatgroupdiscussion, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func AllRoomChatGroupDiscussion(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
            http.NotFound(w, r)
            return
        } 
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("chatgroupdiscussion")
        d := session.DB("ccs").C("matakuliah")
        e := session.DB("ccs").C("jadwalkuliah")
        f := session.DB("ccs").C("tempcgd")

        type DataSend struct {
            DataCGD        []ChatGroupDiscussion                    `json:"datacgd"`
            DataMataKuliah          []matakuliah.MataKuliah            `json:"datamatakuliah"`
            DataJadwalKuliah           []jadwalkuliah.JadwalKuliah            `json:"datajadwalkuliah"`
            Unread      []int               `json:"unread"`
        }

        var dataSend DataSend

        err := c.Find(bson.M{}).All(&dataSend.DataCGD)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all chatgroupdiscussion: ", err)
            return
        }

        for _,data := range dataSend.DataCGD {
            var tempjadwal jadwalkuliah.JadwalKuliah
            e.Find(bson.M{"idjadwalkuliah": data.IdJadwalKuliah}).One(&tempjadwal)
            dataSend.DataJadwalKuliah = append(dataSend.DataJadwalKuliah, tempjadwal)
            var tempunread tempcgd.TempCGD
            f.Find(bson.M{"idcgd": data.IdCGD, "iduser": claims.IdUser}).One(&tempunread)
            totalunread := data.JumlahPesan - tempunread.JumlahPesan
            dataSend.Unread = append(dataSend.Unread, totalunread)
        }

        for _,data := range dataSend.DataJadwalKuliah {
            var tempmatkul matakuliah.MataKuliah
            d.Find(bson.M{"idmatakuliah": data.IdMataKuliah}).One(&tempmatkul)
            dataSend.DataMataKuliah = append(dataSend.DataMataKuliah, tempmatkul)
        }

        respBody, err := json.MarshalIndent(dataSend, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}



