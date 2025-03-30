package storage

import "errors"

var (
	ErrLoginNotFound = errors.New("Login not found")
	ErrURLExists   = errors.New("URL exists")
)
