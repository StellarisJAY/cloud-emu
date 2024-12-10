// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_NOT_FOUND.String() && e.Code == 404
}

func ErrorNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, ErrorReason_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsAccessDenied(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ACCESS_DENIED.String() && e.Code == 403
}

func ErrorAccessDenied(format string, args ...interface{}) *errors.Error {
	return errors.New(403, ErrorReason_ACCESS_DENIED.String(), fmt.Sprintf(format, args...))
}

func IsLoginFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_LOGIN_FAILED.String() && e.Code == 501
}

func ErrorLoginFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(501, ErrorReason_LOGIN_FAILED.String(), fmt.Sprintf(format, args...))
}
