package leadtime

import (
	"fmt"
	"reflect"
	"time"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
)

type Input struct {
	Name     string
	Resource string
}

type Repository interface {
	Nodes(string, string) []Node
	InputsForNode(Node) []Input
	LinkForInput(Node, Input) (Link, error)
	NodesForLink(Link) []Node
}

func NewRepository(client concourse.Client) Repository {
	return &apiRepository{
		client:       client,
		teamName:     "main",
		pipelineName: "concourse",
	}
}

type apiRepository struct {
	client       concourse.Client
	teamName     string
	pipelineName string
}

type Node struct {
	ID        int
	Start     time.Time
	Finish    time.Time
	Job       string
	Succeeded bool
}

func NewNode(build atc.Build) Node {
	return Node{
		ID:        build.ID,
		Finish:    time.Unix(build.EndTime, 0),
		Start:     time.Unix(build.StartTime, 0),
		Job:       build.JobName,
		Succeeded: build.Status == "succeeded",
	}
}

func (r *apiRepository) Nodes(pipelineName, jobName string) []Node {
	nodes := []Node{}
	builds, _, _, _ := r.client.Team(r.teamName).JobBuilds(
		pipelineName,
		jobName,
		concourse.Page{},
	)
	for _, build := range builds {
		nodes = append(nodes, NewNode(build))
	}
	return nodes
}

func (r *apiRepository) InputsForNode(node Node) []Input {
	job, _, _ := r.client.Team(r.teamName).Job(r.pipelineName, node.Job)
	inputs := []Input{}
	for _, jobInput := range job.Inputs {
		if len(jobInput.Passed) != 0 {
			inputs = append(inputs, Input{
				Name:     jobInput.Name,
				Resource: jobInput.Resource,
			})
		}
	}
	return inputs
}

type Link struct {
	Resource string
	ID       int
}

func (r *apiRepository) LinkForInput(node Node, input Input) (Link, error) {
	resources, _, _ := r.client.BuildResources(node.ID)
	var version atc.Version
	for _, buildInput := range resources.Inputs {
		if buildInput.Name == input.Name {
			version = buildInput.Version
			break
		}
	}
	if version == nil {
		return Link{}, fmt.Errorf(
			"couldn't find build input to match job input %s",
			input.Name,
		)
	}
	versions, _, _, _ := r.client.Team(r.teamName).ResourceVersions(
		r.pipelineName,
		input.Resource,
		concourse.Page{},
	)
	for _, v := range versions {
		if reflect.DeepEqual(v.Version, version) {
			return Link{Resource: input.Resource, ID: v.ID}, nil
		}
	}
	return Link{}, fmt.Errorf(
		"couldn't find resource version matching build input: %v",
		version,
	)
}

func (r *apiRepository) NodesForLink(link Link) []Node {
	nodes := []Node{}
	builds, _, _ := r.client.Team(r.teamName).BuildsWithVersionAsOutput(
		r.pipelineName,
		link.Resource,
		link.ID,
	)
	for _, build := range builds {
		nodes = append(nodes, NewNode(build))
	}
	builds, _, _ = r.client.Team(r.teamName).BuildsWithVersionAsInput(
		r.pipelineName,
		link.Resource,
		link.ID,
	)
	for _, build := range builds {
		nodes = append(nodes, NewNode(build))
	}
	return nodes
}
