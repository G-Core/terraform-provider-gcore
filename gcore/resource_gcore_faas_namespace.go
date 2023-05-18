package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	faasPoint = "faas/namespaces"

	faaSNamespaceCreateTimeout = 1200
	faaSNamespaceDeleteTimeout = 1200
)

func resourceFaaSNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFaaSNamespaceCreate,
		ReadContext:   resourceFaaSNamespaceRead,
		UpdateContext: resourceFaaSNamespaceUpdate,
		DeleteContext: resourceFaaSNamespaceDelete,
		Description:   "Represent FaaS namespace",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, nsName, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(nsName)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"envs": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFaaSNamespaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS namespace creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	nsName := d.Get("name").(string)
	opts := faas.CreateNamespaceOpts{
		Name:        nsName,
		Description: d.Get("description").(string),
		Envs:        map[string]string{},
	}
	envsRaw := d.Get("envs").(map[string]interface{})
	if len(envsRaw) > 0 {
		envs := make(map[string]string, len(envsRaw))
		for k, v := range envsRaw {
			envs[k] = v.(string)
		}
		opts.Envs = envs
	}

	results, err := faas.CreateNamespace(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, faaSNamespaceCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return "", nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(nsName)
	resourceFaaSNamespaceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish FaaS namespace creating (%s)", nsName)
	return diags
}

func resourceFaaSNamespaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS namespace reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	nsName := d.Id()
	log.Printf("[DEBUG] namespace = %s", nsName)

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	ns, err := faas.GetNamespace(client, nsName).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			log.Printf("[WARN] Removing namesapce %s because resource doesn't exist anymore", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", ns.Name)
	d.Set("description", ns.Description)
	d.Set("status", ns.Status)
	d.Set("created_at", ns.CreatedAt.Format(time.RFC3339))

	if err := d.Set("envs", ns.Envs); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish FaaS namespace reading")
	return diags
}

func resourceFaaSNamespaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS namespace updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	nsName := d.Id()
	var needUpdate bool
	opts := faas.UpdateNamespaceOpts{}
	if d.HasChange("envs") {
		envsRaw := d.Get("envs").(map[string]interface{})
		envs := make(map[string]string, len(envsRaw))
		for k, v := range envsRaw {
			envs[k] = v.(string)
		}

		opts.Envs = envs
		needUpdate = true
	}

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
		needUpdate = true
	}

	if needUpdate {
		results, err := faas.UpdateNamespace(client, nsName, opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		taskID := results.Tasks[0]
		_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, faaSNamespaceCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
			_, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			return "", nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Println("[DEBUG] Finish FaaS namespace updating")
	return resourceFaaSNamespaceRead(ctx, d, m)
}

func resourceFaaSNamespaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS namespace deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	nsName := d.Id()
	log.Printf("[DEBUG] Namespace = %s", nsName)

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := faas.DeleteNamespace(client, nsName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, faaSNamespaceDeleteTimeout, func(task tasks.TaskID) (interface{}, error) {
		// workaround because of cloud-api
		time.Sleep(time.Second * 5)
		_, err := faas.GetNamespace(client, nsName).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete namespace: %s", nsName)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of FaaS namespace deleting")
	return diags
}
