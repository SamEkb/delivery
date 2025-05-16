package commands

type AssignOrderCommand struct {
	isValid bool
}

func NewAssignOrderCommand() (*AssignOrderCommand, error) {
	return &AssignOrderCommand{
		isValid: true,
	}, nil
}

func (c *AssignOrderCommand) IsValid() bool {
	return c.isValid
}
