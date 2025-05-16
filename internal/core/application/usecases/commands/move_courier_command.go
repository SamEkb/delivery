package commands

type MoveCourierCommand struct {
	isValid bool
}

func NewMoveCourierCommand() (*MoveCourierCommand, error) {
	return &MoveCourierCommand{
		isValid: true,
	}, nil
}
