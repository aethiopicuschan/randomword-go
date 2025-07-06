package randomword

import "errors"

var ErrNumberLessThanOne = errors.New("number must be greater than 0")
var ErrLengthLessThanOne = errors.New("length must be greater than 0")
var ErrDoFuncCannotBeNil = errors.New("do function cannot be nil")
var ErrUnexpectedResponse = errors.New("unexpected response from server")
var ErrInternal = errors.New("internal error")
