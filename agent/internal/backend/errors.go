package backend

import (
	"fmt"
	"net/http"
)

var (
	ErrLinkExpired  = fmt.Errorf("link expired")
	ErrDuplicateKey = fmt.Errorf("user with email account already exit")
)

func (a App) clientError(w http.ResponseWriter, errStatus int, err error) {
	http.Error(w, err.Error(), errStatus)
	fmt.Println("clientError", err)
}

func (a App) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	fmt.Println("serverError", err)
}
