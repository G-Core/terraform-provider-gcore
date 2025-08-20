// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_billing_reservation_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_billing_reservation"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudBillingReservationsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_billing_reservation.CloudBillingReservationsDataSourceModel)(nil)
	schema := cloud_billing_reservation.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
