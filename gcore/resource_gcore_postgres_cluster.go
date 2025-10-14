package gcore

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/dbaas/postgres/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	postgresClustersPoint = "dbaas/postgres/clusters"

	// timeout for creating or updating a PostgreSQL cluster, in seconds (1h)
	postgresClusterOpTimeoutSecs = 3600

	defaultPGConfigSettings = `
huge_pages=off
max_connections=100
shared_buffers=256MB
effective_cache_size=768MB
maintenance_work_mem=64MB
work_mem=2MB
checkpoint_completion_target=0.9
wal_buffers=-1
min_wal_size=1GB
max_wal_size=4GB
random_page_cost=1.2
effective_io_concurrency=200
`
)

var postgresClusterOperationTimeout = time.Second * time.Duration(postgresClusterOpTimeoutSecs)

func resourcePostgresCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePostgresClusterCreate,
		ReadContext:   resourcePostgresClusterRead,
		UpdateContext: resourcePostgresClusterUpdate,
		DeleteContext: resourcePostgresClusterDelete,
		Description:   "Represents a PostgreSQL cluster resource in Gcore Cloud.",
		Timeouts: &schema.ResourceTimeout{
			Create: &postgresClusterOperationTimeout,
			Update: &postgresClusterOperationTimeout,
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, clusterName, err := ImportStringParser(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(clusterName)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema:        resourceSchema(),
		CustomizeDiff: validateDatabaseOwners,
	}
}

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeInt,
			Optional: true,
			ExactlyOneOf: []string{
				"project_id",
				"project_name",
			},
			DiffSuppressFunc: suppressDiffProjectID,
		},
		"region_id": {
			Type:     schema.TypeInt,
			Optional: true,
			ExactlyOneOf: []string{
				"region_id",
				"region_name",
			},
			DiffSuppressFunc: suppressDiffRegionID,
		},
		"project_name": {
			Type:     schema.TypeString,
			Optional: true,
			ExactlyOneOf: []string{
				"project_id",
				"project_name",
			},
		},
		"region_name": {
			Type:     schema.TypeString,
			Optional: true,
			ExactlyOneOf: []string{
				"region_id",
				"region_name",
			},
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Cluster name.",
			Required:    true,
			ForceNew:    true,
		},
		"flavor": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Flavor of the cluster instance.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:        schema.TypeInt,
						Description: "Number of CPU cores.",
						Required:    true,
					},
					"memory": {
						Type:        schema.TypeInt,
						Description: "Amount of RAM in GiB.",
						Required:    true,
					},
				},
			},
		},
		"database": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "List of databases to create in the cluster.",
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "Database name.",
						Required:    true,
					},
					"owner": {
						Type:        schema.TypeString,
						Description: "Owner of the database. Must be one of the users defined in the `users` block.",
						Required:    true,
					},
					"size": {
						Type:        schema.TypeInt,
						Description: "Size of the database in MiB.",
						Computed:    true,
					},
				},
			},
			Set: func(v interface{}) int {
				m := v.(map[string]interface{})
				return schema.HashString(m["name"].(string) + "|" + m["owner"].(string))
			},
		},
		"ha_replication_mode": {
			Type:             schema.TypeString,
			Description:      "Replication mode. Possible values are `async` and `sync`.",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"async", "sync"}, false)),
		},
		"network": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Network configuration.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"acl": {
						Type:        schema.TypeSet,
						Description: "List of IP addresses or CIDR blocks allowed to access the cluster.",
						Required:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Set: schema.HashString,
					},
					"network_type": {
						Type:             schema.TypeString,
						Description:      "Network type. Currently, only `public` is supported.",
						Optional:         true,
						Default:          "public",
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"public"}, false)),
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
			Required:    true,
			Description: "PostgreSQL cluster configuration.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"version": {
						Type:        schema.TypeString,
						Description: "PostgreSQL version. Possible values are `13`, `14`, and `15`.",
						Required:    true,
					},
					"pg_conf": {
						Type:        schema.TypeString,
						Description: "PostgreSQL configuration in `key=value` format, one per line.",
						Optional:    true,
						Default:     defaultPGConfigSettings,
						DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
							return suppressDiffMultiline(k, old, new, d)
						},
					},
					"pooler_mode": {
						Type:             schema.TypeString,
						Description:      "Connection pooler mode. Possible values are `session`, `transaction`, and `statement`. If not set, connection pooler is not enabled.",
						Optional:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"session", "transaction", "statement"}, false)),
					},
					"pooler_type": {
						Type:             schema.TypeString,
						Description:      "Connection pooler type. Currently, only `pgbouncer` is supported.",
						Optional:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"pgbouncer"}, false)),
					},
				},
			},
		},
		"storage": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Storage configuration.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"size": {
						Type:        schema.TypeInt,
						Description: "Storage size in GiB. Must be between 1 and 100.",
						Required:    true,
					},
					"type": {
						Type:        schema.TypeString,
						Description: "Storage type.",
						Required:    true,
						ForceNew:    true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"ssd-hiiops", "ssd-lowlatency", "standard"}, false)),
					},
				},
			},
		},
		"user": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "List of users to create in the cluster.",
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "User name.",
						Required:    true,
					},
					"role_attributes": {
						Type: schema.TypeSet,
						Description: fmt.Sprintf("List of role attributes. Possible values are: %s.",
							strings.Join(clusters.RoleAttributeStringList(), ", ")),
						Required: true,
						Elem: &schema.Schema{
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(clusters.RoleAttributeStringList(), false)),
						},
						Set: schema.HashString,
					},
					"is_secret_revealed": {
						Type:        schema.TypeBool,
						Description: "Whether the password for the default `postgres` user is revealed.",
						Computed:    true,
					},
				},
			},
			Set: func(v interface{}) int {
				m := v.(map[string]interface{})
				roleAttrSet := m["role_attributes"].(*schema.Set)
				roleAttrs := roleAttrSet.List()
				roleAttrsStrs := make([]string, len(roleAttrs))
				for i, v := range roleAttrs {
					roleAttrsStrs[i] = v.(string)
				}
				// sort role attributes to ensure consistent hash
				sort.Strings(roleAttrsStrs)
				return schema.HashString(m["name"].(string) + "|" + strings.Join(roleAttrsStrs, ","))
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
	}
}

func resourcePostgresClusterCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start postgres cluster create")
	config := m.(*Config)
	provider := config.Provider

	pgClient, err := CreateClient(provider, data, postgresClustersPoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts, err := expandClusterCreateOpts(data)
	if err != nil {
		return diag.FromErr(err)
	}
	result := clusters.Create(pgClient, createOpts)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}
	if err = waitForTaskResult(result, postgresClusterOpTimeoutSecs, provider, data); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(createOpts.ClusterName)

	log.Println("[DEBUG] Finished postgres cluster create")
	return resourcePostgresClusterRead(ctx, data, m)
}

func expandClusterCreateOpts(data *schema.ResourceData) (*clusters.CreateOpts, error) {
	opts := clusters.CreateOpts{ClusterName: data.Get("name").(string)}

	// extract flavors
	flavorList := data.Get("flavor").([]interface{})
	if len(flavorList) > 0 {
		flavorMap := flavorList[0].(map[string]interface{})
		opts.Flavor = clusters.FlavorOpts{
			CPU:       flavorMap["cpu"].(int),
			MemoryGiB: flavorMap["memory"].(int),
		}
	} else {
		return nil, fmt.Errorf("flavor must be specified")
	}

	// extract databases
	databases := make([]clusters.DatabaseOpts, 0)
	for _, db := range data.Get("database").(*schema.Set).List() {
		dbMap := db.(map[string]interface{})
		databases = append(databases, clusters.DatabaseOpts{
			Name:  dbMap["name"].(string),
			Owner: dbMap["owner"].(string),
		})
	}
	if len(databases) == 0 {
		return nil, fmt.Errorf("at least one database must be specified")
	}
	opts.Databases = databases

	// extract network
	networkList := data.Get("network").([]interface{})
	if len(networkList) > 0 {
		networkMap := networkList[0].(map[string]interface{})
		aclSet := networkMap["acl"].(*schema.Set)
		acl := make([]string, 0)
		for _, v := range aclSet.List() {
			acl = append(acl, v.(string))
		}
		opts.Network = clusters.NetworkOpts{
			ACL:         acl,
			NetworkType: networkMap["network_type"].(string),
		}
	} else {
		return nil, fmt.Errorf("network must be specified")
	}

	// extract pg_config
	pgConfigList := data.Get("pg_config").([]interface{})
	if len(pgConfigList) > 0 {
		pgConfigMap := pgConfigList[0].(map[string]interface{})
		opts.PGServerConfiguration = clusters.PGServerConfigurationOpts{
			PGConf:  pgConfigMap["pg_conf"].(string),
			Version: pgConfigMap["version"].(string),
		}
		if poolerMode, ok := pgConfigMap["pooler_mode"].(string); ok && poolerMode != "" {
			opts.PGServerConfiguration.Pooler = &clusters.PoolerOpts{
				Mode: clusters.PoolerMode(poolerMode),
				Type: pgConfigMap["pooler_type"].(string),
			}
		}
	} else {
		return nil, fmt.Errorf("pg_config must be specified")
	}

	// extract storage
	storageList := data.Get("storage").([]interface{})
	if len(storageList) > 0 {
		storageMap := storageList[0].(map[string]interface{})
		opts.Storage = clusters.PGStorageConfigurationOpts{
			SizeGiB: storageMap["size"].(int),
			Type:    storageMap["type"].(string),
		}
	} else {
		return nil, fmt.Errorf("storage must be specified")
	}

	// extract users
	users := make([]clusters.PgUserOpts, 0)
	for _, user := range data.Get("user").(*schema.Set).List() {
		userMap := user.(map[string]interface{})
		roleAttrSet := userMap["role_attributes"].(*schema.Set)
		roleAttributes := make([]clusters.RoleAttribute, 0)
		for _, v := range roleAttrSet.List() {
			roleAttributes = append(roleAttributes, clusters.RoleAttribute(v.(string)))
		}
		users = append(users, clusters.PgUserOpts{
			Name:           userMap["name"].(string),
			RoleAttributes: roleAttributes,
		})
	}
	if len(users) > 0 {
		opts.Users = users
	}

	// extract ha_replication_mode
	if v, ok := data.GetOk("ha_replication_mode"); ok {
		opts.HighAvailability = &clusters.HighAvailabilityOpts{
			ReplicationMode: clusters.HighAvailabilityReplicationMode(v.(string)),
		}
	}

	return &opts, nil
}

func resourcePostgresClusterRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start postgres cluster read")
	config := m.(*Config)
	provider := config.Provider
	clusterName := data.Id()
	pgClient, err := CreateClient(provider, data, postgresClustersPoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	cluster, err := clusters.Get(pgClient, clusterName).Extract()
	if err != nil {
		var errDefault404 gcorecloud.ErrDefault404
		if errors.As(err, &errDefault404) {
			// removing from state because it doesn't exist anymore
			data.SetId("")
			return nil
		}
		return diag.Errorf("cannot get postgres cluster with ID: %s. Error: %s", clusterName, err)
	}
	data.Set("name", cluster.ClusterName)
	data.Set("status", cluster.Status)
	data.Set("created_at", cluster.CreatedAt.String())

	flavor := make(map[string]interface{})
	flavor["cpu"] = cluster.Flavor.CPU
	flavor["memory"] = cluster.Flavor.MemoryGiB
	data.Set("flavor", []interface{}{flavor})

	databases := make([]map[string]interface{}, 0, len(cluster.Databases))
	for _, db := range cluster.Databases {
		database := make(map[string]interface{})
		database["name"] = db.Name
		database["owner"] = db.Owner
		database["size"] = db.Size
		databases = append(databases, database)
	}
	data.Set("database", databases)

	if cluster.HighAvailability != nil {
		data.Set("ha_replication_mode", cluster.HighAvailability.ReplicationMode)
	}

	network := make(map[string]interface{})
	network["acl"] = cluster.Network.ACL
	network["network_type"] = cluster.Network.NetworkType
	network["connection_string"] = cluster.Network.ConnectionString
	network["host"] = cluster.Network.Host
	data.Set("network", []interface{}{network})

	pgConfig := make(map[string]interface{})
	pgConfig["version"] = cluster.PGServerConfiguration.Version
	pgConfig["pg_conf"] = cluster.PGServerConfiguration.PGConf
	if cluster.PGServerConfiguration.Pooler != nil {
		pgConfig["pooler_mode"] = cluster.PGServerConfiguration.Pooler.Mode
		pgConfig["pooler_type"] = cluster.PGServerConfiguration.Pooler.Type
	}
	data.Set("pg_config", []interface{}{pgConfig})

	storage := make(map[string]interface{})
	storage["size"] = cluster.Storage.SizeGiB
	storage["type"] = cluster.Storage.Type
	data.Set("storage", []interface{}{storage})

	users := make([]map[string]interface{}, 0, len(cluster.Users))
	for _, user := range cluster.Users {
		u := make(map[string]interface{})
		u["name"] = user.Name
		u["role_attributes"] = clusters.RoleAttributeSliceToStrings(user.RoleAttributes)
		u["is_secret_revealed"] = user.IsSecretRevealed
		users = append(users, u)
	}
	data.Set("user", users)
	log.Printf("[DEBUG] Read postgres cluster %s", data.Id())
	return nil
}

func resourcePostgresClusterUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start postgres cluster update")
	config := m.(*Config)
	provider := config.Provider
	clusterName := data.Id()
	pgClient, err := CreateClient(provider, data, postgresClustersPoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := clusters.UpdateOpts{}

	if data.HasChange("storage") {
		storageList := data.Get("storage").([]interface{})
		if len(storageList) > 0 {
			storageMap := storageList[0].(map[string]interface{})
			updateOpts.Storage = &clusters.PGStorageConfigurationUpdateOpts{
				SizeGiB: storageMap["size"].(int),
			}
		} else {
			return diag.FromErr(fmt.Errorf("storage must be specified"))
		}
	}

	if data.HasChange("pg_config") {
		pgConfigList := data.Get("pg_config").([]interface{})
		if len(pgConfigList) > 0 {
			pgConfigMap := pgConfigList[0].(map[string]interface{})
			updateOpts.PGServerConfiguration = &clusters.PGServerConfigurationUpdateOpts{
				PGConf:  pgConfigMap["pg_conf"].(string),
				Version: pgConfigMap["version"].(string),
			}
			if poolerMode, ok := pgConfigMap["pooler_mode"].(string); ok && poolerMode != "" {
				updateOpts.PGServerConfiguration.Pooler = &clusters.PoolerOpts{
					Mode: clusters.PoolerMode(poolerMode),
					Type: pgConfigMap["pooler_type"].(string),
				}
			}
		} else {
			return diag.FromErr(fmt.Errorf("pg_config must be specified"))
		}
	}

	if data.HasChange("ha_replication_mode") {
		if v, ok := data.GetOk("ha_replication_mode"); ok {
			updateOpts.HighAvailability = &clusters.HighAvailabilityOpts{
				ReplicationMode: clusters.HighAvailabilityReplicationMode(v.(string)),
			}
		} else {
			updateOpts.HighAvailability = nil
		}
	}

	if data.HasChange("database") {
		databases := make([]clusters.DatabaseOpts, 0)
		for _, db := range data.Get("database").(*schema.Set).List() {
			dbMap := db.(map[string]interface{})
			databases = append(databases, clusters.DatabaseOpts{
				Name:  dbMap["name"].(string),
				Owner: dbMap["owner"].(string),
			})
		}
		if len(databases) > 0 {
			updateOpts.Databases = databases
		}
	}

	if data.HasChange("user") {
		users := make([]clusters.PgUserOpts, 0)
		for _, user := range data.Get("user").(*schema.Set).List() {
			userMap := user.(map[string]interface{})
			roleAttrSet := userMap["role_attributes"].(*schema.Set)
			roleAttributes := make([]clusters.RoleAttribute, 0)
			for _, v := range roleAttrSet.List() {
				roleAttributes = append(roleAttributes, clusters.RoleAttribute(v.(string)))
			}
			users = append(users, clusters.PgUserOpts{
				Name:           userMap["name"].(string),
				RoleAttributes: roleAttributes,
			})
		}
		if len(users) > 0 {
			updateOpts.Users = users
		}
	}

	if data.HasChange("flavor") {
		flavorList := data.Get("flavor").([]interface{})
		if len(flavorList) > 0 {
			flavorMap := flavorList[0].(map[string]interface{})
			updateOpts.Flavor = &clusters.FlavorOpts{
				CPU:       flavorMap["cpu"].(int),
				MemoryGiB: flavorMap["memory"].(int),
			}
		} else {
			return diag.FromErr(fmt.Errorf("flavor must be specified"))
		}
	}

	if data.HasChange("network") {
		networkList := data.Get("network").([]interface{})
		if len(networkList) > 0 {
			networkMap := networkList[0].(map[string]interface{})
			aclSet := networkMap["acl"].(*schema.Set)
			acl := make([]string, 0)
			for _, v := range aclSet.List() {
				acl = append(acl, v.(string))
			}
			updateOpts.Network = &clusters.NetworkOpts{
				ACL:         acl,
				NetworkType: networkMap["network_type"].(string),
			}
		} else {
			return diag.FromErr(fmt.Errorf("network must be specified"))
		}
	}

	result := clusters.Update(pgClient, clusterName, updateOpts)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	if err = waitForTaskResult(result, postgresClusterOpTimeoutSecs, provider, data); err != nil {
		return diag.FromErr(err)
	}

	return resourcePostgresClusterRead(ctx, data, m)
}

func resourcePostgresClusterDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start postgres cluster delete")
	config := m.(*Config)
	provider := config.Provider
	clusterName := data.Id()
	pgClient, err := CreateClient(provider, data, postgresClustersPoint, "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	result := clusters.Delete(pgClient, clusterName, clusters.DeleteOpts{})
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	taskResults, err := result.Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := taskResults.Tasks[0]

	tasksClient, err := CreateClient(provider, data, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tasks.WaitTaskAndProcessResult(tasksClient, taskID, true, postgresClusterOpTimeoutSecs, func(task tasks.TaskID) error {
		_, err := clusters.Get(pgClient, clusterName).Extract()
		if err == nil {
			return fmt.Errorf("cannot delete cluster with name: %s", clusterName)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil
		default:
			return err
		}
	})
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("")
	log.Println("[DEBUG] End postgres cluster delete")
	return nil
}

func suppressDiffMultiline(k, old, new string, d *schema.ResourceData) bool {
	// normalize function converts multiline string in "key=value" format to map
	// ignores empty lines and comments, trims spaces around keys and values
	// if a line is malformed (does not contain '='), the whole line is kept as key with empty value
	// this way, any change in malformed lines will be detected
	normalize := func(s string) map[string]string {
		m := make(map[string]string)
		sc := bufio.NewScanner(strings.NewReader(s))
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				// keep whole line to detect changes in malformed lines
				m[line] = ""
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			m[key] = val
		}
		return m
	}

	oldMap := normalize(old)
	newMap := normalize(new)

	if len(oldMap) != len(newMap) {
		return false
	}
	for key, oldVal := range oldMap {
		if newMap[key] != oldVal {
			return false
		}
	}
	return true
}

func validateDatabaseOwners(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	dbRaw := diff.Get("database")
	userRaw := diff.Get("user")
	if dbRaw == nil || userRaw == nil {
		return nil
	}

	userSet, ok := userRaw.(*schema.Set)
	if !ok {
		return nil
	}
	userNames := make(map[string]struct{}, userSet.Len())
	for _, u := range userSet.List() {
		umap := u.(map[string]interface{})
		if name, ok := umap["name"].(string); ok && name != "" {
			userNames[name] = struct{}{}
		}
	}

	dbSet, ok := dbRaw.(*schema.Set)
	if !ok {
		return nil
	}
	var err error
	for _, db := range dbSet.List() {
		if db == nil {
			continue
		}
		dbMap := db.(map[string]interface{})
		owner, _ := dbMap["owner"].(string)
		if _, found := userNames[owner]; !found {
			return fmt.Errorf("database owner %q not found among defined user names. "+
				"Add a user with name %q in the user block", owner, owner)
		}
	}
	return err
}
