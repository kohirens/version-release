package circleci

import (
	"github.com/kohirens/stdlib/log"
)

func (c *Client) RunWorkflow(branch, nameWorkflow string) error {
	// Build pipeline parameters to trigger the tag-and-release workflow.
	pp, e1 := PipelineParameters(branch, nameWorkflow)
	if e1 != nil {
		return e1
	}

	log.Logf(stdout.TriggerWorkflow, nameWorkflow)

	return c.TriggerWorkflow(pp)
}
