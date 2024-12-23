package sqlc

import (
	u "backend-masterclass/util"
	"database/sql"
	"errors"
)

const (
	uniquViolation                                = "23505"
	foreignKeyViolation                           = "23503"
	connectionException                           = "08000"
	connectionDoesNotExist                        = "08003"
	connectionFailure                             = "08006"
	sqlclientUnableToEstablishSqlConnection       = "08001"
	sqlserverRejectedEstablishmentOfSqlConnection = "08004"
	transactionResolutionUnknown                  = "08007"
	protocolViolation                             = "08P01"
)

var (
	// Group of possible violation error codes.
	constraintViolations = u.StringSlice{uniquViolation, foreignKeyViolation}

	// Group of possible connection error code.
	connectionErrors = u.StringSlice{
		connectionException, connectionDoesNotExist,
		connectionFailure, sqlclientUnableToEstablishSqlConnection,
		sqlserverRejectedEstablishmentOfSqlConnection,
		transactionResolutionUnknown, protocolViolation,
	}

	// When the package is imported, these variables get an address.
	ErrForbiddenInput = errors.New("forbidden input")
	ErrConnection     = errors.New("connection error")
	ErrRecordNotFound = sql.ErrNoRows
)
