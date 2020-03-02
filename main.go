package main

import (
	"fmt"
	"net/http"

	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/concourse/leadtime/leadtime"
)

const url = "https://ci.concourse-ci.org"

func main() {
	client := concourse.NewClient(url, http.DefaultClient, false)
	repository := leadtime.NewRepository(client)
	leadTime := leadtime.LeadTime(repository)
	fmt.Println(leadTime)
}
