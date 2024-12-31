package gapi

import (
	"backend-masterclass/db/sqlc"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error) error {

	if errors.Is(err, sqlc.ErrRecordNotFound) {
		return status.Error(codes.FailedPrecondition, err.Error())

	} else if errors.Is(err, sqlc.ErrConnection) {
		return status.Error(codes.Unavailable, err.Error())

	} else if errors.Is(err, sqlc.ErrForbiddenInput) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// Any other error.
	return status.Error(codes.Internal, err.Error())
}
