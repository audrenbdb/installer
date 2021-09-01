package installer

import "fmt"

type lang string

const (
	fr lang = "fr"
	en lang = "en"
	vi lang = "vi"
)

type texts struct {
	AcceptButton   string `json:"acceptButton"`
	Fail           string `json:"fail"`
	Success        string `json:"success"`
	CompletedSteps string `json:"completedSteps"`
	ReadAllConditionsTooltip string `json:"readAllConditionsTooltip"`
}

func (i *installer) setDefaultTexts() {
	i.texts = &texts{
		AcceptButton:   i.getAcceptButtonText(),
		Fail:           i.getInstallationFailText(),
		Success:        i.getInstallationSuccessText(),
		CompletedSteps: i.getCompletedStepsText(),
		ReadAllConditionsTooltip: i.getReadAllConditionsToolTip(),
	}
}

func (i *installer) SetAcceptButtonText(txt string) {
	i.texts.AcceptButton = txt
}

func (i *installer) SetFailText(txt string) {
	i.texts.Fail = txt
}

func (i *installer) SetSuccessText(txt string) {
	i.texts.Success = txt
}

func (i *installer) getRegisterSchemeText(protocol string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Nous installons le schema %s."
	case vi:
		msg = "Cài đặt chương trình %s."
	default:
		msg = "Installation of scheme %s."
	}
	return fmt.Sprintf(msg, protocol)
}

func (i *installer) getCopyFilesText(dirPath string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Des fichiers nécessaires seront installés ici : %s."
	case vi:
		msg = "Các tệp cần thiết sẽ được cài đặt tại đây : %s."
	default:
		msg = "Required files will be installed here : %s."
	}
	return fmt.Sprintf(msg, dirPath)
}

func (i *installer) getRmkDirText(dirPath string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Le dossier %s va être créé."
	case vi:
		msg = "Thư mục %s đang được tạo."
	default:
		msg = "Directory %s is being created."
	}
	return fmt.Sprintf(msg, dirPath)
}

func (i *installer) getRmvDirText(dirPath string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Suppression du dossier %s."
	case vi:
		msg = "xóa thư mục %s."
	default:
		msg = "Deleting folder %s."
	}
	return fmt.Sprintf(msg, dirPath)
}

func (i *installer) getUninstallOptText() string {
	switch i.lang {
	case fr:
		return "Ajout d'une option de désinstallation."
	case vi:
		return "Thêm tùy chọn gỡ cài đặt."
	default:
		return "Adding uninstall option."
	}
}

func (i *installer) getRemoveUninstallOptText() string {
	switch i.lang {
	case fr:
		return "Suppression de l'option de désinstallation."
	case vi:
		return "Xóa tùy chọn gỡ cài đặt."
	default:
		return "Removing uninstall option."
	}
}

func (i *installer) getUnregisterSchemeText(protoc string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Suppression du scheme %s."
	case vi:
		msg = "Xóa lược đồ %s."
	default:
		msg = "Deleting scheme %s."
	}
	return fmt.Sprintf(msg, protoc)
}

func (i *installer) getAcceptButtonText() string {
	switch i.lang {
	case fr:
		return "J'ai lu et j'accepte"
	case vi:
		return "Tôi đã đọc và tôi chấp nhận"
	default:
		return "I have read and I accept"
	}
}

func (i *installer) getRemoveFolderText(path string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Le dossier %s va être supprimé."
	case vi:
		msg = "Thư mục %s sẽ bị xóa."
	default:
		msg = "Folder %s is going to be deleted."
	}
	return fmt.Sprintf(msg, path)
}

func (i *installer) getRemoveFolderAfterInstallText(path string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Le dossier : %s sera supprimé quelques secondes après la fermeture de cette fenêtre."
	case vi:
		msg = "Thư mục TOTO sẽ bị xóa vài giây sau khi đóng cửa sổ này."
	default:
		msg = "The folder : %s will be deleted a few seconds after closing this window."
	}
	return fmt.Sprintf(msg, path)
}

func (i *installer) getInstallationSuccessText() string {
	switch i.lang {
	case fr:
		return "<p>Le processus s'est déroulé correctement jusqu'à son terme.</p><p><b>Vous pouvez fermer cette fenêtre.</b></p><p>Si vous le souhaitez, vous pouvez voir l'historique'des étapes achevées via le bouton ci-dessous.</p>"
	case vi:
		return "<p>Quá trình diễn ra suôn sẻ để hoàn tất.</p><p><b>Bạn có thể đóng cửa sổ này.</b></p><p>Nếu muốn, bạn có thể xem lịch sử của các bước đã hoàn thành qua nút bên dưới.</p>"
	default:
		return "<p>The process went smoothly to completion.</p><p><b>You may close this window.</b></p><p>If you want, you may see the history of the steps completed via the button below.</p>"
	}
}

func (i *installer) getCompletedStepsText() string {
	switch i.lang {
	case fr:
		return "Étapes réalisées"
	case vi:
		return "Các bước đã hoàn thành"
	default:
		return "Steps completed"
	}
}

func (i *installer) getInstallationFailText() string {
	switch i.lang {
	case fr:
		return "Le processus a rencontré une erreur et n'a pu arriver à son terme."
	case vi:
		return "Quá trình gặp lỗi. Cửa sổ này sẽ tự đóng sau vài giây."
	default:
		return "Process encountered an error and could not complete."
	}
}

func (i *installer) getShortcutCreatingText(src, dest string) string {
	var msg string
	switch i.lang {
	case fr:
		msg = "Un raccourci de %s sera créé ici : %s."
	case vi:
		msg = "Một lối tắt từ %s sẽ được tạo ở đây: %s."
	default:
		msg = "A shortcut from %s is going to be created here : %s."
	}
	return fmt.Sprintf(msg, src, dest)
}

func (i *installer) getReadAllConditionsToolTip() string {
	switch i.lang {
	case fr:
		return "Vous devez avoir fait défiler l'ensemble des conditions pour accepter"
	case vi:
		return "Bạn phải cuộn qua tất cả các điều kiện để tiếp tục"
	default:
		return "You must have scrolled through all of the conditions to continue"
	}
}