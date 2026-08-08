package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/internal/testutils"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/lpernett/godotenv"
)

var srv *server
var tdb *testutils.TestDB

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env_test")
	if err != nil {
		log.Fatal(err)
	}
	URI := os.Getenv("URI")
	srv, err = newServer(URI)
	if err != nil {
		log.Fatal(err)
	}
	tdb, err = testutils.Connect(URI)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	tdb.Close()
	srv.db.Close()
	os.Exit(code)
}

func handleUser(t *testing.T, user models.RegisterRequest) *httptest.ResponseRecorder {
	t.Helper()
	bodyBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Didn't manage to encode body into bytes: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.handleRegister(w, req)

	return w
}

func checkStatus(t *testing.T, w *httptest.ResponseRecorder, status int) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("The code of the response didn't match: got %v, wanted %v", w.Code, status)
	}
}

func validateResponseBody(t *testing.T, resp *http.Response, correct serverResponseBody) {
	var body serverResponseBody
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("Failed to decode json: %v", err)
	}
	if body.Success != correct.Success || body.ErrorKind != correct.ErrorKind {
		t.Fatalf("Body of the response is incorrect:\n Got %+v\n Wanted %+v", body, correct)
	}

}

func TestHandleRequest_Valid(t *testing.T) {
	user := models.RegisterRequest {
		Email:    "pidoras_@mail.ru",
		Username: "pidoras",
		Password: "2313112313",
	}

	correct := newServerResponseBody(true, apperrors.KindErrNone)
	serverRecorder := handleUser(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	serverResp := serverRecorder.Result()
	checkStatus(t, serverRecorder, http.StatusOK)
	validateResponseBody(t, serverResp, correct)
}

func TestHandleRequest_Duplicate(t *testing.T) {
	user := models.RegisterRequest {
		Email:    "pidoras_@mail.ru",
		Username: "pidoras",
		Password: "2313112313",
	}
	correct := newServerResponseBody(false, apperrors.KindErrEmailTaken)
	//creating the first user
	handleUser(t, user)
	//creating a duplicate
	serverRecorder := handleUser(t, user)
	defer tdb.DeleteUsersByEmail(t, user.Email)
	serverResp := serverRecorder.Result()
	checkStatus(t, serverRecorder, http.StatusBadRequest)
	validateResponseBody(t, serverResp, correct)
}
