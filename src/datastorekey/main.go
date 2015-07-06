// Google Datastore Key viewer with decoder/encoder.
//
// [WTFPL](http://en.wikipedia.org/wiki/WTFPL)
package datastorekey

import (
	"html/template"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"
)

var templates *template.Template

func init() {
	templates = template.New("datastore-keys")
	templates, _ = templates.ParseGlob("template/*.html")
	http.HandleFunc("/decodeAndJump", decodeAndJump)
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	data := extractGetParameters(r)

	keyString := data["keystring"].(string)
	if keyString != "" {
		c.Infof("Decoding %v\n", keyString)
		err := autodecode(keyString, data)
		if err == nil {
			c.Infof("Decoded %v\n", data)
		} else {
			c.Errorf("Failed: %v\n", err.Error())
			// If autodecode failed, render the page with key not decoded
		}
	}
	render(w, data)
}

func render(w http.ResponseWriter, data Parameters) {
	templates.ExecuteTemplate(w, "index", data)
}

func extractGetParameters(r *http.Request) Parameters {
	data := Parameters{
		"kind":      trimmedFormValue(r, "kind"),
		"stringid":  trimmedFormValue(r, "stringid"),
		"intid":     trimmedFormValue(r, "intid"),
		"appid":     trimmedFormValue(r, "appid"),
		"namespace": trimmedFormValue(r, "namespace"),
		"keystring": trimmedFormValue(r, "keystring"),
		"kind2":     trimmedFormValue(r, "kind2"),
		"stringid2": trimmedFormValue(r, "stringid2"),
		"intid2":    trimmedFormValue(r, "intid2"),
		"kind3":     trimmedFormValue(r, "kind3"),
		"stringid3": trimmedFormValue(r, "stringid3"),
		"intid3":    trimmedFormValue(r, "intid3"),
	}
	return data
}

// IF keystring was given as GET parameter
// THEN it is nice that all decoded values are directly served in the html
func autodecode(keystring string, data Parameters) error {
	if keystring == "" {
		// Nothing to decode
		return nil
	}
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
	keystring := trimmedFormValue(r, "keystring")
	key, err := datastore.DecodeKey(keystring)
	if err != nil {
		w.WriteHeader(400)
		templates.ExecuteTemplate(w, "error-bad-keystring", keystring)
		return
	}
	url := "https://appengine.google.com/datastore/explorer?submitted=1&app_id=" + key.AppID() +
		"&show_options=yes&viewby=gql&query=SELECT+*+FROM+" + key.Kind() +
		"+WHERE+__key__%3DKEY%28%27" + keystring + "%27%29" +
		"&namespace=" + key.Namespace() +
		"&options=Run+Query"
	data := Response{
		"keystring": keystring,
		"url":       url,
	}
	fillFields(key, data)
	templates.ExecuteTemplate(w, "jump", data)

	// This very direct redirect caused a long blank screen because Google Console is so slow.
	//http.Redirect(w, r, url, 301)
}

type Parameters map[string]interface{}

func trimmedFormValue(r *http.Request, paramName string) string {
	return strings.TrimSpace(r.FormValue(paramName))
}
