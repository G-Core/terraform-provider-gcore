// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_insight_silence"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainInsightSilenceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_insight_silence.WaapDomainInsightSilenceModel)(nil)
	schema := waap_domain_insight_silence.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
