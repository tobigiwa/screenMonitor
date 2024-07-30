package backend

import (
	"fmt"
	"log"
	"net/http"
)

var (
	ErrLinkExpired  = fmt.Errorf("link expired")
	ErrDuplicateKey = fmt.Errorf("user with email account already exit")
)

func (a App) clientError(w http.ResponseWriter, errStatus int, err error) {
	http.Error(w, err.Error(), errStatus)
	log.Println("clientError", err)
}

func (a App) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	log.Println("serverError", err)
}
