package gcore

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iam "github.com/G-Core/gcore-iam-sdk-go"
)

func resourceIamAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamAccountCreate,
		ReadContext:   resourceIamAccountRead,
		UpdateContext: resourceIamAccountUpdate,
		DeleteContext: resourceIamAccountDelete,
		Description:   "Represent IAM Account with authentication and management capabilities",

		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Token to validate user.",
			},
			// Account creation fields
			"company": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The company name.",
			},
			// Account details from API response
			"account_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Account details retrieved from the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The account ID.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account email.",
						},
						"phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone of a user who registered the requested account.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a user who registered the requested account.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the account (new, trial, trialend, active, integration, paused, preparation, ready).",
						},
						"company_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The company name from the API response.",
						},
						"website": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The company website.",
						},
						"current_user": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the current user.",
						},
						"capabilities": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of services available for the account (CDN, STORAGE, STREAMING, DNS, DDOS, CLOUD).",
						},
						"deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The field shows the status of the account: true – the account has been deleted, false – the account is not deleted.",
						},
						"bill_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System field. Billing type of the account.",
						},
						"custom_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account custom ID.",
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System field. The company country (ISO 3166-1 alpha-2 format).",
						},
						"is_test": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "System field: true — a test account; false — a production account.",
						},
						"has_active_admin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "System field indicating if the account has an active admin.",
						},
						"entry_base_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System field. Control panel domain.",
						},
						"signup_process": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System field. Type of the account registration process (sign_up_full, sign_up_simple).",
						},
						"users": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of account users.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "User's ID.",
									},
									"email": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User's email address.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User's name.",
									},
									"phone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User's phone.",
									},
									"company": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User's company.",
									},
									"activated": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Email confirmation status.",
									},
									"deleted": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Deletion flag.",
									},
									"sso_auth": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "SSO authentication flag.",
									},
									"is_active": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "User active status.",
									},
									"groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "User's groups in the current account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Group's ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group's name.",
												},
											},
										},
									},
								},
							},
						},
						"service_statuses": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "An object containing information about all services available for the requested account.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"login_activity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Login activity log for the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The login record ID.",
						},
						"login_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Login timestamp in ISO 8086/RFC 3339 format.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP Address of the login attempt.",
						},
						"user_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User-Agent string of the login attempt.",
						},
						"is_successful": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the login was successful or not.",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of authentication used.",
						},
						"user": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User ID associated with the login attempt.",
						},
					},
				},
			},
		},
	}
}

func resourceIamAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	apiToken := d.Get("api_token").(string)

	authFunc := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "APIKey "+apiToken)
		return nil
	}

	if company, ok := d.GetOk("company"); ok {
		createReq := iam.PostClientsCreateJSONRequestBody{
			Company: company.(string),
		}

		createResp, err := client.PostClientsCreateWithResponse(ctx, createReq, authFunc)
		if err != nil {
			return diag.Errorf("Failed to create IAM account: %v", err)
		}

		if createResp.StatusCode() != http.StatusCreated {
			return diag.Errorf("Failed to create IAM account. Status code: %d with error: %s", createResp.StatusCode(), createResp.Body)
		}

		if createResp.JSON201 == nil || createResp.JSON201.Id == nil {
			return diag.Errorf("Failed to get account ID from create response")
		}

		d.SetId(strconv.Itoa(*createResp.JSON201.Id))
		return resourceIamAccountRead(ctx, d, m)
	}
	return diag.Errorf("Company name is required to create an IAM account")
}

func resourceIamAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	apiToken := d.Get("api_token").(string)

	authFunc := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "APIKey "+apiToken)
		return nil
	}

	resp, err := client.GetClientsMeWithResponse(ctx, authFunc)
	if err != nil {
		return diag.Errorf("Failed to get IAM account details: %v", err)
	}

	if resp.StatusCode() != 200 {
		return diag.Errorf("Failed to get IAM account details. Status code: %d", resp.StatusCode())
	}

	if resp.JSON200 == nil {
		return diag.Errorf("Empty response from IAM account details API")
	}

	account := resp.JSON200

	if account.Id != nil {
		d.SetId(strconv.Itoa(*account.Id))
	}

	accountDetails := map[string]interface{}{}
	if account.Id != nil {
		accountDetails["account_id"] = *account.Id
	}
	if account.Email != nil {
		accountDetails["email"] = string(*account.Email)
	}
	if account.Phone != nil {
		accountDetails["phone"] = *account.Phone
	}
	if account.Name != nil {
		accountDetails["name"] = *account.Name
	}
	if account.Status != nil {
		accountDetails["status"] = string(*account.Status)
	}
	if account.CompanyName != nil {
		accountDetails["company_name"] = *account.CompanyName
	}
	if account.Website != nil {
		accountDetails["website"] = *account.Website
	}
	if account.CurrentUser != nil {
		accountDetails["current_user"] = *account.CurrentUser
	}
	if account.Deleted != nil {
		accountDetails["deleted"] = *account.Deleted
	}
	if account.BillType != nil {
		accountDetails["bill_type"] = *account.BillType
	}
	if account.CustomId != nil {
		accountDetails["custom_id"] = *account.CustomId
	}
	if account.CountryCode != nil {
		accountDetails["country_code"] = *account.CountryCode
	}
	if account.IsTest != nil {
		accountDetails["is_test"] = *account.IsTest
	}
	if account.HasActiveAdmin != nil {
		accountDetails["has_active_admin"] = *account.HasActiveAdmin
	}
	if account.EntryBaseDomain != nil {
		accountDetails["entry_base_domain"] = *account.EntryBaseDomain
	}
	if account.SignupProcess != nil {
		accountDetails["signup_process"] = string(*account.SignupProcess)
	}
	if account.Capabilities != nil {
		capabilities := make([]string, len(*account.Capabilities))
		for i, cap := range *account.Capabilities {
			capabilities[i] = string(cap)
		}
		accountDetails["capabilities"] = capabilities
	}
	if account.Users != nil {
		users := make([]map[string]interface{}, len(*account.Users))
		for i, user := range *account.Users {
			userMap := map[string]interface{}{}
			if user.Id != nil {
				userMap["id"] = *user.Id
			}
			if user.Email != nil {
				userMap["email"] = string(*user.Email)
			}
			if user.Name != nil {
				userMap["name"] = *user.Name
			}
			if user.Phone != nil {
				userMap["phone"] = *user.Phone
			}
			if user.Company != nil {
				userMap["company"] = *user.Company
			}
			if user.Activated != nil {
				userMap["activated"] = *user.Activated
			}
			if user.Deleted != nil {
				userMap["deleted"] = *user.Deleted
			}
			if user.SsoAuth != nil {
				userMap["sso_auth"] = *user.SsoAuth
			}
			if user.Groups != nil {
				groups := make([]map[string]interface{}, len(*user.Groups))
				for j, group := range *user.Groups {
					groupMap := map[string]interface{}{}
					if group.Id != nil {
						groupMap["id"] = *group.Id
					}
					if group.Name != nil {
						groupMap["name"] = string(*group.Name)
					}
					groups[j] = groupMap
				}
				userMap["groups"] = groups
			}
			users[i] = userMap
		}
		accountDetails["users"] = users
	}
	if account.ServiceStatuses != nil {
		serviceStatuses := make(map[string]string)
		if account.ServiceStatuses.CDN != nil && account.ServiceStatuses.CDN.Status != nil {
			serviceStatuses["CDN"] = string(*account.ServiceStatuses.CDN.Status)
		}
		if account.ServiceStatuses.CLOUD != nil && account.ServiceStatuses.CLOUD.Status != nil {
			serviceStatuses["CLOUD"] = string(*account.ServiceStatuses.CLOUD.Status)
		}
		if account.ServiceStatuses.DNS != nil && account.ServiceStatuses.DNS.Status != nil {
			serviceStatuses["DNS"] = string(*account.ServiceStatuses.DNS.Status)
		}
		if account.ServiceStatuses.DDOS != nil && account.ServiceStatuses.DDOS.Status != nil {
			serviceStatuses["DDOS"] = string(*account.ServiceStatuses.DDOS.Status)
		}
		if account.ServiceStatuses.STORAGE != nil && account.ServiceStatuses.STORAGE.Status != nil {
			serviceStatuses["STORAGE"] = string(*account.ServiceStatuses.STORAGE.Status)
		}
		if account.ServiceStatuses.STREAMING != nil && account.ServiceStatuses.STREAMING.Status != nil {
			serviceStatuses["STREAMING"] = string(*account.ServiceStatuses.STREAMING.Status)
		}
		accountDetails["service_statuses"] = serviceStatuses
	}
	d.Set("account_details", []map[string]interface{}{accountDetails})

	httpResp, err := client.GetActivityLogLogins(ctx, nil, authFunc)
	if err != nil {
		return diag.Errorf("Failed to get login activity: %v", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 200 {
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return diag.Errorf("failed to read login activity response body: %v", err)
		}

		trimmedBody := bytes.TrimSpace(body)
		var resultsArray []interface{}

		if len(trimmedBody) > 0 {
			if trimmedBody[0] == '{' {
				var loginActivityResponse map[string]interface{}
				if err := json.Unmarshal(body, &loginActivityResponse); err != nil {
					return diag.Errorf("failed to unmarshal login activity response object: %s", err)
				}

				if resultsData, ok := loginActivityResponse["results"]; ok {
					if data, ok := resultsData.([]interface{}); ok {
						resultsArray = data
					} else {
						return diag.Errorf("login activity 'results' data is not an array, but is %T", resultsData)
					}
				}
			} else if trimmedBody[0] == '[' {
				if err := json.Unmarshal(body, &resultsArray); err != nil {
					return diag.Errorf("failed to unmarshal login activity response array: %s", err)
				}
			} else {
				return diag.Errorf("unexpected format for login activity response body, starts with %c", trimmedBody[0])
			}
		}

		logins := make([]map[string]interface{}, 0, len(resultsArray))
		for _, item := range resultsArray {
			loginItem, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			login := map[string]interface{}{}
			if id, ok := loginItem["id"].(float64); ok {
				login["id"] = int(id)
			}
			if loginAt, ok := loginItem["login_at"].(string); ok {
				login["login_at"] = loginAt
			}
			if ipAddress, ok := loginItem["ip_address"].(string); ok {
				login["ip_address"] = ipAddress
			}
			if userAgent, ok := loginItem["user_agent"].(string); ok {
				login["user_agent"] = userAgent
			}
			if isSuccessful, ok := loginItem["is_successful"].(bool); ok {
				login["is_successful"] = isSuccessful
			}
			if authType, ok := loginItem["auth_type"].(string); ok {
				login["auth_type"] = authType
			}
			if user, ok := loginItem["user"].(float64); ok {
				login["user"] = int(user)
			}
			logins = append(logins, login)
		}
		if err := d.Set("login_activity", logins); err != nil {
			return diag.Errorf("Failed to set login activity: %v", err)
		}
	} else if httpResp.StatusCode != 404 {
		body, _ := io.ReadAll(httpResp.Body)
		return diag.Errorf("Failed to get login activity. Status code: %d, Body: %s", httpResp.StatusCode, string(body))
	}

	return nil
}

func resourceIamAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIamAccountRead(ctx, d, m)
}

func resourceIamAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
