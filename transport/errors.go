package transport

import (
	"errors"
	"fmt"
)

type errWithReason error

var (
	errEstablishSecureConn errWithReason = errors.New("error establishing secure conn")
)

func WrapErrWithReason(reason errWithReason, explanation error) error {
	return fmt.Errorf(reason.Error()+": %w", explanation)
}
