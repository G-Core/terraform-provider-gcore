package gcore

import (
	"context"
	"errors"
	"log"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/dbaas/postgres/v1/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePostgresCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePostgresClusterRead,
		Description: "Get information about a PostgreSQL cluster in Gcore Cloud.",
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ExactlyOneOf:     []string{"project_id", "project_name"},
				DiffSuppressFunc: suppressDiffProjectID,
				Description:      "Project ID, only one of project_id or project_name should be set",
			},
			"region_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ExactlyOneOf:     []string{"region_id", "region_name"},
				DiffSuppressFunc: suppressDiffRegionID,
				Description:      "Region ID, only one of region_id or region_name should be set",
			},
			"project_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"project_id", "project_name"},
				Description:  "Project name, only one of project_id or project_name should be set",
			},
			"region_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"region_id", "region_name"},
				Description:  "Region name, only one of region_id or region_name should be set",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the PostgreSQL cluster. It must be unique within the project and region.",
			},
			"flavor": {
				Type:        schema.TypeList,
				Description: "Flavor of the cluster instance.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeInt,
							Description: "Number of CPU cores.",
							Computed:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Amount of RAM in GiB.",
							Computed:    true,
						},
					},
				},
			},
			"database": {
				Type:        schema.TypeSet,
				Description: "List of databases in the cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Database name.",
							Computed:    true,
						},
						"owner": {
							Type:        schema.TypeString,
							Description: "Owner of the database.",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "Size of the database in MiB.",
							Computed:    true,
						},
					},
				},
			},
			"ha_replication_mode": {
				Type:        schema.TypeString,
				Description: "Replication mode. Possible values are `async` and `sync`.",
				Computed:    true,
			},
			"network": {
				Type:        schema.TypeList,
				Description: "Network configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl": {
							Type:        schema.TypeSet,
							Description: "List of IP addresses or CIDR blocks allowed to access the cluster.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"network_type": {
							Type:        schema.TypeString,
							Description: "Network type. Currently, only `public` is supported.",
							Computed:    true,
						},
						"connection_string": {
							Type:        schema.TypeString,
							Description: "Connection string for the cluster.",
							Computed:    true,
						},
						"host": {
							Type:        schema.TypeString,
							Description: "Host address for the cluster.",
							Computed:    true,
						},
					},
				},
			},
			"pg_config": {
				Type:        schema.TypeList,
				Description: "PostgreSQL cluster configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Description: "PostgreSQL version. Possible values are `13`, `14`, and `15`.",
							Computed:    true,
						},
						"pg_conf": {
							Type:        schema.TypeString,
							Description: "PostgreSQL configuration in `key=value` format, one per line.",
							Computed:    true,
						},
						"pooler_mode": {
							Type:        schema.TypeString,
							Description: "Connection pooler mode. Possible values are `session`, `transaction`, and `statement`. If not set, connection pooler is not enabled.",
							Computed:    true,
						},
						"pooler_type": {
							Type:        schema.TypeString,
							Description: "Connection pooler type. Currently, only `pgbouncer` is supported.",
							Computed:    true,
						},
					},
				},
			},
			"storage": {
				Type:        schema.TypeList,
				Description: "Storage configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Description: "Storage size in GiB.",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Storage type.",
							Computed:    true,
						},
					},
				},
			},
			"user": {
				Type:        schema.TypeSet,
				Description: "List of users in the cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "User name.",
							Computed:    true,
						},
						"role_attributes": {
							Type:        schema.TypeSet,
							Description: "List of role attributes.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_secret_revealed": {
							Type:        schema.TypeBool,
							Description: "Whether the password for the default `postgres` user is revealed.",
							Computed:    true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the cluster.",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Cluster creation date.",
				Computed:    true,
			},
		},
	}
}

func dataSourcePostgresClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start PostgreSQL cluster data source reading")
	config := m.(*Config)
	provider := config.Provider

	projectID, regionID, err := getProjectAndRegionID(provider, d)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	client, err := CreateClient(provider, d, postgresClustersPoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	cluster, err := clusters.Get(client, name).Extract()
	if err != nil {
		var errDefault404 gcorecloud.ErrDefault404
		if errors.As(err, &errDefault404) {
			return diag.Errorf("PostgreSQL cluster with name '%s' not found in project %d and region %d", name, projectID, regionID)
		}
		return diag.FromErr(err)
	}

	d.SetId(cluster.ClusterName)
	d.Set("name", cluster.ClusterName)
	d.Set("status", cluster.Status)
	d.Set("created_at", cluster.CreatedAt.String())

	// Set flavor
	flavor := make(map[string]interface{})
	flavor["cpu"] = cluster.Flavor.CPU
	flavor["memory"] = cluster.Flavor.MemoryGiB
	d.Set("flavor", []interface{}{flavor})

	// Set databases
	databases := make([]map[string]interface{}, 0, len(cluster.Databases))
	for _, db := range cluster.Databases {
		database := make(map[string]interface{})
		database["name"] = db.Name
		database["owner"] = db.Owner
		database["size"] = db.Size
		databases = append(databases, database)
	}
	d.Set("database", databases)

	// Set high availability replication mode
	if cluster.HighAvailability != nil {
		d.Set("ha_replication_mode", cluster.HighAvailability.ReplicationMode)
	}

	// Set network
	network := make(map[string]interface{})
	network["acl"] = cluster.Network.ACL
	network["network_type"] = cluster.Network.NetworkType
	network["connection_string"] = cluster.Network.ConnectionString
	network["host"] = cluster.Network.Host
	d.Set("network", []interface{}{network})

	// Set pg_config
	pgConfig := make(map[string]interface{})
	pgConfig["version"] = cluster.PGServerConfiguration.Version
	pgConfig["pg_conf"] = cluster.PGServerConfiguration.PGConf
	if cluster.PGServerConfiguration.Pooler != nil {
		pgConfig["pooler_mode"] = cluster.PGServerConfiguration.Pooler.Mode
		pgConfig["pooler_type"] = cluster.PGServerConfiguration.Pooler.Type
	}
	d.Set("pg_config", []interface{}{pgConfig})

	// Set storage
	storage := make(map[string]interface{})
	storage["size"] = cluster.Storage.SizeGiB
	storage["type"] = cluster.Storage.Type
	d.Set("storage", []interface{}{storage})

	// Set users
	users := make([]map[string]interface{}, 0, len(cluster.Users))
	for _, user := range cluster.Users {
		u := make(map[string]interface{})
		u["name"] = user.Name
		u["role_attributes"] = clusters.RoleAttributeSliceToStrings(user.RoleAttributes)
		u["is_secret_revealed"] = user.IsSecretRevealed
		users = append(users, u)
	}
	d.Set("user", users)

	log.Printf("[DEBUG] Read PostgreSQL cluster %s", d.Id())
	return nil
}
