package gcore

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iam "github.com/G-Core/gcore-iam-sdk-go"
)

func resourceIamApiToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamApiTokenCreate,
		ReadContext:   resourceIamApiTokenRead,
		UpdateContext: resourceIamApiTokenUpdate,
		DeleteContext: resourceIamApiTokenDelete,
		Description:   "Represent IAM ApiToken for authentication and management capabilities",

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Account ID.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API token value",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API token name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API token description.",
			},
			"exp_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Date when the API token becomes expired (ISO 8086/RFC 3339 format), UTC. If null, then the API token will never expire.",
			},
			"client_user": {
				Type:        schema.TypeList,
				Required:    true,
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
						"deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Deletion flag. If true, then the API token was deleted.",
						},
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User's ID who issued the API token.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User's name who issued the API token.",
						},
						"user_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User's email who issued the API token.",
						},
						"client_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account's ID.",
						},
					},
				},
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Deletion flag. If true, then the API token was deleted.",
			},
			"expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Expiration flag. If true, then the API token has expired. When an API token expires it will be automatically deleted.",
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
	if v, ok := d.GetOk("exp_date"); ok {
		expDateStr := v.(string)
		if expDateStr != "" {
			expDatePtr = &expDateStr
		} else {
			expDatePtr = nil
		}
	} else {
		expDatePtr = nil
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

	if getResp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to get Api Token. Status code: %d with error: %s", getResp.StatusCode(), getResp.Body)
	}

	if getResp.JSON200 == nil {
		return diag.Errorf("Failed to get token")
	}

	token := getResp.JSON200

	if token.Id != nil {
		d.Set("token_id", *token.Id)
	}

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

	if token.Deleted != nil {
		d.Set("deleted", *token.Deleted)
	}

	if token.Expired != nil {
		d.Set("expired", *token.Expired)
	}

	clientUserData := map[string]interface{}{}

	if token.ClientUser.Deleted != nil {
		clientUserData["deleted"] = *token.ClientUser.Deleted
	}

	if token.ClientUser.UserId != nil {
		clientUserData["user_id"] = *token.ClientUser.UserId
	}

	if token.ClientUser.UserName != nil {
		clientUserData["user_name"] = *token.ClientUser.UserName
	}

	if token.ClientUser.UserEmail != nil {
		clientUserData["user_email"] = *token.ClientUser.UserEmail
	}

	if token.ClientUser.ClientId != nil {
		clientUserData["client_id"] = *token.ClientUser.ClientId
	}

	if token.ClientUser.Role != nil {
		roleData := map[string]interface{}{}

		if token.ClientUser.Role.Id != nil {
			roleData["id"] = *token.ClientUser.Role.Id
		}

		if token.ClientUser.Role.Name != nil {
			roleData["name"] = string(*token.ClientUser.Role.Name)
		}

		clientUserData["role"] = []interface{}{roleData}
	}

	d.Set("client_user", []interface{}{clientUserData})

	return nil
}

func resourceIamApiTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIamApiTokenRead(ctx, d, m)
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
