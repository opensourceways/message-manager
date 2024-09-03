package domain

type MessageRecipientAdapter interface {
	GetRecipientConfig(countPerPage, pageNum int, userName string) ([]MessageRecipientDTO, int64,
		error)
	AddRecipientConfig(cmd CmdToAddRecipient, userName string) error
	UpdateRecipientConfig(cmd CmdToUpdateRecipient, userName string) error
	RemoveRecipientConfig(cmd CmdToDeleteRecipient, userName string) error
	SyncUserInfo(cmd CmdToSyncUserInfo) (uint, error)
}
