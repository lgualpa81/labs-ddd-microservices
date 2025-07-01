package errors

type ErrorCode string

const (
	//User domain errors
	UserNotFound      ErrorCode = "USER_NOT_FOUND"
	UserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	UserInactive      ErrorCode = "USER_INACTIVE"

	//Generic domain errors
	ValidationFailed   ErrorCode = "VALIDATION_FAILED"
	InvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
)
