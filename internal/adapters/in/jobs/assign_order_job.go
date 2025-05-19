package jobs

import (
	"context"

	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrderJob{}

type AssignOrderJob struct {
	command commands.AssignOrderHandler
}

func NewAssignOrderJob(command commands.AssignOrderHandler) (*AssignOrderJob, error) {
	if command == nil {
		return nil, errs.NewValueIsRequiredError("AssignOrderHandler")
	}
	return &AssignOrderJob{
		command: command,
	}, nil
}

func (j *AssignOrderJob) Run() {
	ctx := context.Background()
	command, err := commands.NewAssignOrderCommand()
	if err != nil {
		log.Error("failed to create assign order command: ", err)
	}
	if err := j.command.Handle(ctx, command); err != nil {
		log.Error("failed to handle assign order command: ", err)
	}
}
