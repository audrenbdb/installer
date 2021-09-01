package installer

import (
	"github.com/audrenbdb/locale"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//New creates an installer.
//Title provided is going to be installer
//headline inside window GUI. It accepts HTML tags.
//
//From that installer you can :
//
//- add one or more option the user needs to accept
//- add one or more step that will be triggered
//once user accepts conditions.
//
//Once added, you can open installer window
//with the OpenWindow method
func New(title string) *installer {
	i := &installer{
		lang:   lang(locale.GetLang()),
		title:  title,
		height: 540,
		width:  640,
		mustReadAllConditions: true,
	}
	i.setDefaultTexts()
	return i
}

//SetMustReadAllConditions determines if accept condition button will be grayed
//until user scrolled to the bottom of the condition page.
//It is true by default.
func (i *installer) SetMustReadAllConditions(r bool) {
	i.mustReadAllConditions = r
}

//SetOnCloseFunc binds a function to be called once installer
//closes.
//
//It's going to be called in any scenarios, whether
//installer encounters an error or not.
func (i *installer) SetOnCloseFunc(onClose func()) {
	i.onClose = onClose
}

//SetDimensions replaces default dimensions with custom ones
func (i *installer) SetDimensions(width, height int) {
	i.width = width
	i.height = height
}

//AddCondition adds a new condition to be displayed to the user.
//Each condition added this way will be displayed in the same
//order they were added.
func (i *installer) AddCondition(title, body string) {
	i.conditions = append(i.conditions, condition{
		Title: title,
		Body:  body,
	})
}

//AddStep adds a new step to the installer.
//
//A Step is a function that will be called once the user
//accepts conditions displayed to him.
//If multiple steps are added, they will be executed in the
//same order they were added.
//In order to give end-user sense of progression, an artificial
//delay of 2 seconds is added.
func (i *installer) AddStep(process func() error, desc string) {
	i.steps = append(i.steps, step{
		process: func() error {
			time.Sleep(2 * time.Second)
			return process()
		},
		Description: desc,
	})
}

//AddStepRmkDir adds a step that deletes a dir and its child
//before remaking it.
func (i *installer) AddStepRmkDir(dirPath string) {
	process := func() error { return rmkDir(dirPath) }
	desc := i.getRmkDirText(dirPath)
	i.AddStep(process, desc)
}

//AddStepRmvDir adds a step that removes a directory and its child
func (i *installer) AddStepRmvDir(dirPath string) {
	process := func() error { return rmvDir(dirPath) }
	desc := i.getRmvDirText(dirPath)
	i.AddStep(process, desc)
}

//AddStepCopyFiles copy listed files in a given dir
//Files are in a form of map with key being file name
//and value being its content in form of byte array
func (i *installer) AddStepCopyFiles(dirPath string, files map[string][]byte) {
	process := func() error { return copyFiles(dirPath, files) }
	desc := i.getCopyFilesText(dirPath)
	i.AddStep(process, desc)
}

func mkDirAll(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

func rmkDir(dirPath string) error {
	if err := rmvDir(dirPath); err != nil {
		return err
	}
	return mkDirAll(dirPath)
}

func rmvDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func copyFiles(dirPath string, files map[string][]byte) error {
	for fileName, content := range files {
		file := filepath.Join(dirPath, fileName)
		if err := copyFile(file, content); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(file string, content []byte) error {
	return ioutil.WriteFile(file, content, os.ModePerm)
}

func rmvPath(dirPath string) error {
	return os.RemoveAll(dirPath)
}

type step struct {
	process     func() error
	Description string `json:"description"`
}

//condition that has to be accepted by the user to proceed
type condition struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type installer struct {
	title      string
	conditions []condition
	texts      *texts
	lang       lang
	steps      []step
	height     int
	width      int
	//Callback when installer closes
	onClose func()
	//mustReadAllConditions states if a user should scroll to
	//bottom of conditions list
	mustReadAllConditions bool
}
