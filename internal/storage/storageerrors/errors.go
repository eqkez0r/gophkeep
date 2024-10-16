package storageerrors

import "errors"

var (
	ErrUnknownDatabaseType = errors.New("Unknown database type")

	ErrUserIsExist           = errors.New("User already exists")
	ErrUserNotFound          = errors.New("User not found")
	ErrInvalidAuthParameters = errors.New("Invalid auth parameters")

	ErrCredentialsIsExist  = errors.New("Credentials already exists")
	ErrCredentialsNotFound = errors.New("Credentials not found")

	ErrTextIsExist  = errors.New("Text already exists")
	ErrTextNotFound = errors.New("Text not found")

	ErrAttachmentIsExist  = errors.New("Attachment already exists")
	ErrAttachmentNotFound = errors.New("Attachment not found")

	ErrCardIsExist  = errors.New("Card already exists")
	ErrCardNotFound = errors.New("Card not found")
)
