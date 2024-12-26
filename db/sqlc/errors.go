package sqlc

import (
	u "backend-masterclass/util"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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

	// Group of possible connection error codes/ messages.
	connectionErrors = u.StringSlice{
		connectionException, connectionDoesNotExist,
		connectionFailure, sqlclientUnableToEstablishSqlConnection,
		sqlserverRejectedEstablishmentOfSqlConnection,
		transactionResolutionUnknown, protocolViolation,
		sql.ErrConnDone.Error(),
	}

	// When the package is imported, these variables get an address.
	ErrForbiddenInput = errors.New("forbidden input")
	ErrConnection     = errors.New("connection error")
	ErrRecordNotFound = sql.ErrNoRows
)

// translateSQLError translates a SQL error into a more readable error.
func translateSQLError(err error) (trError error, errNotNil bool) {

	errNotNil = err != nil
	trError = err
	// pgconn.PgError is initialised to nil so we can't use it inside As.
	// That's why we take its address, which still implements the error interface.
	var pgxErr *pgconn.PgError

	// Checking if err's message appears in the constraintViolations slice
	// Type assertion under the hood.
	if errors.As(err, &pgxErr) {
		if constraintViolations.Contains(pgxErr.Code) {
			trError = ErrForbiddenInput

			// Checking if err's message appears in the connectionErrors slice
		} else if connectionErrors.Contains(err.Error()) {
			trError = ErrConnection
		}
	}

	// In all other cases, the original error is returned.
	return
}
