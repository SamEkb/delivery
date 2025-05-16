package commands

type MoveCourierCommand struct {
	isValid bool
}

func NewMoveCourierCommand() (*MoveCourierCommand, error) {
	return &MoveCourierCommand{
		isValid: true,
	}, nil
}

func (c *MoveCourierCommand) IsValid() bool {
	return c.isValid
}
