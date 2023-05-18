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
	faaSFunctionCreateTimeout = 2400
	faaSFunctionDeleteTimeout = 2400
)

func resourceFaaSFunction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFaaSFunctionCreate,
		ReadContext:   resourceFaaSFunctionRead,
		UpdateContext: resourceFaaSFunctionUpdate,
		DeleteContext: resourceFaaSFunctionDelete,
		Description:   "Represent FaaS function",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, nsName, fName, err := ImportStringParserExtended(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(funcID(fName, nsName))

				d.Set("name", fName)
				d.Set("namespace", nsName)

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
			"namespace": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Namespace of the function",
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
			"runtime": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"code_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Autoscaling max number of instances",
				Required:    true,
			},
			"min_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Autoscaling min number of instances",
				Required:    true,
			},
			"main_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Main startup method name",
				Required:    true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"build_message": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deploy_status": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFaaSFunctionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS function creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	fName := d.Get("name").(string)
	nsName := d.Get("namespace").(string)
	opts := faas.CreateFunctionOpts{
		Name:        fName,
		Description: d.Get("description").(string),
		Runtime:     d.Get("runtime").(string),
		Timeout:     d.Get("timeout").(int),
		Flavor:      d.Get("flavor").(string),
		Autoscaling: faas.FunctionAutoscaling{
			MinInstances: d.Get("min_instances").(int),
			MaxInstances: d.Get("max_instances").(int),
		},
		CodeText:   d.Get("code_text").(string),
		MainMethod: d.Get("main_method").(string),
		Envs:       map[string]string{},
	}
	envsRaw := d.Get("envs").(map[string]interface{})
	if len(envsRaw) > 0 {
		envs := make(map[string]string, len(envsRaw))
		for k, v := range envsRaw {
			envs[k] = v.(string)
		}
		opts.Envs = envs
	}

	results, err := faas.CreateFunction(client, nsName, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, faaSFunctionCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return "", nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(funcID(fName, nsName))

	resourceFaaSFunctionRead(ctx, d, m)

	log.Printf("[DEBUG] Finish FaaS function creating (%s)", fName)
	return diags
}

func resourceFaaSFunctionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS function reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	fName := d.Get("name").(string)
	nsName := d.Get("namespace").(string)
	log.Printf("[DEBUG] function = %s in %s", fName, nsName)

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	function, err := faas.GetFunction(client, nsName, fName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := faaSSetState(d, function); err != nil {
		diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish FaaS function reading")
	return diags
}

func resourceFaaSFunctionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS function updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	fName := d.Get("name").(string)
	nsName := d.Get("namespace").(string)

	var needUpdate bool
	opts := faas.UpdateFunctionOpts{}
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

	if d.HasChange("code_text") {
		opts.CodeText = d.Get("code_text").(string)
		needUpdate = true
	}

	if d.HasChange("timeout") {
		opts.Timeout = d.Get("timeout").(int)
		needUpdate = true
	}

	if d.HasChange("main_method") {
		opts.MainMethod = d.Get("main_method").(string)
		needUpdate = true
	}

	if d.HasChange("flavor") {
		opts.Flavor = d.Get("flavor").(string)
		needUpdate = true
	}

	if d.HasChanges("max_instances", "max_instances") {
		opts.Autoscaling = &faas.FunctionAutoscaling{
			MinInstances: d.Get("min_instances").(int),
			MaxInstances: d.Get("max_instances").(int),
		}
		needUpdate = true
	}

	if needUpdate {
		results, err := faas.UpdateFunction(client, nsName, fName, opts).Extract()
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

	log.Println("[DEBUG] Finish FaaS function updating")
	return resourceFaaSFunctionRead(ctx, d, m)
}

func resourceFaaSFunctionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS function deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	fName := d.Get("name").(string)
	nsName := d.Get("namespace").(string)

	log.Printf("[DEBUG] function = %s", fName)

	client, err := CreateClient(provider, d, faasPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := faas.DeleteFunction(client, nsName, fName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, faaSFunctionDeleteTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := faas.GetFunction(client, nsName, fName).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete function: %s", fName)
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
	log.Printf("[DEBUG] Finish of FaaS function deleting")
	return diags
}

func funcID(fName, nsName string) string {
	return fmt.Sprintf("%s_%s", fName, nsName)
}

func faaSSetState(d *schema.ResourceData, function *faas.Function) error {
	d.Set("name", function.Name)
	d.Set("description", function.Description)
	d.Set("runtime", function.Runtime)
	d.Set("code_text", function.CodeText)
	d.Set("timeout", function.Timeout)
	d.Set("max_instances", function.Autoscaling.MaxInstances)
	d.Set("min_instances", function.Autoscaling.MinInstances)
	d.Set("main_method", function.MainMethod)
	d.Set("flavor", function.Flavor)
	d.Set("build_message", function.BuildMessage)
	d.Set("build_status", function.BuildStatus)
	d.Set("status", function.Status)
	d.Set("endpoint", function.Endpoint)
	d.Set("created_at", function.CreatedAt.Format(time.RFC3339))

	if err := d.Set("envs", function.Envs); err != nil {
		return err
	}
	ds := map[string]int{
		"total": function.DeployStatus.Total,
		"ready": function.DeployStatus.Ready,
	}
	if err := d.Set("deploy_status", ds); err != nil {
		return err
	}

	return nil
}
