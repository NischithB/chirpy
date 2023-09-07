package utils

import "errors"

var ErrNotFound = errors.New("error: resource not found")
var ErrForbidden = errors.New("error: Not authorized to modify resource")

var ErrUserExists = errors.New("auth: user already exists")
var ErrUserNotExists = errors.New("auth: user doesn't exist")
var ErrTokenMissing = errors.New("auth: token missing from the request header")
