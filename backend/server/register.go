package server
import (
	"net/http"
	"io"
	"encoding/json"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"errors"

)

type userData struct {
	Email string
	Username string
	Password string
}

func readHTTP(req *http.Request) ([]byte,error){
	defer req.Body.Close()
	buf,err := io.ReadAll(req.Body)
	if err != nil {
		return make([]byte,0),err
	}
	return buf,nil

}

func getUserData(buf []byte) (userData,error) {
	var currUser userData 
	err := json.Unmarshal(buf,&currUser)
	return currUser,err
}

func (s *server) registerUser(userDB db.User) error {
	err := s.db.CreateUser(userDB)
	return err
}

func (s *server) handleRegister( w http.ResponseWriter, req *http.Request) {
	buf,err := readHTTP(req)
	if err != nil {
		body := newServerResponseBody(false,err,"ERR_REQUEST_BODY_READ")
		sendJSON(w,http.StatusBadRequest,body)
		return 
	}


	currUser,err := getUserData(buf)
	if err != nil {
		body := newServerResponseBody(false,err,"ERR_INVALID_JSON")
		sendJSON(w,http.StatusBadRequest,body)
		return 
	}
	
	passwordHash,err := bcrypt.GenerateFromPassword([]byte(currUser.Password),bcrypt.DefaultCost)
	if err != nil {
		body := newServerResponseBody(false,err,"ERR_HASHING_PASSWORD")
		sendJSON(w,http.StatusInternalServerError,body)
		return 
	}
	
	date := time.Now().Year()
	userDB := db.User {
		Email:currUser.Email,
		Username:currUser.Username,
		PasswordHash: string(passwordHash),
		LastYearActive:date,

	}
	err = s.registerUser(userDB)
	if err != nil {
		var body serverResponseBody
		if errors.Is(err,apperrors.ErrEmailTaken) {
			body = newServerResponseBody(false,err,"ERR_EMAIL_TAKEN")
			sendJSON(w,http.StatusBadRequest,body)
		} else {
			body = newServerResponseBody(false,err,"ERR_DB")
			sendJSON(w,http.StatusInternalServerError,body)
		} 

		return 
	}
	
	body := newServerResponseBody(true,nil,"ERR_NONE")
	sendJSON(w,http.StatusOK,body)

}

