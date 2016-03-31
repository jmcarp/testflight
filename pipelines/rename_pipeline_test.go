package pipelines_test

import (
	"fmt"

	"github.com/concourse/testflight/gitserver"
	"github.com/concourse/testflight/guidserver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Renaming a pipeline", func() {
	var guidServer *guidserver.Server
	var originGitServer *gitserver.Server

	BeforeEach(func() {
		guidServer = guidserver.Start(guidServerRootfs, gardenClient)
		originGitServer = gitserver.Start(gitServerRootfs, gardenClient)
	})

	AfterEach(func() {
		guidServer.Stop()
		originGitServer.Stop()
	})

	It("runs scheduled after pipeline is renamed", func() {
		configurePipeline(
			"-c", "fixtures/simple-trigger.yml",
			"-v", "testflight-helper-image="+guidServerRootfs,
			"-v", "guid-server-curl-command="+guidServer.RegisterCommand(),
			"-v", "origin-git-server="+originGitServer.URI(),
		)

		guid1 := originGitServer.Commit()
		Eventually(guidServer.ReportingGuids).Should(ContainElement(guid1))

		newPipelineName := fmt.Sprintf("renamed-test-pipeline-%d", GinkgoParallelNode())
		renamePipeline(newPipelineName)

		guid2 := originGitServer.Commit()
		Eventually(guidServer.ReportingGuids).Should(ContainElement(guid2))
	})
})
