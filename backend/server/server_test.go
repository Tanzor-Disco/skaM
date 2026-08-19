package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/Tanzor-Disco/skaM/auth"
	"github.com/Tanzor-Disco/skaM/db"
	"github.com/Tanzor-Disco/skaM/db/migrations"
	"github.com/Tanzor-Disco/skaM/internal/testutils"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/lpernett/godotenv"
)

func loadURI() string {
	//load the environment
	err := godotenv.Load("../.env_test")
	if err != nil {
		log.Fatal(err)
	}
	testURI := os.Getenv("TEST_URI")
	return testURI
}

func setUpDB(testURI string) (*testutils.TestDB, *db.DB) {

	//reset the db and update to the latest migration
	err := migrations.Reset(testURI)
	if err != nil {
		log.Fatalf("migrations:%v", err)
	}

	//connect to db
	tdb, err := testutils.Connect(testURI)
	if err != nil {
		log.Fatalf("setUpDB: %v", err)
	}
	db, err := db.Connect(testURI)
	if err != nil {
		log.Fatalf("setUpDB: %v", err)
	}

	return tdb, db
}

func runSMTPServer() *exec.Cmd {
	cmd := exec.Command("mailpit")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func createServerData(testURI string) *server {
	TestBaseURL := "http://localhost:8080"
	TestSMTPData := models.NewSMTPData("", "", "localhost", "localhost:1025", "test@example.com", TestBaseURL)
	TestServerData := models.NewServerData(testURI, TestSMTPData)
	srv, err := newServer(TestServerData)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

var srv *server
var tdb *testutils.TestDB
var database *db.DB

func TestMain(m *testing.M) {
	testURI := loadURI()
	tdb, database = setUpDB(testURI)
	cmd := runSMTPServer()
	srv = createServerData(testURI)

	//launch and finish the tests
	code := m.Run()
	tdb.Close()
	srv.db.Close()
	cmd.Process.Kill()
	os.Exit(code)
}

func createDecodeRequest[T any](t *testing.T, body T, sessionString *string) *http.Request {
	t.Helper()
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Didn't manage to encode body into bytes: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/", bytes.NewReader(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	if sessionString != nil {
		cookie := auth.CreateSessionCookie(*sessionString)
		r.AddCookie(&cookie)
	}
	return r
}

func createRawRequest(t *testing.T, body string, sessionString *string) *http.Request {
	t.Helper()
	r := httptest.NewRequest("POST", "/api/", bytes.NewReader([]byte(body)))
	r.Header.Set("Content-Type", "application/json")
	if sessionString != nil {
		cookie := auth.CreateSessionCookie(*sessionString)
		r.AddCookie(&cookie)
	}
	return r
}

type wantedResult struct {
	Code        int
	Success     bool
	ErrKind     string
	CookieExist bool
}

func verifyResponse(t *testing.T, w *httptest.ResponseRecorder, wanted wantedResult) {
	t.Helper()
	if w.Code != wanted.Code {
		t.Fatalf("verifyResponse: code of the response doesn't match:\n got %v\n wanted %v", w.Code, wanted.Code)
	}
	var gotBody serverResponseBody
	err := json.NewDecoder(w.Body).Decode(&gotBody)
	if err != nil {
		t.Fatalf("verifyResponse: couldn't decode the server response")
	}
	if gotBody.Success != wanted.Success {
		t.Fatalf("verifyResponse: success of the response doesn't match:\n got %v\n wanted %v", gotBody.Success, wanted.Success)
	}
	if gotBody.ErrorKind != wanted.ErrKind {
		t.Fatalf("verifyResponse: error kind of the response doesn't match:\n got %v\n wanted %v", gotBody.ErrorKind, wanted.ErrKind)
	}
	cookies := w.Result().Cookies()
	if (len(cookies) == 1) != wanted.CookieExist {
		t.Fatalf("got %v cookies, wanted cookies: %v", len(cookies), wanted.CookieExist)
	}
	if len(cookies) == 0 {
		return
	}
	cookie := cookies[0]
	if cookie.Name != "session_string" {
		t.Fatalf("verifyResponse: cookie name mismatch:\n got %v\n wanted %v", cookie.Name, "session_string")
	}
}
