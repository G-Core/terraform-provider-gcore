package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iam "github.com/G-Core/gcore-iam-sdk-go"
)

func resourceIamApiToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamApiTokenCreate,
		ReadContext:   resourceIamApiTokenRead,
		DeleteContext: resourceIamApiTokenDelete,
		Description:   "Represent IAM ApiToken for authentication and management capabilities",

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Account ID.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API token value",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "API token name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "API token description.",
			},
			"exp_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Date when the API token becomes expired (ISO 8086/RFC 3339 format), UTC. If null, the token will expire in one week.",
			},
			"client_user": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "API token role.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Group's ID",
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Group's name.",
									},
								},
							},
						},
					},
				},
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the API token was issued (ISO 8086/RFC 3339 format), UTC.",
			},
			"last_usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the API token was last used (ISO 8086/RFC 3339 format), UTC.",
			},
		},
	}
}

func resourceIamApiTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	clientId := d.Get("client_id").(int)
	name := d.Get("name").(string)

	var expDatePtr *string
	if v, ok := d.GetOk("exp_date"); ok && v.(string) != "" {
		expDateStr := v.(string)
		expDatePtr = &expDateStr
	} else {
		expTime := time.Now().UTC().Add(7 * 24 * time.Hour)
		expDateStr := expTime.Format(time.RFC3339)
		expDatePtr = &expDateStr
	}

	clientUser := d.Get("client_user").([]interface{})
	roleConfig := clientUser[0].(map[string]interface{})
	roleList := roleConfig["role"].([]interface{})
	roleMap := roleList[0].(map[string]interface{})

	roleId := roleMap["id"].(int)
	roleName := iam.PostClientsClientIdTokensJSONBodyClientUserRoleName(roleMap["name"].(string))

	description := d.Get("description").(string)

	createReq := iam.PostClientsClientIdTokensJSONRequestBody{
		Name:    name,
		ExpDate: expDatePtr,
		ClientUser: struct {
			Role *struct {
				Id   *int                                                     `json:"id,omitempty"`
				Name *iam.PostClientsClientIdTokensJSONBodyClientUserRoleName `json:"name,omitempty"`
			} `json:"role,omitempty"`
		}{
			Role: &struct {
				Id   *int                                                     `json:"id,omitempty"`
				Name *iam.PostClientsClientIdTokensJSONBodyClientUserRoleName `json:"name,omitempty"`
			}{
				Id:   &roleId,
				Name: &roleName,
			},
		},
		Description: &description,
	}

	createResp, err := client.PostClientsClientIdTokensWithResponse(ctx, clientId, createReq)
	if err != nil {
		return diag.Errorf("Failed to create Api Token: %v", err)
	}

	if createResp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to create Api Token. Status code: %d with error: %s", createResp.StatusCode(), createResp.Body)
	}

	if createResp.JSON200 == nil || createResp.JSON200.Token == nil {
		return diag.Errorf("Failed to get token from create response")
	}

	split := strings.SplitN(*createResp.JSON200.Token, "$", 2)

	tokenIdStr := split[0]
	tokenId, err := strconv.Atoi(tokenIdStr)
	if err != nil {
		log.Fatalf("Invalid token id: %v", err)
	}

	d.SetId(strconv.Itoa(tokenId))
	d.Set("token", *createResp.JSON200.Token)

	return resourceIamApiTokenRead(ctx, d, m)
}

func resourceIamApiTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	clientId := d.Get("client_id").(int)
	tokenId, _ := strconv.Atoi(d.Id())

	getResp, err := client.GetClientsClientIdTokensTokenIdWithResponse(ctx, clientId, tokenId)
	if err != nil {
		return diag.Errorf("Failed to get Api Token: %v", err)
	}

	if getResp.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Token (%d) was not found, removed from TF state", tokenId),
			},
		}
	}

	if getResp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to get Api Token. Status code: %d with error: %s", getResp.StatusCode(), getResp.Body)
	}

	if getResp.JSON200 == nil {
		return diag.Errorf("Failed to get token")
	}

	token := getResp.JSON200

	d.Set("name", token.Name)

	if token.Description != nil {
		d.Set("description", *token.Description)
	}

	if token.ExpDate != nil {
		d.Set("exp_date", *token.ExpDate)
	}

	if token.Created != nil {
		d.Set("created", *token.Created)
	}

	if token.LastUsage != nil {
		d.Set("last_usage", *token.LastUsage)
	}

	if token.ClientUser.Role != nil {
		roleData := map[string]interface{}{}

		if token.ClientUser.Role.Id != nil {
			roleData["id"] = *token.ClientUser.Role.Id
		}

		if token.ClientUser.Role.Name != nil {
			roleData["name"] = string(*token.ClientUser.Role.Name)
		}

		if err := d.Set("client_user", []interface{}{
			map[string]interface{}{
				"role": []interface{}{roleData},
			},
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceIamApiTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	clientId := d.Get("client_id").(int)
	tokenId, _ := strconv.Atoi(d.Id())

	resp, err := client.DeleteClientsClientIdTokensTokenIdWithResponse(ctx, clientId, tokenId)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to delete Api Token. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	return nil
}