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
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	data := Parameters{
		"kind": r.FormValue("kind"),
		"stringid": r.FormValue("stringid"),
		"intid": r.FormValue("intid"),
		"appid": r.FormValue("appid"),
		"namespace": r.FormValue("namespace"),
		"keystring": r.FormValue("keystring"),
	}
	templates.ExecuteTemplate(w, "index", data)
}

type Parameters map[string]interface{}
