package id

import "github.com/lithammer/shortuuid/v4"

func New() string {
	return shortuuid.New()
}

func IsValid(value string) bool {
	if len(value) != 22 {
		return false
	}
	_, err := shortuuid.DefaultEncoder.Decode(value)
	return err == nil
}
