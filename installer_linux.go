package installer

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	shareAppPath     = ".local/share/applications"
	xdgMime          = "xdg-mime"
	xdgSchemeHandler = "x-scheme-handler"
	desktopExt       = ".desktop"
)

//AddStepCreateScheme registers a new scheme by
//copying a customprotocol.desktop file in application dir
//located in ~/.local/shared/applications.
//A desktop file contains a link or a bash cmd that will
//be triggered when scheme is called
func (i *installer) AddStepCreateScheme(protoc string, content []byte) {
	process := func() error { return createScheme(protoc, content) }
	description := i.getRegisterSchemeText(protoc)
	i.AddStep(process, description)
}

func (i *installer) AddStepDeleteScheme(scheme string) {
	process := func() error { return deleteScheme(scheme) }
	description := i.getUnregisterSchemeText(scheme)
	i.AddStep(process, description)
}

func deleteScheme(scheme string) error {
	return deleteDesktopFile(getSchemeDesktopFileName(scheme))
}

func deleteDesktopFile(file string) error {
	path, err := getDotDesktopFilePath(file)
	if err != nil {
		return err
	}
	return rmvPath(path)
}

func getShareAppFullPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, shareAppPath), nil
}

func getDotDesktopFilePath(file string) (string, error) {
	p, err := getShareAppFullPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, file), nil
}

func getSchemeDesktopFileName(scheme string) string {
	return scheme + desktopExt
}

func createScheme(protoc string, content []byte) error {
	desktopFile := protoc + desktopExt
	if err := mkAllShareAppDirPath(); err != nil {
		return err
	}
	err := copyDesktopFile(desktopFile, content)
	if err != nil {
		return err
	}
	return runXDGMime(desktopFile, protoc)
}

func mkAllShareAppDirPath() error {
	path, err := getShareAppFullPath()
	if err != nil {
		return err
	}
	return mkDirAll(path)
}

//runXDGMime calls xdg mime app to bind a scheme to a .desktop file
func runXDGMime(desktopFile, scheme string) error {
	handler := filepath.Join(xdgSchemeHandler, scheme)
	cmd := exec.Command(xdgMime, "default", desktopFile, handler)
	return cmd.Run()
}

func copyDesktopFile(file string, content []byte) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(home, shareAppPath, file)
	return copyFile(filePath, content)
}
