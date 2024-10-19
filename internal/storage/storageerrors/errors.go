package storageerrors

import "errors"

var (
	ErrUnknownDatabaseType = errors.New("unknown database type")

	ErrUserIsExist           = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidAuthParameters = errors.New("invalid auth parameters")

	ErrCredentialsIsExist  = errors.New("credentials already exists")
	ErrCredentialsNotFound = errors.New("credentials not found")

	ErrTextIsExist  = errors.New("text already exists")
	ErrTextNotFound = errors.New("text not found")

	ErrAttachmentIsExist  = errors.New("attachment already exists")
	ErrAttachmentNotFound = errors.New("attachment not found")

	ErrCardIsExist  = errors.New("card already exists")
	ErrCardNotFound = errors.New("card not found")
)
