package service

import (
	"42cli/conf"
)

func Autologin() (string, string, error) {

	user, err := conf.GetString("user_login")
	if err != nil {
		return "", "", err
	}

	pass, err := conf.GetString("passwd_login")
	if err != nil {
		return "", "", err
	}

	if user == "" || pass == "" {
		return "", "", nil
	}

	return user, pass, nil
}