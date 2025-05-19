package jobs

import (
	"context"

	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCourierJob{}

type MoveCourierJob struct {
	command commands.MoveCourierHandler
}

func NewMoveCourierJob(command commands.MoveCourierHandler) (*MoveCourierJob, error) {
	if command == nil {
		return nil, errs.NewValueIsRequiredError("MoveCourierHandler")
	}
	return &MoveCourierJob{
		command: command,
	}, nil
}

func (j *MoveCourierJob) Run() {
	ctx := context.Background()
	command, err := commands.NewMoveCourierCommand()
	if err != nil {
		log.Error("failed to create move courier command: ", err)
	}
	if err := j.command.Handle(ctx, command); err != nil {
		log.Error("failed to handle move courier command: ", err)
	}
}
