package views

import (
	"os"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
)

func Autologin() string {
	ascii := `
	  .o     .oooo.             oooo   o8o  
    .d88   .dP""Y88b            '888   '"'  
  .d'888         ]8P'  .ooooo.   888  oooo  
.d'  888       .d8P'  d88' '"Y8  888  '888  
88ooo888oo   .dP'     888        888   888  
     888   .oP     .o 888   .o8  888   888  
    o888o  8888888888 'Y8bod8P' o888o o888o
`
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	root := lipgloss.NewStyle().
		Width(w).
		AlignHorizontal(lipgloss.Center).
		Padding(2, 0)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#688990")).
		Bold(true)


	return root.Render(style.Render(ascii),)
}