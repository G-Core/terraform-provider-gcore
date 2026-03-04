#!/usr/bin/env bash
# Validates a Terraform example file by running terraform validate and optionally terraform plan.
# Usage: ./scripts/validate-example.sh <path-to-tf-file> [--plan]
#
# Outputs JSON: {"file": "...", "validate_valid": true/false, "plan_result": "...", "errors": [...]}

set -euo pipefail

TF_FILE="${1:?Usage: validate-example.sh <path-to-tf-file> [--plan]}"
DO_PLAN="${2:-}"

if [[ ! -f "$TF_FILE" ]]; then
  echo "{\"file\": \"$TF_FILE\", \"validate_valid\": false, \"plan_result\": \"skipped\", \"errors\": [\"file not found\"]}"
  exit 1
fi

# Resolve to absolute path
TF_FILE="$(cd "$(dirname "$TF_FILE")" && pwd)/$(basename "$TF_FILE")"
EXAMPLE_DIR="$(dirname "$TF_FILE")"

# Create temp working directory
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Write provider configuration
cat > "$TMPDIR/provider.tf" <<'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}
EOF

# Copy the target .tf file
cp "$TF_FILE" "$TMPDIR/"

# Copy any prereq files from the same directory
for f in "$EXAMPLE_DIR"/*.tf; do
  base="$(basename "$f")"
  if [[ "$base" != "$(basename "$TF_FILE")" && "$base" != "provider.tf" ]]; then
    cp "$f" "$TMPDIR/"
  fi
done

# Initialize terraform (suppress output)
cd "$TMPDIR"
terraform init -input=false -no-color >/dev/null 2>&1 || true

# Run terraform validate
VALIDATE_OUTPUT=$(terraform validate -json -no-color 2>/dev/null || true)
VALIDATE_VALID=$(echo "$VALIDATE_OUTPUT" | python3 -c "import sys,json; d=json.load(sys.stdin); print('true' if d.get('valid',False) else 'false')" 2>/dev/null || echo "false")

# Collect validation errors
VALIDATE_ERRORS=$(echo "$VALIDATE_OUTPUT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
errs = []
for diag in d.get('diagnostics', []):
    if diag.get('severity') == 'error':
        errs.append(diag.get('summary', '') + ': ' + diag.get('detail', ''))
print(json.dumps(errs))
" 2>/dev/null || echo "[]")

# Run terraform plan if requested
PLAN_RESULT="skipped"
PLAN_ERRORS="[]"
if [[ "$DO_PLAN" == "--plan" && -n "${GCORE_API_KEY:-}" ]]; then
  PLAN_OUTPUT=$(terraform plan -json -input=false -no-color 2>&1 || true)
  # Check for errors in plan output
  PLAN_HAS_ERROR=$(echo "$PLAN_OUTPUT" | python3 -c "
import sys, json
has_err = False
errs = []
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        d = json.loads(line)
        if d.get('@level') == 'error' or d.get('type') == 'diagnostic' and d.get('diagnostic',{}).get('severity') == 'error':
            has_err = True
            msg = d.get('diagnostic',{}).get('summary','') or d.get('@message','')
            errs.append(msg)
    except:
        pass
print(json.dumps({'has_error': has_err, 'errors': errs}))
" 2>/dev/null || echo '{"has_error": true, "errors": ["plan parse failed"]}')

  PLAN_HAS_ERR=$(echo "$PLAN_HAS_ERROR" | python3 -c "import sys,json; print(json.load(sys.stdin)['has_error'])" 2>/dev/null || echo "true")
  PLAN_ERRORS=$(echo "$PLAN_HAS_ERROR" | python3 -c "import sys,json; print(json.dumps(json.load(sys.stdin)['errors']))" 2>/dev/null || echo "[]")

  if [[ "$PLAN_HAS_ERR" == "True" || "$PLAN_HAS_ERR" == "true" ]]; then
    PLAN_RESULT="error"
  else
    PLAN_RESULT="success"
  fi
fi

# Combine all errors
ALL_ERRORS=$(python3 -c "
import json, sys
v = json.loads('$VALIDATE_ERRORS')
p = json.loads('$PLAN_ERRORS')
print(json.dumps(v + p))
" 2>/dev/null || echo "[]")

# Output result as JSON
echo "{\"file\": \"$TF_FILE\", \"validate_valid\": $VALIDATE_VALID, \"plan_result\": \"$PLAN_RESULT\", \"errors\": $ALL_ERRORS}"
