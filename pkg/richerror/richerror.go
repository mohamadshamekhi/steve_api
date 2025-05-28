package richerror

import "errors"

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type RichError struct {
	operation string
	warpError error
	message   string
	kind      Kind
	meta      map[string]any
}

func New(operation string) RichError {
	return RichError{operation: operation}
}

func (r RichError) WithWarpError(err error) RichError {
	r.warpError = err

	return r
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message

	return r
}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind

	return r
}

func (r RichError) WithMeta(meta map[string]any) RichError {
	r.meta = meta

	return r
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	var re RichError

	ok := errors.As(r.warpError, &re)
	if !ok {
		return 0
	}

	return re.Kind()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	if r.warpError == nil {
		return ""
	}

	var re RichError

	ok := errors.As(r.warpError, &re)
	if !ok {
		return r.warpError.Error()
	}

	return re.Message()
}

func (r RichError) Error() string {
	if r.message != "" {
		return r.message
	}

	if r.warpError == nil {
		return ""
	}

	var re RichError
	if errors.As(r.warpError, &re) {
		return re.Error()
	}

	return r.warpError.Error()
}
