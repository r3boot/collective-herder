package utils

import (
	"github.com/satori/go.uuid"
)

func Uuidgen() string {
	return uuid.NewV4().String()
}
