package main

import (
	"fmt"
)

type NotFound struct {
	Msg string `json:"msg"`
}

func (e NotFound) Error() string {
	return e.Msg
}

func CreateNotFoundError(id string) error {
	return NotFound{fmt.Sprintf("Resource %s is not found", id)}
}

func IsNotFoundError(err error) bool {
	switch err.(type) {
	case NotFound:
		return true
	default:
		return false
	}
}
