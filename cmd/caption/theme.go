package main

import (
	"fyne.io/fyne"
)

type myTheme struct {
	fyne.Theme
}

func (myTheme) TextFont() fyne.Resource { return resourceMplus1pRegularTtf }
