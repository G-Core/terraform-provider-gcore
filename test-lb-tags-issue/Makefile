.PHONY: help init plan apply destroy clean test-all

help:
	@echo "Available targets:"
	@echo "  init           - Initialize Terraform"
	@echo "  plan           - Run terraform plan"
	@echo "  apply          - Run terraform apply"
	@echo "  destroy        - Destroy all resources"
	@echo "  clean          - Clean up test artifacts"
	@echo ""
	@echo "Test scenarios:"
	@echo "  test-scenario-1 - Create LB with tags from start"
	@echo "  test-scenario-2 - Create LB, then add tags later"
	@echo "  test-scenario-3 - Exact Jira reproduction"
	@echo "  test-scenario-4 - Modify existing tags"
	@echo "  test-all       - Run all scenarios"

init:
	@echo "==> Initializing Terraform (skipped with dev overrides)"
	@echo "Provider dev overrides are active, no init needed"

plan:
	@echo "==> Running terraform plan..."
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform plan"

apply:
	@echo "==> Running terraform apply..."
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform apply -auto-approve"

destroy:
	@echo "==> Destroying resources..."
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform destroy -auto-approve"

clean:
	@echo "==> Cleaning up..."
	rm -f terraform.tfstate terraform.tfstate.backup
	rm -f *.log
	rm -f reproduce_tags_error.tf.backup

test-scenario-1:
	@echo "========================================"
	@echo "SCENARIO 1: Create LB with tags"
	@echo "========================================"
	@echo "This tests creating an LB with tags from the start"
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform apply -target=gcore_cloud_load_balancer.lb_with_tags -auto-approve"

test-scenario-2:
	@echo "========================================"
	@echo "SCENARIO 2: Add tags to existing LB"
	@echo "========================================"
	@echo "Step 1: Creating LB without tags..."
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform apply -target=gcore_cloud_load_balancer.lb_add_tags_later -auto-approve"
	@echo ""
	@echo "Step 2: Uncomment tags and apply again..."
	@echo "Manual step required: Edit reproduce_tags_error.tf and uncomment tags for lb_add_tags_later"
	@echo "Then run: make apply"

test-scenario-3:
	@echo "========================================"
	@echo "SCENARIO 3: Exact Jira reproduction"
	@echo "========================================"
	@echo "Testing with exact manifest from Jira issue"
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform apply -target=gcore_cloud_load_balancer.lb -auto-approve"

test-scenario-4:
	@echo "========================================"
	@echo "SCENARIO 4: Modify existing tags"
	@echo "========================================"
	@echo "Step 1: Creating LB with initial tag..."
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform apply -target=gcore_cloud_load_balancer.lb_modify_tags -auto-approve"
	@echo ""
	@echo "Step 2: Modify tags and apply again..."
	@echo "Manual step required: Edit reproduce_tags_error.tf and change tags for lb_modify_tags"
	@echo "Then run: make apply"

test-all: test-scenario-1 test-scenario-3
	@echo "========================================"
	@echo "All automated tests complete"
	@echo "========================================"
	@echo "Check outputs for any inconsistency errors"
	@bash -c "source ../.env && export TF_CLI_CONFIG_FILE=../.terraformrc && terraform output"
