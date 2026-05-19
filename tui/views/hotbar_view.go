package views

import (
	"fmt"
	"charm.land/lipgloss/v2"

	"42cli/conf"
)

func Hotbar() string {

	loggedUser, err := conf.GetString("logged_with")
	if err != nil || loggedUser == "" {
		return ""
	}

	line := lipgloss.JoinHorizontal(lipgloss.Left,fmt.Sprintf("Logged with %s", loggedUser))
	return line
}