package lucifer

import "errors"

var ErrUnsupportedOperation = errors.New("lucifer: operation not supported for this driver")

// ErrUnsupportedDriver is returned if an invalid driver type is requested.
var ErrUnsupportedDriver = errors.New("lucifer: driver not supported")

// ErrBridgeNotFound is returned if a bridge cannot be found and some operations fails because of it.
var ErrBridgeNotFound = errors.New("lucifer: bridge not found")
