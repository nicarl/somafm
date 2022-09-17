package prompt

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

func InitScreen() (tcell.Screen, func(), error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, nil, err
	}
	if err := s.Init(); err != nil {
		return nil, nil, err
	}

	s.Clear()
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	return s, quit, nil
}
