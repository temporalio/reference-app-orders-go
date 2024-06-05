package billing_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temporalio/reference-app-orders-go/app/billing"
)

func TestChargeWorkflowID(t *testing.T) {
	wfid := billing.ChargeWorkflowID(billing.ChargeInput{
		IdempotencyKey: "test",
	})

	assert.Equal(t, "Charge:test", wfid)

	wfid = billing.ChargeWorkflowID(billing.ChargeInput{
		IdempotencyKey: "",
	})

	assert.Regexp(t, regexp.MustCompile("Charge:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+"), wfid)
}
