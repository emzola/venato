package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	internalServerError = status.Error(codes.Internal, "the server encountered a problem and could not process your request")
	notFoundError       = status.Error(codes.NotFound, "the requested resource could not be found")
	nilRequestError     = status.Error(codes.InvalidArgument, "nil request")
	editConflictError   = status.Error(codes.AlreadyExists, "unable to update the record due to an edit conflict, please try again")
)

// failedValidationError returns a failed validation error message.
func (h *Handler) failedValidationError(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}
