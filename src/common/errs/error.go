package errs

import "errors"

type Error struct {
	error
	code string
}

func New(text string, code ...string) *Error {
	er := &Error{
		error: errors.New(text),
		code:  "",
	}
	if len(code) > 0{
		er.code = code[0]
	}
	return er
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Wrap(err error) *Error {
	return &Error{error: errors.New(e.Error()+" "+ err.Error()) , code: e.code}
}