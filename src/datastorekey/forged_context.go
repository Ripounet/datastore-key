package datastorekey

import (
	"appengine"
)

// forgedContext embeds an appengine.Context, and has a custom appId.
// It overrides method FullyQualifiedAppID()
type forgedContext struct {
	appengine.Context
	appId string
}

func (c *forgedContext) FullyQualifiedAppID() string {
	return c.appId
}
