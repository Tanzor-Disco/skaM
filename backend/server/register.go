package server
import (
	"net/http"
	"io"
	"encoding/json"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/Tanzor-Disco/skaM/db"

)

type userData struct {
	Email string
	Username string
	Password string
}

func (s *server) handleRegister( w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	buf,err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var currUser userData 
	err = json.Unmarshal(buf,&currUser)
	if err != nil {
		sendJSON(w,false,"an error occured while parsing the json",err,400)
		return 
	}
	
	passwordHash,err := bcrypt.GenerateFromPassword([]byte(currUser.Password),bcrypt.DefaultCost)
	if err != nil {
		sendJSON(w,false,"an error occured while hashing the password",err,500)
		return 
	}

	userDB := db.User {
		Name:currUser.Username,
		Email:currUser.Email,
		PasswordHash:string(passwordHash),
		LastYearActive:time.Now().Year(),
	}
	
	err = s.db.CreateUser(userDB)
	if err != nil {
		sendJSON(w,false,"The email is already taken",err,500)
		return 
	}

	sendJSON(w,true,"successfully added a user",nil,200)

}

