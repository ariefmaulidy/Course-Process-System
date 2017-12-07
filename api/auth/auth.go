package auth

import (
  "fmt"
  "net/http"
  "context"
  "time"
  "log"
  "encoding/json"
  "gopkg.in/mgo.v2"

  "github.com/dgrijalva/jwt-go"
  "../jsonhandler"
)

type Key int

const MyKey Key = 0


type Claims struct {
	IdUser 			int 	`json:"iduser"`
	Username        string  `json:"username"`
  	Class			string	`json:"class"`
	jwt.StandardClaims
}

type Cookie struct
{
	Name		string 		`json:"name"`
	Value 		string 		`json:"value"`
	Expires 	time.Time 	`json:"expires"`
	HttpOnly	bool		`json:"httponly"`
	Class 		string 		`json:"class"`

}

func SetToken(w http.ResponseWriter, r *http.Request,iduser int ,username string,class string) {
  expireToken := time.Now().Add(time.Hour * 5).Unix()
  expireCookie := time.Now().Add(time.Hour * 5)

  claims := Claims{
  	iduser,
    username,
    class,
    jwt.StandardClaims{
        ExpiresAt: expireToken,
        Issuer:    "localhost:8080",
      },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signedToken, _ := token.SignedString([]byte("ccs"))
  cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
  http.SetCookie(w, &cookie)

  respBody, err := json.MarshalIndent(Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true,Class: class}, "", "  ")
      if err != nil {
        log.Fatal(err)
  }

  jsonhandler.ResponseWithJSON(w, respBody, http.StatusCreated)
  return 
}


func Validate(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			jsonhandler.SendWithJSON(res, "belum login", http.StatusOK)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("ccs"), nil
		})
		if err != nil {
			jsonhandler.SendWithJSON(res, "belum login", http.StatusOK)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), MyKey, *claims)
			page(res, req.WithContext(ctx))
		} else {
			jsonhandler.SendWithJSON(res, "belum login", http.StatusOK)
			return
		}
	})
}

func CheckExpiredToken (s *mgo.Session) func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request){
    	session := s.Copy()
        defer session.Close()

		cookie, err := req.Cookie("Auth")
		if err != nil {
			deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
    		http.SetCookie(res, &deleteCookie)
			jsonhandler.SendWithJSON(res, "belum login", http.StatusOK)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("ccs"), nil
		})
		if err != nil {
			deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
    		http.SetCookie(res, &deleteCookie)
			jsonhandler.SendWithJSON(res, "belum login", http.StatusOK)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			context.WithValue(req.Context(), MyKey, *claims)
			jsonhandler.SendWithJSON(res,  claims.Username, http.StatusOK)
		} else {
			deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
    		http.SetCookie(res, &deleteCookie)
			jsonhandler.SendWithJSON(res, "belum login", http.StatusOK)
			return
		}
	}
}

func CheckToken(unchecktoken string,res http.ResponseWriter, req *http.Request) (map[string]interface{}){
  var isok bool
  token, err := jwt.ParseWithClaims(unchecktoken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method")
    }
    return []byte("ccs"), nil
  })
  if err != nil {
    isok = false
  }
  if claims, ok := token.Claims.(*Claims); ok && token.Valid {
    message := map[string]interface{} {
			"username"	  : 	claims.Username,
			"class"	      : 	claims.Class, //customized new name of files comma separated
		}
		return message
  } else {
    isok = false
  }
	if isok{
    errMsg := map[string]interface{} {
	  "message"     :   "Autentikasi",
	  }
	return errMsg
  }
  errMsg := map[string]interface{} {
  "message"     :   "Tidak Autentikasi",
  }
  return errMsg

}

