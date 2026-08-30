package servertest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tanzor-Disco/skaM/internal/apperrors"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/Tanzor-Disco/skaM/server/register"
)

func handleUser(t *testing.T, user models.RegisterRequest) *httptest.ResponseRecorder {
	t.Helper()
	r := createEncodeRequest(t, user, nil)
	w := httptest.NewRecorder()

	registerHandler := register.NewRegisterHandler(srv.DB, srv.BaseURL, srv.SMTPData)
	registerHandler.HandleRegister(w, r)
	return w
}

func TestHandleRequest_Valid(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras",
		Password: "2313112313",
	}

	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusOK,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleRequest_Empty(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "",
		Username: "",
		Password: "",
	}

	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrInvalidEmail,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleRequest_InvalidEmail(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "huylo",
		Username: "testname",
		Password: "123312123",
	}
	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrInvalidEmail,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleRequest_Duplicate(t *testing.T) {
	userPrev := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras",
		Password: "2313112313",
	}
	user := models.RegisterRequest{
		Email:    "pidoras_@mail.ru",
		Username: "pidoras_new",
		Password: "1321313213131",
	}

	_ = handleUser(t, userPrev)
	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusOK,
		ErrKind:     apperrors.KindErrNone,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleRequest_InvalidUsernameLength(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "example_@mail.ru",
		Username: "125151512515353255151551142144241241",
		Password: "12321331123",
	}
	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrInvalidUsernameLength,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleRequest_InvalidPasswordChars(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "example@mail.ru",
		Username: "huylo",
		Password: "фывфывфйцуйцуол",
	}

	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrForbiddenPasswordChars,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}

func TestHandleRequest_InvalidPasswordLength(t *testing.T) {
	user := models.RegisterRequest{
		Email:    "example@mail.ru",
		Username: "huylo",
		Password: "12313131331313113131313132131313131313131313131313131313131331313131313131313131313131313",
	}
	w := handleUser(t, user)
	defer tdb.DeletePendingUsersByEmail(t, user.Email)

	wanted := wantedResult{
		Code:        http.StatusBadRequest,
		ErrKind:     apperrors.KindErrInvalidPasswordLength,
		CookieExist: false,
	}
	verifyResponse(t, w, wanted)
}
