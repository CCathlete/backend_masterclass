package gapi

import (
	"backend-masterclass/rpc"
	"backend-masterclass/validation"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateCreateUserRequest(req *rpc.CreateUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {

	err := NewPropagatedError()
	if validation.ValidateUsername(req.GetUsername(), err) {
		violations = append(violations, fieldViolation("username", err))
	}

	err = NewPropagatedError()
	if validation.ValidateFullName(req.GetFullName(), err) {
		violations = append(violations, fieldViolation("full_name", err))
	}

	err = NewPropagatedError()
	if validation.ValidateEmail(req.GetEmail(), err) {
		violations = append(violations, fieldViolation("email", err))
	}

	err = NewPropagatedError()
	if validation.ValidatePassword(req.GetPassword(), err) {
		violations = append(violations, fieldViolation("password", err))
	}

	return
}

// TODO: Implement the rest of the validation functions.
func validateLoginUserRequest(req *rpc.LoginUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	return nil
}

func validateUpdateUserRequest(req *rpc.UpdateUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	return nil
}

func validateDeleteUserRequest(req *rpc.DeleteUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	return nil
}

func validateGetUserRequest(req *rpc.GetUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	return nil
}

func validateListUsersRequest(req *rpc.ListUsersRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	return nil
}

func violationsFound(violations []*errdetails.BadRequest_FieldViolation) bool {
	return len(violations) > 0
}

func badRequestError(violations []*errdetails.BadRequest_FieldViolation,
) error {

	badRequestErr := errdetails.BadRequest{
		FieldViolations: violations,
	}

	statusInvalid :=
		status.New(codes.InvalidArgument, "invalid parameters")

	// We try to add the details to the status.
	statusDetails, err := statusInvalid.WithDetails(&badRequestErr)
	if err != nil {
		// If we can't add the details, we return the status without them.
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}
