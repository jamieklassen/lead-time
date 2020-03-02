package leadtime

import (
	"time"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
)

type Build interface {
	Done() time.Time
	Neighbours() []Build
}

func NewBuild(client concourse.Client, atcBuild atc.Build) Build {
	return &build{
		client: client,
		id:     atcBuild.ID,
		done:   time.Unix(atcBuild.EndTime, 0),
	}
}

type build struct {
	client concourse.Client
	id     int
	done   time.Time
}

func (b *build) Done() time.Time {
	return b.done
}
func (b *build) Neighbours() []Build {
	b.client.BuildResources(b.id)
	return []Build{}
}
