package dao

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")

	ErrCompanyAlreadyExists = errors.New("company already exists")
	ErrCompanyNotFound      = errors.New("company not found")

	ErrProfilePictureNotFound = errors.New("profile picture not found")
	ErrInvalidProfilePicture  = errors.New("invalid profile picture")
)
