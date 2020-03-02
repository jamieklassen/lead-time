package leadtime

import (
	"github.com/concourse/concourse/go-concourse/concourse"
)

type Repository interface {
	FindBuild() Build
}

func NewRepository(client concourse.Client) Repository {
	return &repository{
		client:       client,
		teamName:     "main",
		pipelineName: "concourse",
		jobName:      "publish-binaries",
	}
}

type repository struct {
	client       concourse.Client
	teamName     string
	pipelineName string
	jobName      string
}

func (r *repository) FindBuild() Build {
	bs, _, _, _ := r.client.Team(r.teamName).JobBuilds(
		r.pipelineName,
		r.jobName,
		concourse.Page{Limit: 1},
	)
	return NewBuild(r.client, bs[0])
}
