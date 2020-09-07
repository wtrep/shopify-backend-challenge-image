/*
  This file defines commonly used errors that are shared by the frontend and the backend
 */

package common

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error_ *ErrorResponseError `json:"error,omitempty"`
}

type ErrorResponseError struct {
	// Id of the error
	Id int32 `json:"id,omitempty"`
	// Name of the error
	Name string `json:"name,omitempty"`
	// Details related to the error
	Detail string `json:"detail,omitempty"`
	// HTTP Code related to the error
	Code int32 `json:"-"`
}

var InvalidRequestBodyError = ErrorResponseError{
	Id:     1200,
	Name:   "InvalidRequestBodyError",
	Detail: "The request body doesn't respect the valid format",
	Code:   http.StatusBadRequest,
}

var UserDoesNotExistError = ErrorResponseError{
	Id:     1201,
	Name:   "UserDoesNotExistError",
	Detail: "The provided user doesn't exist",
	Code:   http.StatusNotFound,
}

var WrongPasswordError = ErrorResponseError{
	Id:     1202,
	Name:   "WrongPasswordError",
	Detail: "The provided password doesn't match our records",
	Code:   http.StatusUnauthorized,
}

var DatabaseInsertionError = ErrorResponseError{
	Id:     1203,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var JSONEncoderError = ErrorResponseError{
	Id:     1204,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var PasswordTooLongError = ErrorResponseError{
	Id:     1205,
	Name:   "PasswordTooLongError",
	Detail: "Password length is more than 32 characters",
	Code:   http.StatusBadRequest,
}

var UserAlreadyExistError = ErrorResponseError{
	Id:     1206,
	Name:   "UserAlreadyExistError",
	Detail: "The provided user already exist",
	Code:   http.StatusConflict,
}

var MissingTokenError = ErrorResponseError{
	Id:     1207,
	Name:   "MissingTokenError",
	Detail: "No bearer token was provided",
	Code:   http.StatusUnauthorized,
}

var InvalidTokenError = ErrorResponseError{
	Id:     1208,
	Name:   "InvalidTokenError",
	Detail: "The bearer token provided is invalid",
	Code:   http.StatusBadRequest,
}

var TokenGenerationError = ErrorResponseError{
	Id:     1209,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var WrongUserError = ErrorResponseError{
	Id:     1211,
	Name:   "WrongUserError",
	Detail: "The user associated with that token doesn't match the image owner",
	Code:   http.StatusUnauthorized,
}

var InvalidImageBodyError = ErrorResponseError{
	Id:     1212,
	Name:   "InvalidImageBodyError",
	Detail: "Invalid image body. Maximum file size is 10 MB",
	Code:   http.StatusUnauthorized,
}

var FileUploadError = ErrorResponseError{
	Id:     1213,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var URLGenerationError = ErrorResponseError{
	Id:     1214,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var InvalidUUIDError = ErrorResponseError{
	Id:     1215,
	Name:   "InvalidUUIDError",
	Detail: "The provided uuid is invalid",
	Code:   http.StatusBadRequest,
}

var ImageNotFoundError = ErrorResponseError{
	Id:     1216,
	Name:   "ImageNotFoundError",
	Detail: "No image was found with the provided uuid",
	Code:   http.StatusNotFound,
}

var ImageNotUploadedError = ErrorResponseError{
	Id:     1217,
	Name:   "ImageNotUploadedError",
	Detail: "The image was found but isn't currently uploaded to the service",
	Code:   http.StatusNotFound,
}

var UserPermissionDeniedError = ErrorResponseError{
	Id:     1218,
	Name:   "UserPermissionDeniedError",
	Detail: "You aren't the owner of that image",
	Code:   http.StatusUnauthorized,
}

var FileDeletionError = ErrorResponseError{
	Id:     1219,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var DBDeletionError = ErrorResponseError{
	Id:     1220,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var GetImagesDBError = ErrorResponseError{
	Id:     1221,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

func RespondWithError(w http.ResponseWriter, error *ErrorResponseError) {
	w.WriteHeader(int(error.Code))
	response := ErrorResponse{
		Error_: error,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		http.Error(w, "Server unhandled error", int(error.Code))
	}
}
