package gapi

import (
	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropagatedError = u.PropagatedError

// Function aliases.
var (
	WrapError          = u.WrapError
	NewPropagatedError = u.NewPropagatedError
	GetMessage         = u.GetMessage
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

func fieldViolation(field string, err PropagatedError,
) *errdetails.BadRequest_FieldViolation {

	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: GetMessage(err),
	}
}
