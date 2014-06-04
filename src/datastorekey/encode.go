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
	kind := r.FormValue("kind")
	stringID := r.FormValue("stringid")
	intIDstr := r.FormValue("intid")
	appID := r.FormValue("appid")
	namespace := r.FormValue("namespace")
	keyString, err := encodeKey(c, appID, namespace, kind, stringID, intID64(intIDstr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, keyString)
}

// See https://developers.google.com/appengine/docs/go/datastore/entities#Go_Kinds_and_identifiers
func encodeKey(c appengine.Context, appID string, namespace string, kind string, stringID string, intID int64) (string, error) {
	// Function Context.Namespace does not seem to work??
	//c.Debugf("Current context is %v \n", c)
	cc, err := appengine.Namespace(c, namespace)
	if err != nil {
		return "", err
	}
	//c.Debugf("Tempo context is %v \n", cc)
	key := datastore.NewKey(
		cc,        // appengine.Context.
		kind,     // Kind.
		stringID, // String ID; empty means no string ID.
		intID,        // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		nil)      // Parent Key; nil means no parent.

	return key.Encode(), nil
}

func intID64(intIDstr string) int64{
	if intIDstr=="" {
		return 0
	}
	intID64, _ := strconv.ParseInt(intIDstr, 10, 64)
	return intID64
}


// Copied-pasted from appengine/datastore/key.go, and modified :
//func NewKey(appID string, namespace string, kind, stringID string, intID int64, parent *datastore.Key) *datastore.Key {
//	if parent != nil {
//		namespace = parent.namespace
//	} else {
//		s := &basepb.StringProto{}
//		c.Call("__go__", "GetNamespace", &basepb.VoidProto{}, s, nil)
//		namespace = s.GetValue() // "" if the RPC fails
//	}
//
//	return &datastore.Key{
//		kind:      kind,
//		stringID:  stringID,
//		intID:     intID,
//		parent:    parent,
//		appID:     c.FullyQualifiedAppID(),
//		namespace: namespace,
//	}
//}