package datastorekey

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"strconv"
)

func init() {
	http.HandleFunc("/encode", ajaxEncode)
}

func ajaxEncode(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var err error

	kind := trimmedFormValue(r, "kind")
	stringID := trimmedFormValue(r, "stringid")
	intIDstr := trimmedFormValue(r, "intid")
	appID := trimmedFormValue(r, "appid")
	namespace := trimmedFormValue(r, "namespace")

	// Parent (optional)
	kind2 := trimmedFormValue(r, "kind2")
	stringID2 := trimmedFormValue(r, "stringid2")
	intIDstr2 := trimmedFormValue(r, "intid2")

	// Grand-parent (optional)
	kind3 := trimmedFormValue(r, "kind3")
	stringID3 := trimmedFormValue(r, "stringid3")
	intIDstr3 := trimmedFormValue(r, "intid3")

	var key, parent, grandparent *datastore.Key

	if kind2 != "" {
		if kind3 != "" {
			grandparent, err = CreateKey(c, appID, namespace, kind3, stringID3, intID64(intIDstr3), nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		parent, err = CreateKey(c, appID, namespace, kind2, stringID2, intID64(intIDstr2), grandparent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	key, err = CreateKey(c, appID, namespace, kind, stringID, intID64(intIDstr), parent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Fprint(w, keyString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{
		"keystring": key.Encode(),
	})
}

// See https://developers.google.com/appengine/docs/go/datastore/entities#Go_Kinds_and_identifiers
func CreateKey(c appengine.Context, appID string, namespace string, kind string, stringID string, intID int64, parent *datastore.Key) (*datastore.Key, error) {
	// c is the true context of the current request
	// forged is a wrapper context with our custom appID
	forged := &ForgedContext{c, appID}
	// cc is a wrapper context with our custom namespace
	cc, err := appengine.Namespace(forged, namespace)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(
		cc,       // appengine.Context.
		kind,     // Kind.
		stringID, // String ID; empty means no string ID.
		intID,    // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		parent,   // Parent Key; nil means no parent.
	)
	return key, nil
}

func intID64(intIDstr string) int64 {
	if intIDstr == "" {
		return 0
	}
	intID64, _ := strconv.ParseInt(intIDstr, 10, 64)
	return intID64
}
