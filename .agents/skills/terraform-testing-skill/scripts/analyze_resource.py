#!/usr/bin/env python3
"""
analyze_resource.py - Analyze a resource to determine what needs testing
This helps Claude identify special operations, computed fields, and edge cases
"""

import os
import re
import subprocess
from pathlib import Path

def analyze_old_provider(resource_name):
    """Analyze old provider implementation"""
    
    old_file = f"old_terraform_provider/gcore/resource_gcore_{resource_name}.go"
    
    if not Path(old_file).exists():
        print(f"Old provider file not found: {old_file}")
        return {}
    
    with open(old_file, 'r') as f:
        content = f.read()
    
    analysis = {
        "uses_tasks": "WaitTaskAndReturn" in content,
        "has_unset": "Unset(" in content,
        "special_operations": [],
        "complex_updates": [],
        "computed_fields": []
    }
    
    # Find special operations
    special_ops = re.findall(r'func.*?(attach|detach|add|remove).*?\(', content, re.IGNORECASE)
    analysis["special_operations"] = list(set(special_ops))
    
    # Find complex update logic
    if "if d.HasChange" in content:
        changes = re.findall(r'if d\.HasChange\("(\w+)"\)', content)
        analysis["complex_updates"] = changes
    
    # Find fields that might be computed
    computed_patterns = re.findall(r'd\.Set\("(\w+)",.*?(?:computed|default)', content, re.IGNORECASE)
    analysis["computed_fields"] = computed_patterns
    
    return analysis

def analyze_new_provider(resource_name):
    """Analyze new provider implementation"""
    
    resource_dir = f"internal/services/cloud_{resource_name}"
    
    if not Path(resource_dir).exists():
        print(f"New provider directory not found: {resource_dir}")
        return {}
    
    analysis = {
        "uses_polling": False,
        "computed_optional_fields": [],
        "required_fields": [],
        "updatable_fields": []
    }
    
    # Check resource.go for polling
    resource_file = Path(resource_dir) / "resource.go"
    if resource_file.exists():
        with open(resource_file, 'r') as f:
            content = f.read()
            analysis["uses_polling"] = "AndPoll" in content
    
    # Check model.go for field types
    model_file = Path(resource_dir) / "model.go"
    if model_file.exists():
        with open(model_file, 'r') as f:
            content = f.read()
            computed_optional = re.findall(r'(\w+)\s+types\.\w+\s+`.*?json:".*?,computed_optional"', content)
            analysis["computed_optional_fields"] = computed_optional
    
    # Check schema.go for field definitions
    schema_file = Path(resource_dir) / "schema.go"
    if schema_file.exists():
        with open(schema_file, 'r') as f:
            content = f.read()
            required = re.findall(r'"(\w+)":\s+schema\.\w+Attribute\{[^}]*Required:\s+true', content)
            analysis["required_fields"] = required
    
    return analysis

def suggest_test_cases(old_analysis, new_analysis):
    """Suggest test cases based on analysis"""
    
    test_cases = []
    
    # Drift testing for computed_optional fields
    if new_analysis.get("computed_optional_fields"):
        test_cases.append({
            "name": "Drift Test - Computed Fields",
            "description": "Test that computed_optional fields don't cause drift",
            "fields": new_analysis["computed_optional_fields"],
            "approach": "Apply with minimal config, then plan to check for drift"
        })
    
    # Async operation testing
    if old_analysis.get("uses_tasks") or new_analysis.get("uses_polling"):
        test_cases.append({
            "name": "Async Operations",
            "description": "Verify polling works correctly",
            "approach": "Monitor task completion in logs"
        })
    
    # Special operations
    if old_analysis.get("special_operations"):
        test_cases.append({
            "name": "Special Operations",
            "description": f"Test special endpoints: {', '.join(old_analysis['special_operations'])}",
            "approach": "Use mitmproxy to verify correct API calls"
        })
    
    # Complex updates
    if old_analysis.get("complex_updates"):
        for field in old_analysis["complex_updates"]:
            test_cases.append({
                "name": f"Update {field}",
                "description": f"Test updating {field} field",
                "approach": "Update field and verify PATCH operation, not recreate"
            })
    
    # Unset operations
    if old_analysis.get("has_unset"):
        test_cases.append({
            "name": "Field Unsetting",
            "description": "Test clearing optional fields",
            "approach": "Set field then clear it, verify Unset API call"
        })
    
    return test_cases

def generate_test_report(resource_name):
    """Generate analysis report for the resource"""
    
    print(f"Analyzing {resource_name}...")
    print("=" * 50)
    
    old = analyze_old_provider(resource_name)
    new = analyze_new_provider(resource_name)
    test_cases = suggest_test_cases(old, new)
    
    print("\n📋 Old Provider Analysis:")
    print(f"  - Uses async tasks: {old.get('uses_tasks', False)}")
    print(f"  - Has Unset operation: {old.get('has_unset', False)}")
    print(f"  - Special operations: {', '.join(old.get('special_operations', [])) or 'None'}")
    print(f"  - Complex update fields: {', '.join(old.get('complex_updates', [])) or 'None'}")
    
    print("\n📋 New Provider Analysis:")
    print(f"  - Uses polling: {new.get('uses_polling', False)}")
    print(f"  - Computed optional fields: {', '.join(new.get('computed_optional_fields', [])) or 'None'}")
    print(f"  - Required fields: {', '.join(new.get('required_fields', [])) or 'None'}")
    
    print("\n🧪 Suggested Test Cases:")
    for i, test in enumerate(test_cases, 1):
        print(f"\n  {i}. {test['name']}")
        print(f"     Description: {test['description']}")
        print(f"     Approach: {test['approach']}")
        if 'fields' in test:
            print(f"     Fields: {', '.join(test['fields'])}")
    
    print("\n" + "=" * 50)
    print("Use this analysis to create targeted Terraform test configurations")

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Usage: analyze_resource.py <resource_name>")
        print("Example: analyze_resource.py router")
        sys.exit(1)
    
    resource = sys.argv[1]
    generate_test_report(resource)