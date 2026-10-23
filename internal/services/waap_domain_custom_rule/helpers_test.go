package waap_domain_custom_rule_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const resourceName = "gcore_waap_domain_custom_rule.test"

func testAccCheckWaapDomainCustomRuleDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_waap_domain_custom_rule" {
			continue
		}
		ruleID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse rule ID: %w", err)
		}
		domainID, err := strconv.ParseInt(rs.Primary.Attributes["domain_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse domain_id: %w", err)
		}
		_, err = client.Waap.Domains.CustomRules.Get(context.Background(), ruleID, waap.DomainCustomRuleGetParams{
			DomainID: domainID,
		})
		if err == nil {
			return fmt.Errorf("WAAP custom rule %d still exists", ruleID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking WAAP custom rule deletion: %w", err)
		}
	}
	return nil
}

// testAccWaapDomainCustomRuleConfig renders a rule with the given body: the
// action, the conditions and any optional attributes.
func testAccWaapDomainCustomRuleConfig(domainID int64, name string, enabled bool, body string) string {
	return fmt.Sprintf(`
resource "gcore_waap_domain_custom_rule" "test" {
  domain_id = %d
  name      = %q
  enabled   = %t
%s
}
`, domainID, name, enabled, body)
}
