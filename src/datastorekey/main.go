package datastorekey

import (
//	"appengine"
	"net/http"
	"html/template"
)

var templates *template.Template

func init() {
	templates = template.New("datastore-keys")
	templates, _ = templates.ParseGlob("template/*.html")
	http.HandleFunc("/decodeAndJump", decodeAndJump)
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	data := extractGetParameters(r)
	render(w, data)
}

func render(w http.ResponseWriter, data Parameters) {
	templates.ExecuteTemplate(w, "index", data)
}

func extractGetParameters(r *http.Request) Parameters{
	data := Parameters{
		"kind": r.FormValue("kind"),
		"stringid": r.FormValue("stringid"),
		"intid": r.FormValue("intid"),
		"appid": r.FormValue("appid"),
		"namespace": r.FormValue("namespace"),
		"keystring": r.FormValue("keystring"),
		"kind2": r.FormValue("kind2"),
		"stringid2": r.FormValue("stringid2"),
		"intid2": r.FormValue("intid2"),
		"kind3": r.FormValue("kind3"),
		"stringid3": r.FormValue("stringid3"),
		"intid3": r.FormValue("intid3"),
	}
	return data
}

func decodeAndJump(w http.ResponseWriter, r *http.Request) {
	keystring := r.FormValue("keystring")
	key, err := datastore.DecodeKey(keyString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := "https://appengine.google.com/datastore/explorer?submitted=1&app_id=" + key.AppID() + "&show_options=yes&viewby=gql&query=SELECT+*+FROM+"+key.Kind()+"+WHERE+__key__%3DKEY%28%27"+keystring+"%27%29&options=Run+Query"
	http.Redirect(w, r, url, 301)
}

type Parameters map[string]interface{}
