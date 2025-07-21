package gcore

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	sdk "github.com/G-Core/gcore-iam-sdk-go"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIamUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamUserCreate,
		ReadContext:   resourceIamUserRead,
		UpdateContext: resourceIamUserUpdate,
		DeleteContext: resourceIamUserDelete,
		Description:   "Manage IAM user with complete lifecycle support",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User's email address.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User's name.",
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "en",
				ValidateFunc: validation.StringInSlice([]string{"de", "en", "ru", "zh", "az"}, false),
				Description:  "User's language. Defines language of the control panel and email messages. Possible values: de, en, ru, zh, az.",
			},
			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User's phone.",
			},
			"company": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User's company.",
			},
			"client_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of account to invite user to.",
			},
			"user_role": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Group name.",
						},
					},
				},
				Description: "User's group in the account. IAM supports 5 groups: Users, Administrators, Engineers, Purge and Prefetch only (API), Purge and Prefetch only (API+Web).",
			},
			"auth_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"password", "sso", "github", "google-oauth2"}, false),
				},
				Description: "List of auth types available for the account. Possible values: password, sso, github, google-oauth2.",
			},
			// Read-only computed fields
			"user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User's ID.",
			},
			"current_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User's current email address.",
			},
			"reseller": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Services provider ID.",
			},
			"client": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "User's account ID.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Deletion flag. If true then user was deleted.",
			},
			"activated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Email confirmation: true – user confirmed the email; false – user did not confirm the email.",
			},
			"sso_auth": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "SSO authentication flag. If true then user can login via SAML SSO.",
			},
			"two_fa": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Two-step verification: true – user enabled two-step verification; false – user disabled two-step verification.",
			},
			"user_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User's type.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "User activity flag.",
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
							Description: "Group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group name.",
						},
					},
				},
			},
			"client_and_roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user's clients and roles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Client ID.",
						},
						"client_company_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User's role in the client.",
						},
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User ID.",
						},
						"user_roles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User roles in this client.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceIamUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	email := d.Get("email").(string)
	clientID := d.Get("client_id").(int)
	// Prepare user role
	userRoleConfigs := d.Get("user_role").([]interface{})
	userRoleConfig := userRoleConfigs[0].(map[string]interface{})

	// Prepare user role properly
	roleID := userRoleConfig["id"].(int)
	roleName := userRoleConfig["name"].(string)

	roleNameTyped := sdk.PostClientsInviteUserJSONBodyUserRoleName(roleName)

	// Prepare invite request using the correct SDK types
	inviteReq := sdk.PostClientsInviteUserJSONRequestBody{
		Email:    openapi_types.Email(email),
		ClientId: clientID,
		UserRole: struct {
			Id   *int                                           `json:"id,omitempty"`
			Name *sdk.PostClientsInviteUserJSONBodyUserRoleName `json:"name,omitempty"`
		}{
			Id:   &roleID,
			Name: &roleNameTyped,
		},
	}

	if name, ok := d.GetOk("name"); ok {
		nameStr := name.(string)
		inviteReq.Name = &nameStr
	}

	if lang, ok := d.GetOk("lang"); ok {
		langStr := lang.(string)
		langTyped := sdk.PostClientsInviteUserJSONBodyLang(langStr)
		inviteReq.Lang = &langTyped
	}

	// Invite user
	resp, err := client.PostClientsInviteUserWithResponse(ctx, inviteReq)
	if err != nil {
		return diag.Errorf("failed to invite user: %s", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("failed to invite user: status code %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil {
		return diag.Errorf("failed to invite user: empty response")
	}

	// Set the user ID as resource ID
	userID := resp.JSON200.UserId
	if userID == nil {
		return diag.Errorf("failed to invite user: no user ID returned")
	}
	d.SetId(fmt.Sprintf("%d", *userID))
	d.Set("user_id", *userID)

	// Update user with additional details if provided
	// Use forceUpdate=true for creation to check GetOk instead of HasChanges
	updateDiags := updateUserDetails(ctx, d, m, *userID, true)
	if updateDiags.HasError() {
		return updateDiags
	}

	return resourceIamUserRead(ctx, d, m)
}

func resourceIamUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	userID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert user ID: %s", err)
	}

	// Get user details
	resp, err := client.GetUsersUserIdWithResponse(ctx, userID)
	if err != nil {
		return diag.Errorf("failed to get user: %s", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("User (%d) was not found, removed from TF state", userID)},
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("failed to read user: status code %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil {
		return diag.Errorf("failed to read user: empty response")
	}

	user := resp.JSON200

	// Set user fields - always set to ensure proper state refresh
	var userIDValue int
	if user.Id != nil {
		userIDValue = *user.Id
	}
	d.Set("user_id", userIDValue)

	var currentEmailValue string
	if user.Email != nil {
		currentEmailValue = string(*user.Email)
	}
	d.Set("current_email", currentEmailValue)

	var resellerValue int
	if user.Reseller != nil {
		resellerValue = *user.Reseller
	}
	d.Set("reseller", resellerValue)

	var clientValue float64
	if user.Client != nil {
		clientValue = float64(*user.Client)
	}
	d.Set("client", clientValue)

	var deletedValue bool
	if user.Deleted != nil {
		deletedValue = *user.Deleted
	}
	d.Set("deleted", deletedValue)

	var activatedValue bool
	if user.Activated != nil {
		activatedValue = *user.Activated
	}
	d.Set("activated", activatedValue)

	var ssoAuthValue bool
	if user.SsoAuth != nil {
		ssoAuthValue = *user.SsoAuth
	}
	d.Set("sso_auth", ssoAuthValue)

	var twoFaValue bool
	if user.TwoFa != nil {
		twoFaValue = *user.TwoFa
	}
	d.Set("two_fa", twoFaValue)

	var userTypeValue string
	if user.UserType != nil {
		userTypeValue = string(*user.UserType)
	}
	d.Set("user_type", userTypeValue)

	var isActiveValue bool
	if user.IsActive != nil {
		isActiveValue = *user.IsActive
	}
	d.Set("is_active", isActiveValue)

	// Always set configurable fields to ensure proper drift detection
	// Use empty string for nil string fields
	var nameValue string
	if user.Name != nil {
		nameValue = *user.Name
	}
	d.Set("name", nameValue)

	langValue := "en" // default value
	if user.Lang != nil {
		langValue = string(*user.Lang)
	}
	d.Set("lang", langValue)

	var phoneValue string
	if user.Phone != nil {
		phoneValue = *user.Phone
	}
	d.Set("phone", phoneValue)

	var companyValue string
	if user.Company != nil {
		companyValue = *user.Company
	}
	d.Set("company", companyValue)

	// Always set auth_types, use empty slice for nil
	authTypesValue := []string{}
	if user.AuthTypes != nil {
		authTypesValue = make([]string, len(*user.AuthTypes))
		for i, authType := range *user.AuthTypes {
			authTypesValue[i] = string(authType)
		}
	}
	d.Set("auth_types", authTypesValue)

	// Always set groups and user_role to ensure proper state refresh
	var userRoleValue []interface{}
	var groupsValue []interface{}

	if user.Groups != nil && len(*user.Groups) > 0 {
		// Set primary user role (first group for the primary client)
		group := (*user.Groups)[0]
		userRoleMap := map[string]interface{}{}
		if group.Id != nil {
			userRoleMap["id"] = *group.Id
		}
		if group.Name != nil {
			userRoleMap["name"] = string(*group.Name)
		}
		userRoleValue = []interface{}{userRoleMap}

		// Set all groups for visibility
		groupsValue = make([]interface{}, 0, len(*user.Groups))
		for _, grp := range *user.Groups {
			groupMap := map[string]interface{}{}
			if grp.Id != nil {
				groupMap["id"] = *grp.Id
			}
			if grp.Name != nil {
				groupMap["name"] = string(*grp.Name)
			}
			groupsValue = append(groupsValue, groupMap)
		}
	}
	d.Set("user_role", userRoleValue)
	d.Set("groups", groupsValue)

	// Always set client_and_roles to ensure proper state refresh
	var clientRolesValue []interface{}
	if user.ClientAndRoles != nil {
		clientRolesValue = make([]interface{}, 0, len(*user.ClientAndRoles))
		for _, cr := range *user.ClientAndRoles {
			clientRoleMap := map[string]interface{}{
				"client_id":           cr.ClientId,
				"client_company_name": cr.ClientCompanyName,
				"user_id":             cr.UserId,
				"user_roles":          cr.UserRoles,
			}
			clientRolesValue = append(clientRolesValue, clientRoleMap)
		}
	}
	d.Set("client_and_roles", clientRolesValue)

	return nil
}

func resourceIamUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert user ID: %s", err)
	}

	// Handle user details update
	// Use forceUpdate=false for normal updates to check HasChanges
	if d.HasChanges("name", "lang", "phone", "company", "auth_types") {
		updateDiags := updateUserDetails(ctx, d, m, userID, false)
		if updateDiags.HasError() {
			return updateDiags
		}
	}

	// Handle primary user role change
	if d.HasChange("user_role") {
		clientID := d.Get("client_id").(int)
		updateDiags := updateUserRoleInClient(ctx, d, m, userID, clientID)
		if updateDiags.HasError() {
			return updateDiags
		}
	}

	return resourceIamUserRead(ctx, d, m)
}

func resourceIamUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).IamClient

	userID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert user ID: %s", err)
	}

	// Revoke access from primary client
	clientID := d.Get("client_id").(int)
	resp, err := client.DeleteClientsClientIdClientUsersUserIdWithResponse(ctx, clientID, userID)
	if err != nil {
		return diag.Errorf("failed to revoke user access from primary client: %s", err)
	}

	// Accept both 200 and 204 as success for deletion
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("failed to revoke user access: status code %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	d.SetId("")
	return nil
}

// Helper functions

func updateUserDetails(ctx context.Context, d *schema.ResourceData, m interface{}, userID int, forceUpdate bool) diag.Diagnostics {
	client := m.(*Config).IamClient

	updateReq := sdk.PatchUsersUserIdJSONRequestBody{}
	hasChanges := false

	// Helper function to check fields based on context
	checkField := func(field string) bool {
		if forceUpdate {
			// During creation: check if user provided the field
			_, ok := d.GetOk(field)
			return ok
		} else {
			// During update: check if field changed
			return d.HasChange(field)
		}
	}

	if checkField("name") {
		name := d.Get("name").(string)
		updateReq.Name = &name
		hasChanges = true
	}

	if checkField("lang") {
		lang := d.Get("lang").(string)
		langTyped := sdk.PatchUsersUserIdJSONBodyLang(lang)
		//return diag.Errorf("lang is %s and langTyped is %s", lang, langTyped)
		updateReq.Lang = &langTyped
		hasChanges = true
	}

	if checkField("phone") {
		phone := d.Get("phone").(string)
		updateReq.Phone = &phone
		hasChanges = true
	}

	if checkField("company") {
		company := d.Get("company").(string)
		updateReq.Company = &company
		hasChanges = true
	}

	if checkField("auth_types") {
		authTypesRaw := d.Get("auth_types").(*schema.Set)
		authTypes := make([]sdk.PatchUsersUserIdJSONBodyAuthTypes, 0, authTypesRaw.Len())
		for _, authType := range authTypesRaw.List() {
			authTypes = append(authTypes, sdk.PatchUsersUserIdJSONBodyAuthTypes(authType.(string)))
		}
		updateReq.AuthTypes = &authTypes
		hasChanges = true
	}

	if hasChanges {
		resp, err := client.PatchUsersUserIdWithResponse(ctx, userID, updateReq)
		if err != nil {
			return diag.Errorf("failed to update user: %s", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return diag.Errorf("failed to update user: status code %d, body: %s", resp.StatusCode(), string(resp.Body))
		}
	}

	return nil
}

func updateUserRoleInClient(ctx context.Context, d *schema.ResourceData, m interface{}, userID int, clientID int) diag.Diagnostics {
	client := m.(*Config).IamClient

	userRoleConfigs := d.Get("user_role").([]interface{})
	userRoleConfig := userRoleConfigs[0].(map[string]interface{})

	changeRoleReq := sdk.PutClientsClientIdClientUsersUserIdJSONRequestBody{
		UserRole: &struct {
			Id   *int                                                         `json:"id,omitempty"`
			Name *sdk.PutClientsClientIdClientUsersUserIdJSONBodyUserRoleName `json:"name,omitempty"`
		}{
			Id:   intPtr(userRoleConfig["id"].(int)),
			Name: (*sdk.PutClientsClientIdClientUsersUserIdJSONBodyUserRoleName)(stringPtr(userRoleConfig["name"].(string))),
		},
	}

	resp, err := client.PutClientsClientIdClientUsersUserIdWithResponse(ctx, clientID, userID, changeRoleReq)
	if err != nil {
		return diag.Errorf("failed to change user role: %s", err)
	}

	// Accept both 200 and 204 as success
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("failed to change user role: status code %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	return nil
}

// Helper function to convert int to *int
func intPtr(i int) *int {
	return &i
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}
