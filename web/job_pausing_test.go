package web_test

import (
	"fmt"
	"time"

	"github.com/concourse/atc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("JobPausing", func() {
	var (
		build          atc.Build
		loadingTimeout time.Duration
	)

	Context("with a job in the configuration", func() {
		BeforeEach(func() {
			_, _, _, err := client.CreateOrUpdatePipelineConfig(pipelineName, "0", atc.Config{
				Jobs: []atc.JobConfig{
					{Name: "some-job"},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = client.UnpausePipeline(pipelineName)
			Expect(err).NotTo(HaveOccurred())

			build, err = client.CreateJobBuild(pipelineName, "some-job")
			Expect(err).NotTo(HaveOccurred())

			loadingTimeout = 10 * time.Second
		})

		It("can view the resource", func() {
			pipelineURL := atcRoute(fmt.Sprintf("/pipelines/%s", pipelineName))
			// pipeline url -> job detail w/build info
			Expect(page.Navigate(pipelineURL)).To(Succeed())

			Eventually(page.FindByLink("some-job"), loadingTimeout).Should(BeFound())
			Expect(page.FindByLink("some-job").Click()).To(Succeed())

			// job detail w/build info -> job detail
			Eventually(page, loadingTimeout).Should(HaveURL(atcRoute(fmt.Sprintf(
				"/pipelines/%s/jobs/some-job/builds/%s",
				pipelineName,
				build.Name,
			))))

			Eventually(page.Find("h1"), loadingTimeout).Should(HaveText(fmt.Sprintf("some-job #%s", build.Name)))
			Expect(page.Find("h1 a").Click()).To(Succeed())
			Eventually(page, loadingTimeout).Should(HaveURL(atcRoute(fmt.Sprintf("/pipelines/%s/jobs/some-job", pipelineName))))

			// job-detail pausing
			Eventually(page.Find("#job-state.btn-pause"), loadingTimeout).Should(BeFound())
			Expect(page.Find("#job-state.btn-pause").Click()).To(Succeed())
			Eventually(page.Find("#job-state.btn-pause.enabled"), loadingTimeout).Should(BeFound())
			Eventually(page.Find("#job-state.btn-pause.disabled"), loadingTimeout).ShouldNot(BeFound())

			page.Refresh()

			Eventually(page.Find("#job-state.btn-pause.enabled"), loadingTimeout).Should(BeFound())
			Eventually(page.Find("#job-state.btn-pause.disabled"), loadingTimeout).ShouldNot(BeFound())

			Expect(page.Navigate(pipelineURL)).To(Succeed())
			Eventually(page.Find(".job.paused"), loadingTimeout).Should(BeFound())

			// job-detail unpausing
			Expect(page.Navigate(atcRoute(fmt.Sprintf("/pipelines/%s/jobs/some-job", pipelineName)))).To(Succeed())
			Eventually(page.Find("#job-state.btn-pause"), loadingTimeout).Should(BeFound())
			Expect(page.Find("#job-state.btn-pause").Click()).To(Succeed())
			Eventually(page.Find("#job-state.btn-pause.disabled"), loadingTimeout).Should(BeFound())
			Eventually(page.Find("#job-state.btn-pause.enabled"), loadingTimeout).ShouldNot(BeFound())
		})

		Describe("paused pipeline", func() {
			BeforeEach(func() {
				_, err := client.PausePipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
			})

			It("displays a blue header", func() {
				// pipeline URL -> job detail w/build info
				jobURL := atcRoute(fmt.Sprintf("/pipelines/%s/jobs/some-job", pipelineName))
				Expect(page.Navigate(jobURL)).To(Succeed())
				Eventually(page, loadingTimeout).Should(HaveURL(jobURL))

				// top bar should show the pipeline is paused
				Eventually(page.Find(".js-groups.paused")).Should(BeFound())
			})
		})
	})
})