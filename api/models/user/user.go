package user

import (


	"goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "../../auth"
    "../../jsonhandler"
)

type User struct {
	IdUser		int		`json:"iduser"`
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

func RoutesUser(mux *goji.Mux, session *mgo.Session) {
	

    mux.HandleFunc(pat.Get("/user"), AllUser(session)) //untuk retrieve smua yang di db
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
            jsonhandler.SendJSON(w, "Database error", http.StatusInternalServerError)
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
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&user)
        if err != nil {
            jsonhandler.SendJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("ccs").C("user")

        //untuk auto increment
        var lastUser User
        var lastId  int

        err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastUser)
        if err != nil {
            lastId = 0
        } else {
            lastId,err = strconv.Atoi(lastUser.IdUser)
        }
        currentId := lastId + 1
        user.IdUser = strconv.Itoa(currentId)

        passEncrypt := sha256.Sum256([]byte(user.Password))
        user.Password = fmt.Sprintf("%x", passEncrypt)

        err = c.Insert(user)
        if err != nil {
            if mgo.IsDup(err) {
                jsonhandler.SendJSON(w, "duplicate", http.StatusOK)
                return
            }

            jsonhandler.SendJSON(w, "Database error", http.StatusNotFound)
            log.Println("Failed insert user: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+user.Username)
        w.WriteHeader(http.StatusCreated)
    }
}

func login(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var user User
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&user)
        if err != nil {
            jsonhandler.SendJSON(w, "Incorrect body", http.StatusBadRequest)
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
    jsonhandler.SendJSON(w, "logout", http.StatusOK);
    return
  }
}
