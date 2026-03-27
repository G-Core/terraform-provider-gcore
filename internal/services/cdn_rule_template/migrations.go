// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_rule_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithUpgradeState = (*CDNRuleTemplateResource)(nil)

func (r *CDNRuleTemplateResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
