package httpmsg

import (
	"errors"
	"fmt"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"net/http"
)

// TODO: this temperary to ignore linter error.(maggic number).
const (
	internalStatus = 500
)

func Error(err error) (message string, code int) {
	var re richerror.RichError
	if !errors.As(err, &re) {
		return err.Error(), http.StatusBadRequest
	}
	msg := re.Message()
	code = mapKindToHTTPStatusCode(re.Kind())
	// we should not expose unexpected error messages
	if code >= internalStatus {
		// TODO - we have to use log instead of print
		fmt.Println("internal error: ", msg)
		msg = errmsg.ErrorMsgSomethingWentWrong
	}

	return msg, code
}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:

		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:

		return http.StatusNotFound
	case richerror.KindForbidden:

		return http.StatusForbidden
	case richerror.KindUnexpected:

		return http.StatusInternalServerError
	default:

		return http.StatusBadRequest
	}
}
