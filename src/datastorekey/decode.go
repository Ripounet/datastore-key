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

	response := recursiveJsonString(key) 
	fmt.Fprint(w, response)
}

func jsonifyKey(key *datastore.Key) (s string) {
	b, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func recursiveJsonString(key *datastore.Key) string{
	return recursiveJson(key).String() 
}

func recursiveJson(key *datastore.Key) Response{
	var parentJson Response
	if key.Parent() != nil {
		parentJson = recursiveJson(key.Parent())
	}
	return Response{
		"stringID": key.StringID(),
		"intID": key.IntID(),
		"kind": key.Kind(),
		"appID": key.AppID(),
		"namespace": key.Namespace(),
		"parent": parentJson,
	}
}

// See http://nesv.blogspot.fr/2012/09/super-easy-json-http-responses-in-go.html
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



// Copied-pasted from appengine/datastore/key.go :
//func DecodeKey(encoded string) (*Key, error) {
//	// Re-add padding.
//	if m := len(encoded) % 4; m != 0 {
//		encoded += strings.Repeat("=", 4-m)
//	}
//
//	b, err := base64.URLEncoding.DecodeString(encoded)
//	if err != nil {
//		return nil, err
//	}
//
//	ref := new(pb.Reference)
//	if err := proto.Unmarshal(b, ref); err != nil {
//		return nil, err
//	}
//
//	return protoToKey(ref)
//}
