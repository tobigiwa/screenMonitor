package webserver

import (
	"fmt"
	"net/http"
	"pkg/types"
	"strings"
)

func (a *App) SetCategory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	val := r.Form.Get("category")
	appName := r.Form.Get("appName")
	category := types.Category(val)

	msg := types.Message{
		Endpoint: strings.TrimPrefix(r.URL.Path, "/"),
		SetCategoryRequest: types.SetCategoryRequest{
			AppName:  appName,
			Category: category,
		},
	}
	res, err := a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}
	if !res.SetCategoryResponse.IsCategorySet {
		a.serverError(w, fmt.Errorf("error setting app category"))
		return
	}

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}
