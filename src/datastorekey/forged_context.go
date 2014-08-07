package datastorekey

import (
	"appengine"
)

// ForgedContext embeds an appengine.Context, and has a custom appId.
// It overrides method FullyQualifiedAppID()
type ForgedContext struct {
	appengine.Context
	appId string
}

func (c *ForgedContext) FullyQualifiedAppID() string {
	return c.appId
}
