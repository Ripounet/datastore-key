package datastorekey

import (
	"appengine/datastore"
	"html/template"
	"net/http"
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
	_ = autodecode(data["keystring"].(string), data)
	// If autodecode failed, render the page with key not decoded
	render(w, data)
}

func render(w http.ResponseWriter, data Parameters) {
	templates.ExecuteTemplate(w, "index", data)
}

func extractGetParameters(r *http.Request) Parameters {
	data := Parameters{
		"kind":      r.FormValue("kind"),
		"stringid":  r.FormValue("stringid"),
		"intid":     r.FormValue("intid"),
		"appid":     r.FormValue("appid"),
		"namespace": r.FormValue("namespace"),
		"keystring": r.FormValue("keystring"),
		"kind2":     r.FormValue("kind2"),
		"stringid2": r.FormValue("stringid2"),
		"intid2":    r.FormValue("intid2"),
		"kind3":     r.FormValue("kind3"),
		"stringid3": r.FormValue("stringid3"),
		"intid3":    r.FormValue("intid3"),
	}
	return data
}

// IF keystring was given as GET parameter
// THEN it would be nice that all decoded values are directly served in the html
func autodecode(keystring string, data Parameters) error {
	if keystring == "" {
		// Nothing to decode
		return nil
	}
	//b := true
	if data["appid"] != "" || data["kind"] != "" || data["intid"] != "" || data["stringid"] != "" {
		// Don't overwrite user-provided values
		return nil
	}
	key, err := datastore.DecodeKey(keystring)
	if err != nil {
		return err
	}
	fillFields(key, data)
	return nil
}

func fillFields(key *datastore.Key, data map[string]interface{}) {
	data["kind"] = key.Kind()
	data["stringid"] = key.StringID()
	data["intid"] = key.IntID()
	data["appid"] = key.AppID()
	data["namespace"] = key.Namespace()
	if key.Parent() != nil {
		data["kind2"] = key.Parent().Kind()
		data["stringid2"] = key.Parent().StringID()
		data["intid2"] = key.Parent().IntID()
		if key.Parent().Parent() != nil {
			data["kind3"] = key.Parent().Parent().Kind()
			data["stringid3"] = key.Parent().Parent().StringID()
			data["intid3"] = key.Parent().Parent().IntID()
		}
	}
}

func decodeAndJump(w http.ResponseWriter, r *http.Request) {
	keystring := r.FormValue("keystring")
	key, err := datastore.DecodeKey(keystring)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := "https://appengine.google.com/datastore/explorer?submitted=1&app_id=" + key.AppID() + "&show_options=yes&viewby=gql&query=SELECT+*+FROM+" + key.Kind() + "+WHERE+__key__%3DKEY%28%27" + keystring + "%27%29&options=Run+Query"
	data := Response{
	    "keystring": keystring,
		"url": url,
	}
	fillFields(key, data)
	templates.ExecuteTemplate(w, "jump", data)

	// This very direct redirect caused a long blank screen because Google Console is so slow.
	//http.Redirect(w, r, url, 301)
}

type Parameters map[string]interface{}
