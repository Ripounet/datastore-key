package datastorekey

import (
	//	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/decode", ajaxDecode)
}

func decode(w http.ResponseWriter, r *http.Request) {
}

func ajaxDecode(w http.ResponseWriter, r *http.Request) {
	keyString := r.FormValue("keystring")
	key, err := datastore.DecodeKey(keyString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := recursiveJsonResponse(key) 
	fmt.Fprint(w, response)
}

func jsonifyKey(key *datastore.Key) (s string) {
	b, err := json.Marshal(key)
	if err != nil {
		return ""
	}
	return string(b)
}

func recursiveJsonResponse(key *datastore.Key) string{
	parentJson := ""
	if key.Parent() != nil {
		parentJson = recursiveJsonResponse(key.Parent())
	}
	return Response{
		"stringID": key.StringID(),
		"intID": key.IntID(),
		"kind": key.Kind(),
		"appID": key.AppID(),
//		"namespace": key.AppNamespace(),
		"parent": parentJson,
	}.String()
}

// See http://nesv.blogspot.fr/2012/09/super-easy-json-http-responses-in-go.html
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

