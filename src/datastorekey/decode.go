package datastorekey

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/decode", ajaxDecode)
}

func ajaxDecode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	keyString := trimmedFormValue(r, "keystring")
	c := appengine.NewContext(r)
	c.Infof("Decoding %v\n", keyString)

	key, err := datastore.DecodeKey(keyString)
	if err != nil {
		c.Errorf("Failed: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := recursiveJsonString(key)
	c.Infof("%v\n", response)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, response)
}

func jsonifyKey(key *datastore.Key) (s string) {
	b, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func recursiveJsonString(key *datastore.Key) string {
	return recursiveJson(key).String()
}

func recursiveJson(key *datastore.Key) Response {
	var parentJson Response
	if key.Parent() != nil {
		parentJson = recursiveJson(key.Parent())
	}
	return Response{
		"stringID":  key.StringID(),
		"intID":     key.IntID(),
		"kind":      key.Kind(),
		"appID":     key.AppID(),
		"namespace": key.Namespace(),
		"parent":    parentJson,
	}
}

//  See http://nesv.blogspot.fr/2012/09/super-easy-json-http-responses-in-go.html
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
