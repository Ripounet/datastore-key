package aquiet

import (
	"io/ioutil"
	"log"
)

func init() {
	log.Println("a")
	log.SetOutput(ioutil.Discard)
	log.Println("b")
}
