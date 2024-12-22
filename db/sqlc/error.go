package sqlc

import (
	u "backend-masterclass/util"
	"database/sql"
	"errors"
)

const (
	uniquViolation                                    = "23505"
	foreignKeyViolation                               = "23503"
	connection_exception                              = "08000"
	connection_does_not_exist                         = "08003"
	connection_failure                                = "08006"
	sqlclient_unable_to_establish_sqlconnection       = "08001"
	sqlserver_rejected_establishment_of_sqlconnection = "08004"
	transaction_resolution_unknown                    = "08007"
	protocol_violation                                = "08P01"
)

var (
	constraintViolations = u.StringSlice{uniquViolation, foreignKeyViolation}
	connectionErrors     = u.StringSlice{
		connection_exception, connection_does_not_exist,
		connection_failure, sqlclient_unable_to_establish_sqlconnection,
		sqlserver_rejected_establishment_of_sqlconnection,
		transaction_resolution_unknown, protocol_violation,
	}

	// When the package is imported, these variables get an address.
	ErrForbiddenInput = errors.New("forbidden input")
	ErrConnection     = errors.New("connection error")
	ErrRecordNotFound = sql.ErrNoRows
)
