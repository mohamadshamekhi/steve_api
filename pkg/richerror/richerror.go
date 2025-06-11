package richerror

import (
	"errors"
	"fmt"
)

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
	KindBadRequest
)

type RichError struct {
	operation    string
	wrappedError error
	messageKey   string
	templateData map[string]any
	kind         Kind
	meta         map[string]any
}

func New(op string) *RichError {
	return &RichError{operation: op}
}

func (r *RichError) WithWarpError(err error) *RichError {
	r.wrappedError = err
	return r
}

func (r *RichError) WithMessageKey(key string) *RichError {
	r.messageKey = key
	return r
}

func (r *RichError) WithKind(kind Kind) *RichError {
	r.kind = kind
	return r
}

func (r *RichError) WithMeta(meta map[string]any) *RichError {
	r.meta = meta
	return r
}

func (r *RichError) WithTemplateData(data map[string]any) *RichError {
	r.templateData = data
	return r
}

func (r *RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	var re *RichError
	if errors.As(r.wrappedError, &re) {
		return re.Kind()
	}

	return 0
}

func (r *RichError) Error() string {
	var wrappedErrMsg string
	if r.wrappedError != nil {
		wrappedErrMsg = r.wrappedError.Error()
	}
	return fmt.Sprintf("op: %s, kind: %d, key: %s, err: [%s]",
		r.operation, r.Kind(), r.messageKey, wrappedErrMsg)
}

func (r *RichError) Unwrap() error {
	return r.wrappedError
}

func (r *RichError) GetMessageKey() string {
	return r.messageKey
}

func (r *RichError) GetTemplateData() map[string]any {
	return r.templateData
}

func (r *RichError) GetMeta() map[string]any {
	return r.meta
}
func (r *RichError) GetOp() string {
	return r.operation
}

func (r *RichError) MessageKey() string {
	if r.messageKey != "" {
		return r.messageKey
	}
	var re *RichError
	if errors.As(r.wrappedError, &re) {
		return re.MessageKey()
	}
	return ""
}

func (r *RichError) TemplateData() map[string]any {
	if r.templateData != nil {
		return r.templateData
	}
	var re *RichError
	if errors.As(r.wrappedError, &re) {
		return re.TemplateData()
	}
	return nil
}
