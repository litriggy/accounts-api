package utils

import (
	"net/mail"
	"regexp"
	"strconv"
)

func IsMail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

func IsNum(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil

}
