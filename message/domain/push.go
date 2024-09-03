package domain

type MessagePushAdapter interface {
	GetPushConfig(subsIds []string, countPerPage, pageNum int,
		userName string) ([]MessagePushDTO, error)
	AddPushConfig(cmd CmdToAddPushConfig) error
	UpdatePushConfig(cmd CmdToUpdatePushConfig) error
	RemovePushConfig(cmd CmdToDeletePushConfig) error
}
