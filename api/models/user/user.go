package user

import (
    "encoding/json"
    "log"
    "time"
    "net/http"
    "fmt"
    "crypto/sha256"

	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
    "../tatausaha"
    "../dosen"
    "../pengelolaruangan"
    "../mahasiswa"
)

type User struct {
	IdUser		int		`json:"iduser"`
	Username	string	`json:"username"`
	Password	string	`json:"password"`
    Class       string  `json:"class"`
}

func RoutesUser(mux *goji.Mux, session *mgo.Session) {
    mux.HandleFunc(pat.Get("/user"), AllUsers(session)) //untuk retrieve smua yang di db
    mux.HandleFunc(pat.Post("/register"), Register(session))
    mux.HandleFunc(pat.Post("/login"), Login(session)) // login
    mux.HandleFunc(pat.Get("/checkexpiredtoken"), auth.CheckExpiredToken(session))
    mux.HandleFunc(pat.Post("/logout"), auth.Validate(logout(session))) // Logout
}

func EnsureUser(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("ccs").C("user")

    index := mgo.Index{
        Key:        []string{"username"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }

    err := c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }
}

func AllUsers(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("ccs").C("user")

        var users []User
        err := c.Find(bson.M{}).All(&users)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all users: ", err)
            return
        }

        respBody, err := json.MarshalIndent(users, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func Register(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var user User

        r.ParseMultipartForm(500000)

        if r.FormValue("username") != ""{
            user.Username = r.FormValue("username")
        }
        if r.FormValue("class") != ""{
            user.Class = r.FormValue("class")
        }
        if r.FormValue("password") != ""{
            user.Password = r.FormValue("password")
        }

        c := session.DB("ccs").C("user")

        //untuk auto increment
        var lastUser User
        var lastId  int

        err := c.Find(nil).Sort("-$natural").Limit(1).One(&lastUser)
        if err != nil {
            lastId = 0
        } else {
            lastId = lastUser.IdUser
        }
        currentId := lastId + 1
        user.IdUser = currentId

        passEncrypt := sha256.Sum256([]byte(user.Password))
        user.Password = fmt.Sprintf("%x", passEncrypt)

        err = c.Insert(user)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendWithJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert user: ", err)
            return
        }

        if user.Class == "TataUsaha"{
            var tu tatausaha.Tatausaha

            if r.FormValue("nama") != "" {
                tu.Nama = r.FormValue("nama")
            }
            if r.FormValue("departemen") != "" {
                tu.Departemen = r.FormValue("departemen")
            }
            if r.FormValue("fakultas") != "" {
                tu.Fakultas = r.FormValue("fakultas")
            }
            tu.IdUser = currentId

            d := session.DB("ccs").C("tatausaha")
            
            err = d.Insert(tu)
            if err != nil { 
                jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
                log.Println("Failed insert tata usaha: ", err)
                return
            }
        }

        if user.Class == "Dosen"{
            var dos dosen.Dosen

            if r.FormValue("nama") != "" {
                dos.Nama = r.FormValue("nama")
            }
            if r.FormValue("departemen") != "" {
                dos.Departemen = r.FormValue("departemen")
            }
            if r.FormValue("fakultas") != "" {
                dos.Fakultas = r.FormValue("fakultas")
            }
            if r.FormValue("nidn") != "" {
                dos.NIDN = r.FormValue("nidn")
            }
            dos.IdUser = currentId

            e := session.DB("ccs").C("dosen")
            
            err = e.Insert(dos)
            if err != nil { 
                jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
                log.Println("Failed insert dosen: ", err)
                return
            }
        }

        if user.Class == "PengelolaRuangan"{
            var pengelola pengelolaruangan.PengelolaRuangan    
            
            if r.FormValue("nama") != "" {
                pengelola.Nama = r.FormValue("nama")
            }
            pengelola.IdUser = currentId

            e := session.DB("ccs").C("pengelolaruangan")
                        
            err = e.Insert(pengelola)
            if err != nil { 
                jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
                log.Println("Failed insert pengelola ruangan: ", err)
                return
            }
        }

        if user.Class == "Mahasiswa"{
            var mhs mahasiswa.Mahasiswa

            if r.FormValue("nama") != "" {
                mhs.Nama = r.FormValue("nama")
            }
            if r.FormValue("departemen") != "" {
                mhs.Departemen = r.FormValue("departemen")
            }
            if r.FormValue("nim") != "" {
                mhs.NIM = r.FormValue("nim")
            }
            if r.FormValue("status") != "" {
                mhs.Status = r.FormValue("status")
            }
            mhs.IdUser = currentId

            e := session.DB("ccs").C("mahasiswa")
                        
            err = e.Insert(mhs)
            if err != nil { 
                jsonhandler.SendWithJSON(w, "Database error", http.StatusNotFound)
                log.Println("Failed insert mahasiswa: ", err)
                return
            }

        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+user.Username)
        w.WriteHeader(http.StatusCreated)
    }
}

func Login(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var user User
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&user)
        if err != nil {
            jsonhandler.SendWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("user")

        passEncrypt := sha256.Sum256([]byte(user.Password))
        user.Password = fmt.Sprintf("%x", passEncrypt)

        err = c.Find(bson.M{"username": user.Username,"password": user.Password}).One(&user)
        if err != nil {
        respBody , _ := json.MarshalIndent(jsonhandler.MessageJSON{Status: false, Message: "user not found"}, "", "  ")
            jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
            log.Println("Failed find user: ", err)
            return
        } else{
            auth.SetToken(w,r,user.IdUser,user.Username,user.Class)
            return
        }
    }
}

func logout(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    session := s.Copy()
    defer session.Close()


    deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
    http.SetCookie(w, &deleteCookie)
    jsonhandler.SendWithJSON(w, "logout", http.StatusOK);
    return
  }
}
