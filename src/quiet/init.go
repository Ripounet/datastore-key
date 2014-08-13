// This is a desperate attempt to discard the log message
// "appengine: not running under devappserver2; using some default configuration".
// It does not really work, as the order of initialization
// between this package and appengine/datastore is unspecified.
package quiet

import (
	"io/ioutil"
	"log"
)

func init() {
	log.SetOutput(ioutil.Discard)
}
