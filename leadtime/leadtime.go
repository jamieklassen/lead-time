package leadtime

import (
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 github.com/concourse/concourse/go-concourse/concourse.Client

func LeadTime(repository Repository) time.Duration {
	repository.FindBuild()
	return 1 * time.Second
}
