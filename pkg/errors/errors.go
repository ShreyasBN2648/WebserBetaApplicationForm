package errorsCust

import "errors"

var (
	ErrorInterfaceNotPointer        = errors.New("interface is not a pointer")
	ErrInsertingBetaUserApplication = errors.New(("error inserting beta user application to collection"))
)
