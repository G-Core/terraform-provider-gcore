package gcore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func conditionField(fieldName string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false},
				fieldName:  {Type: schema.TypeString, Required: true},
			},
		},
	}
}

func conditionRangeFields(lower, upper string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false},
				lower:      {Type: schema.TypeString, Required: true},
				upper:      {Type: schema.TypeString, Required: true},
			},
		},
	}
}

func conditionMatch(field string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation":   {Type: schema.TypeBool, Optional: true, Default: false},
				field:        {Type: schema.TypeString, Required: true},
				"match_type": {Type: schema.TypeString, Required: true},
			},
		},
	}
}

func conditionHeader(header, value string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation":   {Type: schema.TypeBool, Optional: true, Default: false},
				header:       {Type: schema.TypeString, Required: true},
				value:        {Type: schema.TypeString, Required: true},
				"match_type": {Type: schema.TypeString, Required: true},
			},
		},
	}
}

func conditionExists(field string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false},
				field:      {Type: schema.TypeString, Required: true},
			},
		},
	}
}

func conditionList(field string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false},
				field: {
					Type:     schema.TypeList,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Required: true,
				},
			},
		},
	}
}

func stringField(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Required: true, Description: desc}
}

func intField(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeInt, Required: true, Description: desc}
}

func listString(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: desc,
	}
}
