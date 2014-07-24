package datastorekey

import (
	"appengine_internal"
)

// fakeContext implements appengine.Context, but does nothing.
// It's an empty shell which contains only a custom appId.
type fakeContext struct {
	appId string
}

func (c *fakeContext) FullyQualifiedAppID() string {
	return c.appId
}

func (c *fakeContext) Debugf(format string, args ...interface{}) {
	// NOOP
}

func (c *fakeContext) Infof(format string, args ...interface{}) {
	// NOOP
}

func (c *fakeContext) Warningf(format string, args ...interface{}) {
	// NOOP
}

func (c *fakeContext) Errorf(format string, args ...interface{}) {
	// NOOP
}

func (c *fakeContext) Criticalf(format string, args ...interface{}) {
	// NOOP
}

func (c *fakeContext) Call(service, method string, in, out appengine_internal.ProtoMessage, opts *appengine_internal.CallOptions) error {
	// NOOP
	return nil
}

func (c *fakeContext) Request() interface{} {
	// NOOP
	return nil
}
