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
)

type ChatGroupDiscussion struct {
	IdCGD	        int			`json:"idcgd"`
	IdMataKuliah	int		    `json:"idmatakuliah"`
    JumlahPesan     int         `json:"jumlahpesan"`
}

func RoutesChatGroupDiscussion(mux *goji.Mux, session *mgo.Session) {

    mux.HandleFunc(pat.Get("/cgd"), AllChatGroupDiscussion(session)) //untuk retrieve smua yang di db
    //mux.HandleFunc(pat.Post("/addcgd"), AddChatGroupDiscussion(session))
    mux.HandleFunc(pat.Get("/cgd/:idcgd"), GetAttributeChatGroupDiscussion(session))
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

/*func AddChatGroupDiscussion(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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
}*/

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

