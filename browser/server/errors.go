package webserver

import (
	"fmt"
	"net/http"
)

var (
	ErrLinkExpired  = fmt.Errorf("link expired")
	ErrDuplicateKey = fmt.Errorf("user with email account already exit")
)

func (a App) clientError(w http.ResponseWriter, errStatus int, err error) {
	w.Header().Del("lastSaturday")
	http.Error(w, err.Error(), errStatus)
	fmt.Println(err)
}

func (a App) serverError(w http.ResponseWriter, err error) {
	w.Header().Del("lastSaturday")
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	fmt.Println(err)
}
