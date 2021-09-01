package installer

import (
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	createShortcutShellCmdFmt = `
		$SourceFileLocation = "%s"
		$ShortcutLocation = "%s"
		$WScriptShell = New-Object -ComObject WScript.Shell
		$Shortcut = $WScriptShell.CreateShortcut($ShortcutLocation)
		$Shortcut.TargetPath = $SourceFileLocation
		$Shortcut.Save()
	`
)

func (i *installer) AddStepCreateShortcut(src, dst string) {
	process := func() error { return createShortcut(src, dst) }
	description := i.getShortcutCreatingText(src, dst)
	i.AddStep(process, description)
}

//AddStepRmvFolderAfterInstall sets an onClose function to delete
//a path and its subpaths.
//Process is empty because it's not an actual immediatly processed function
//but rather a delayed one.
func (i *installer) AddStepRmvFolderAfterInstall(path string) {
	process := func() error { return nil }
	description := i.getRemoveFolderAfterInstallText(path)
	i.onClose = func() {
		err := rmvFolderAfterDelay(path)
		if err != nil {
			log.Fatal(err)
		}
	}
	i.AddStep(process, description)
}

func rmvFolderAfterDelay(path string) error {
	cmd := fmt.Sprintf("Start-Sleep -s 5; rm -r %s", path)
	return startHiddenPowerShellCmd(cmd)
}

//AddStepCreateScheme adds a step that creates registry keys to handle a custom scheme.
//shellCmd will be executed when scheme is called.
//Please note that registry keys are user related and not global.
//FriendlyTypeName is the name displayed to the user when he attempts to open that scheme.
func (i *installer) AddStepCreateScheme(scheme string, friendlyTypeName string, shellCmd string) {
	process := func() error {
		return createScheme(scheme, friendlyTypeName, shellCmd)
	}
	description := i.getRegisterSchemeText(scheme)
	i.AddStep(process, description)
}

//UninstallOptions is used to create a registry key with optional options provided
type UninstallOptions struct {
	DisplayName    string
	DisplayVersion string
	Publisher      string
	//UninstallString is the path to the uninstaller
	UninstallString string
	URLInfoAbout    string
	//Path to the display icon. Usually main binary exe who has already an icon set
	DisplayIcon string
	//KeyName to be found in SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\KeyName
	KeyName string
}

func (i *installer) AddStepCreateUninstallOpt(opts UninstallOptions) {
	process := func() error {
		return createUninstallOpt(opts)
	}
	description := i.getUninstallOptText()
	i.AddStep(process, description)
}

//AddStepDeleteScheme adds a step that deletes scheme association registry keys
func (i *installer) AddStepDeleteScheme(scheme string) {
	process := func() error { return deleteSchemeKey(scheme) }
	description := i.getUnregisterSchemeText(scheme)
	i.AddStep(process, description)
}

//AddStepDeleteUninstallOpt deletes uninstall registry keys associated
//with a given program.
func (i *installer) AddStepDeleteUninstallOpt(prog string) {
	process := func() error { return deleteUninstallKey(prog) }
	description := i.getRemoveUninstallOptText()
	i.AddStep(process, description)
}

func createShortcut(src, dst string) error {
	shellCmd := fmt.Sprintf(createShortcutShellCmdFmt, src, dst)
	if err := startHiddenPowerShellCmd(shellCmd); err != nil {
		return err
	}
	return nil
}

func startHiddenPowerShellCmd(shellCmd string) error {
	p, err := getPowershellPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(p, shellCmd)
	hidePowershellCmd(cmd)
	return cmd.Start()
}

func getPowershellPath() (string, error) {
	w, err := windows.GetSystemWindowsDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(w, "System32", "WindowsPowershell", "v1.0", "powershell.exe"), nil
}

func hidePowershellCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
}

func windowsSoftwareKeyPath() string {
	return filepath.Join("SOFTWARE", "Microsoft",
		"Windows", "CurrentVersion")
}

func createUninstallOpt(opts UninstallOptions) error {
	k, err := createUninstallProgKey(opts.KeyName)
	if err != nil {
		return err
	}
	defer k.Close()
	return setUninstKeyValues(k, opts)
}

func setUninstKeyValues(k registry.Key, opts UninstallOptions) error {
	return setRegistryKeyValues(k, map[string]string{
		"NoModify":        "1",
		"NoRepair":        "1",
		"UninstallString": opts.UninstallString,
		"DisplayName":     opts.DisplayName,
		"DisplayIcon":     opts.DisplayIcon,
		"Publisher":       opts.Publisher,
		"URLInfoAbout":    opts.URLInfoAbout,
		"DisplayVersion":  opts.DisplayVersion,
	})
}

func setRegistryKeyValues(k registry.Key, values map[string]string) error {
	for name, val := range values {
		err := k.SetStringValue(name, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func createScheme(scheme, friendlyTypeName, shellCmd string) error {
	schemeK, err := createURLSchemeKey(scheme, friendlyTypeName)
	if err != nil {
		return err
	}
	defer schemeK.Close()
	return createShellCmdKey(schemeK, friendlyTypeName, shellCmd)
}

func deleteSchemeKey(scheme string) error {
	classK, err := openSwClassKey()
	if err != nil {
		return err
	}
	defer classK.Close()
	keys := append([]string{scheme}, "shell", "open", "command")
	return deleteKeys(classK, keys)
}

func deleteUninstallKey(prog string) error {
	uninstallK, err := openUninstallKey()
	if err != nil {
		return err
	}
	defer uninstallK.Close()
	return registry.DeleteKey(uninstallK, prog)
}

func openSwClassKey() (registry.Key, error) {
	return openKey(registry.CURRENT_USER, classKeyPath())
}

func openUninstallKey() (registry.Key, error) {
	return openKey(registry.CURRENT_USER, uninstallKeyPath())
}

func createUninstallProgKey(prog string) (registry.Key, error) {
	return addKey(registry.CURRENT_USER, uninstallProgKeyPath(prog))
}

func createURLSchemeKey(scheme, friendlyTypeName string) (registry.Key, error) {
	k, err := addKey(registry.CURRENT_USER, schemeKeyPath(scheme))
	if err != nil {
		return k, err
	}
	return k, setRegistryKeyValues(k, map[string]string{
		"URL Protocol":     "",
		"FriendlyTypeName": friendlyTypeName + " " + "Protocol",
	})
	return k, k.SetStringValue("URL Protocol", "")
}

func createShellCmdKey(sourceK registry.Key, friendlyTypeName, shellCmd string) error {
	k, err := createShellOpenKey(sourceK, friendlyTypeName)
	if err != nil {
		return err
	}
	cmdK, err := addKey(k, "command")
	if err != nil {
		return err
	}
	defer cmdK.Close()
	return cmdK.SetStringValue("", shellCmd)
}

func createShellOpenKey(sourceK registry.Key, friendlyTypeName string) (registry.Key, error) {
	p := filepath.Join("shell", "open")
	openK, err := addKey(sourceK, p)
	if err != nil {
		return openK, err
	}
	return openK, openK.SetStringValue("FriendlyAppName", friendlyTypeName)
}

func classKeyPath() string {
	return filepath.Join("SOFTWARE", "Classes")
}

func schemeKeyPath(scheme string) string {
	return filepath.Join(classKeyPath(), scheme)
}

func uninstallKeyPath() string {
	return filepath.Join(windowsSoftwareKeyPath(), "Uninstall")
}

func uninstallProgKeyPath(prog string) string {
	return filepath.Join(uninstallKeyPath(), prog)
}

func addKey(sourceK registry.Key, newKPath string) (registry.Key, error) {
	newK, _, err := registry.CreateKey(sourceK, newKPath, registry.ALL_ACCESS)
	return newK, err
}

func openKey(sourceK registry.Key, path string) (registry.Key, error) {
	k, _, err := registry.CreateKey(sourceK, path, registry.ALL_ACCESS)
	return k, err
}

func deleteKeys(sourceK registry.Key, subKeyNames []string) error {
	if len(subKeyNames) == 0 {
		return nil
	}
	err := registry.DeleteKey(sourceK, filepath.Join(subKeyNames...))
	if err != nil {
		return err
	}
	return deleteKeys(sourceK, subKeyNames[:len(subKeyNames)-1])
}
