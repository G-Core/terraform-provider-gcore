package gcore

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateCDNOptions(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
	optionsConfig := diff.Get("options").([]interface{})

	if len(optionsConfig) == 0 {
		return nil
	}

	optionsMap := optionsConfig[0].(map[string]interface{})

	// validate FastEdge
	fastedgeOption, ok := optionsMap["fastedge"].([]interface{})

	if ok && len(fastedgeOption) > 0 {
		fastedgeMap := fastedgeOption[0].(map[string]interface{})
		triggers := []string{"on_request_headers", "on_request_body", "on_response_headers", "on_response_body"}

		for _, trigger := range triggers {
			if triggerConfig, exists := fastedgeMap[trigger].([]interface{}); exists && len(triggerConfig) > 0 {
				return nil
			}
		}

		return fmt.Errorf("at least one of 'on_request_headers', 'on_request_body', 'on_response_headers', or 'on_response_body' must be specified for the FastEdge option")
	}

	return nil
}
