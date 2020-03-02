package leadtime_test

import (
	"time"

	"github.com/concourse/concourse/atc"
	leadtimefakes "github.com/concourse/leadtime/lead-timefakes"
	"github.com/concourse/leadtime/leadtime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build", func() {
	var fakeClient *leadtimefakes.FakeClient
	var build leadtime.Build

	BeforeEach(func() {
		fakeClient = new(leadtimefakes.FakeClient)
		build = leadtime.NewBuild(fakeClient, atc.Build{
			ID:           0,
			TeamName:     "main",
			Name:         "0",
			Status:       "succeeded",
			JobName:      "publish-binaries",
			APIURL:       "/api/v1/builds/0",
			PipelineName: "concourse",
			StartTime:    0,
			EndTime:      time.Now().Unix(),
			ReapTime:     0,
		})
	})

	Context("#Neighbours", func() {
		BeforeEach(func() {
			fakeClient.BuildResourcesReturns(atc.BuildInputsOutputs{}, true, nil)
		})

		It("looks up resources", func() {
			build.Neighbours()
			Expect(fakeClient.BuildResourcesCallCount()).To(Equal(1))
			Expect(fakeClient.BuildResourcesArgsForCall(0)).To(Equal(0))
		})
	})
})
