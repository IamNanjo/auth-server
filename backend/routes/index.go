package routes

import (
	"context"
	"github.com/IamNanjo/authserver/db"
	"github.com/IamNanjo/authserver/pages"
	"net/http"
	"net/url"
	"strings"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	appId := query.Get("app")
	if appId == "" {
		err := pages.Index().Render(context.Background(), w)
		if err != nil {
			Error(w, http.StatusInternalServerError, "Could not render the page")
			return
		}
		return
	}

	redirectTo := query.Get("redirect")
	if redirectTo == "" {
		Error(w, http.StatusBadRequest, "Invalid authentication URL. No redirect page specified")
		return
	}

	redirectURL, err := url.Parse(redirectTo)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid redirect URL")
		return
	}

	app, err := db.GetAppById(appId)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid authentication URL. App not found")
		return
	}

	redirectURLIsAllowed := false

	for _, domain := range app.Domains {
		hostname := redirectURL.Hostname()
		if domain.Name == hostname || (domain.Name[0] == '.' && strings.HasSuffix(domain.Name, redirectURL.Hostname())) {
			redirectURLIsAllowed = true
		}
	}

	if !redirectURLIsAllowed {
		Error(w, http.StatusBadRequest, "Invalid authentication URL. Redirect page is not on the app domains")
		return
	}

	err = pages.Auth(app).Render(context.Background(), w)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Could not render the page")
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path == "/" {
			getIndex(w, r)
		} else {
			pages.Error("Page not found").Render(r.Context(), w)
		}
	default:
		Error(w, http.StatusMethodNotAllowed, "Invalid method "+r.Method+" for route "+r.URL.Path)
	}
}
