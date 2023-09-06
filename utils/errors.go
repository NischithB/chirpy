package utils

import "errors"

var ErrNotFound = errors.New("error: resource not found")

var ErrUserExists = errors.New("auth: user already exists")
var ErrUserNotExists = errors.New("auth: user doesn't exist")
