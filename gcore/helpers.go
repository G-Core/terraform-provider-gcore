package gcore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func conditionField(fieldName, desc string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false, Description: "Whether or not to apply a boolean NOT operation to the rule's condition"},
				fieldName:  {Type: schema.TypeString, Required: true, Description: desc},
			},
		},
	}
}

func conditionMatch(field, desc string, match_types []string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation":   {Type: schema.TypeBool, Optional: true, Default: false, Description: "Whether or not to apply a boolean NOT operation to the rule's condition"},
				field:        {Type: schema.TypeString, Required: true, Description: desc},
				"match_type": {Type: schema.TypeString, Required: true, Description: "The type of matching condition", ValidateFunc: validation.StringInSlice(match_types, false)},
			},
		},
	}
}

func conditionHeader(header, value string, match_types []string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation":   {Type: schema.TypeBool, Optional: true, Default: false, Description: "Whether or not to apply a boolean NOT operation to the rule's condition"},
				header:       {Type: schema.TypeString, Required: true, Description: "The request header name"},
				value:        {Type: schema.TypeString, Required: true, Description: "The request header value"},
				"match_type": {Type: schema.TypeString, Required: true, MaxItems: 1, Description: "The type of matching condition", ValidateFunc: validation.StringInSlice(match_types, false)},
			},
		},
	}
}

func conditionExists(field, desc string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false, Description: "Whether or not to apply a boolean NOT operation to the rule's condition"},
				field:      {Type: schema.TypeString, Required: true, Description: desc},
			},
		},
	}
}

func conditionExistsWithValidation(field, desc string, validation schema.SchemaValidateFunc) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false, Description: "Whether or not to apply a boolean NOT operation to the rule's condition"},
				field:      {Type: schema.TypeString, Required: true, Description: desc, ValidateFunc: validation},
			},
		},
	}
}

func conditionList(field, desc string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"negation": {Type: schema.TypeBool, Optional: true, Default: false, Description: "Whether or not to apply a boolean NOT operation to the rule's condition"},
				field: {
					Type:        schema.TypeList,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: desc,
					Required:    true,
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

func listStringWithValidation(desc string, validation schema.SchemaValidateFunc) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: desc,
		ValidateFunc: validation,
	}
}
