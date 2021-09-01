package installer

import (
	_ "embed"
	"errors"
	"github.com/wailsapp/wails"
)

//go:embed main.js
var js string

//go:embed styles.css
var css string

type wailsBind struct {
	Title      string      `json:"title"`
	Conditions []condition `json:"conditions"`
	Steps      []step      `json:"steps"`
	Texts      *texts      `json:"texts"`
	MustReadAllConditions bool `json:"mustReadAllConditions"`

	//completed is set to true when all steps have
	//been processed successfully
	completed bool
}

//OpenWindow open the GUI installer windows.
//Steps and conditions must have been set prior
//to opening the window.
func (i *installer) OpenWindow(windowTitle string) error {
	err := i.newWailsApp(windowTitle)
	if err != nil {
		return err
	}
	if i.onClose != nil {
		i.onClose()
	}
	return nil
}

func (i *installer) newWailsApp(title string) error {
	app := wails.CreateApp(i.newWailsAppConfig(title))
	bind := i.newWailsBind()
	app.Bind(bind)
	if err := app.Run(); err != nil {
		return err
	}
	if !bind.completed {
		return errors.New("all steps not completed")
	}
	return nil
}

func (i *installer) newWailsBind() *wailsBind {
	return &wailsBind{
		Title:      i.title,
		Conditions: i.conditions,
		Steps:      i.steps,
		Texts:      i.texts,
		MustReadAllConditions: i.mustReadAllConditions,
	}
}

func (i *installer) newWailsAppConfig(title string) *wails.AppConfig {
	return &wails.AppConfig{
		Resizable: true,
		Width:     i.width,
		Height:    i.height,
		Title:     title,
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
	}
}

func (g *wailsBind) Self() *wailsBind {
	return g
}

func (g *wailsBind) InstallStep(i int) error {
	lastIndex := len(g.Steps) - 1
	err := g.Steps[i].process()
	if err != nil {
		return err
	}
	if i == lastIndex {
		g.completed = true
	}
	return nil
}
