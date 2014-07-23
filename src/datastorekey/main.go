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

func index(w http.ResponseWriter, r *http.Request) {
	data := extractGetParameters(r)
	render(w, data)
}


func decodeAndJump(w http.ResponseWriter, r *http.Request) {
	data := extractGetParameters(r)
	data["jumpToDatastoreViewer"] = true
	render(w, data)
}

func render(w http.ResponseWriter, data Parameters) {
	templates.ExecuteTemplate(w, "index", data)
}

type Parameters map[string]interface{}
