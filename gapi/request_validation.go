package gapi

import (
	"backend-masterclass/rpc"
	"backend-masterclass/validation"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func validateCreateUserRequest(req *rpc.CreateUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {

	err := NewPropagatedError()

	if validation.ValidateUsername(req.GetUsername(), err) {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: GetMessage(err),
		})
	}

	return
}

func validateUpdateUserRequest(req *rpc.UpdateUserRequest) error {
	return nil
}

func validateDeleteUserRequest(req *rpc.DeleteUserRequest) error {
	return nil
}

func validateGetUserRequest(req *rpc.GetUserRequest) error {
	return nil
}

func validateListUsersRequest(req *rpc.ListUsersRequest) error {
	return nil
}
