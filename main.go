package main

import (
	"fmt"
	"net/http"

	"github.com/concourse/concourse/fly/rc"
	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/concourse/leadtime/leadtime"
)

const url = "https://ci.concourse-ci.org"
const authenticate = true

func main() {
	var client concourse.Client
	if authenticate {
		target, _, _ := rc.LoadTargetFromURL(url, "main", false)
		client = target.Client()
	} else {
		client = concourse.NewClient(url, http.DefaultClient, false)
	}
	repository := leadtime.NewRepository(client)
	leadTime := leadtime.LeadTime(repository)
	fmt.Println(leadTime)
}
