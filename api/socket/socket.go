package socket

import(
	"net/http"
	"time"
    "log"

	"github.com/gorilla/websocket"
	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"


    "../models/tatausaha"
    "../models/mahasiswa"
    "../models/dosen"
    "../models/pesan"
    "../jsonhandler"
    "../auth"
)

var upgrader = websocket.Upgrader{}
var clients = make(map[*websocket.Conn]bool)

type DataSend struct {
	IsiPesan 		string 	`json:"isipesan"`
	NamaPengirim 	string 	`json:"namapengirim"`
	ClassPengirim 	string 	`json:"classpengirim"`
	Status 			string 	`json:"status"` //menandakan siapa orangnya
}

func RoutesSocket(mux *goji.Mux, session *mgo.Session){
	mux.HandleFunc(pat.New("/roomcgd/:idcgd"), auth.Validate(GetRoomCGD(session)))
}	 	

func remove(s []DataSend, r DataSend) []DataSend {
    for i, v := range s {
        if v == r {
            return append(s[:i], s[i+1:]...)
        } 
    }
    return s
}

func GetRoomCGD(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
	 	claims, ok := r.Context().Value(auth.MyKey).(auth.Claims)
        if !ok {
          http.NotFound(w, r)
          return
        }
    	//for connection
	 	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	 	if error != nil {
        	http.NotFound(w, r)
        	return
    	}
	 	clients[conn] = true
	 	//

	 	idcgd := pat.Param(r, "idcgd")
	 	session := s.Copy()
    	
		cPesan := session.DB("ccs").C("pesan")
		cMahasiswa := session.DB("ccs").C("mahasiswa")
		cDosen := session.DB("ccs").C("dosen")
		cTu	:= session.DB("ccs").C("tatausaha")

		go func(conn *websocket.Conn){
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					for client := range clients {
						if(client == conn){
							delete(clients, client)
						}
					}
					if(clients[conn] == false){
						conn.Close()
						defer session.Close()
						log.Println("Starting Listen server....")
						return
					}
				}
			}
		}(conn)

	 	go func(conn *websocket.Conn){
	 		ch := time.Tick(250 * time.Millisecond)	 		
	 		var pesan []pesan.Pesan
	 		var datasend []DataSend

		 	for range ch {
		 		if(clients[conn] == true){
			 		//log.Println("Lapak " +idlapak+ " Work")
			 		log.Println(clients)

			 		//pesan
				    err := cPesan.Find(bson.M{"idcgd": idcgd}).All(&pesan)
				    if err != nil {
				        jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
				        log.Println("Failed find pesan: ", err)
				        return
				    }

				    for _,message := range pesan {
				    	if message.ClassPengirim == "TataUsaha"{
				    		var pengirim tatausaha.Tatausaha
				    		err := cTu.Find(bson.M{"iduser": message.IdPengirim}).One(&pengirim)
					        if err != nil {
					            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
					            log.Println("Failed find tatausaha: ", err)
					            return
					        }
					        if message.IdPengirim == claims.IdUser {
					        	datasend = append(datasend, DataSend{IsiPesan: message.IsiPesan, NamaPengirim: pengirim.Nama, ClassPengirim: message.ClassPengirim, Status: "Saya"})
				    		} else {
				    			datasend = append(datasend, DataSend{IsiPesan: message.IsiPesan, NamaPengirim: pengirim.Nama, ClassPengirim: message.ClassPengirim, Status: "Bukan Saya"})
				    		}
				    	}
				    	if message.ClassPengirim == "Dosen"{
				    		var pengirim dosen.Dosen
				    		err := cDosen.Find(bson.M{"iduser": message.IdPengirim}).One(&pengirim)
					        if err != nil {
					            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
					            log.Println("Failed find dosen: ", err)
					            return
					        }
					        if message.IdPengirim == claims.IdUser {
					        	datasend = append(datasend, DataSend{IsiPesan: message.IsiPesan, NamaPengirim: pengirim.Nama, ClassPengirim: message.ClassPengirim, Status: "Saya"})
				    		} else {
				    			datasend = append(datasend, DataSend{IsiPesan: message.IsiPesan, NamaPengirim: pengirim.Nama, ClassPengirim: message.ClassPengirim, Status: "Bukan Saya"})
				    		}
				    	}
				    	if message.ClassPengirim == "Mahasiswa"{
				    		var pengirim mahasiswa.Mahasiswa
				    		err := cMahasiswa.Find(bson.M{"iduser": message.IdPengirim}).One(&pengirim)
					        if err != nil {
					            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
					            log.Println("Failed find mahasiswa: ", err)
					            return
					        }
					        if message.IdPengirim == claims.IdUser {
					        	datasend = append(datasend, DataSend{IsiPesan: message.IsiPesan, NamaPengirim: pengirim.Nama, ClassPengirim: message.ClassPengirim, Status: "Saya"})
				    		} else {
				    			datasend = append(datasend, DataSend{IsiPesan: message.IsiPesan, NamaPengirim: pengirim.Nama, ClassPengirim: message.ClassPengirim, Status: "Bukan Saya"})
				    		}
				    	}
					}

					conn.WriteJSON(datasend)
					
					for i := 0; i < len(datasend); i++ {
						datasend = remove(datasend, datasend[i])
						i--;
					}
			 	}
			}	
	 	}(conn)
	 }
}

