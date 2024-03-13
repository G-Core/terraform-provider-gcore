package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/originshielding"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataOriginShieldingLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataOriginShieldingLocationRead,
		Description: "Represent shielding locations data",
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Description: "Displayed shielding location for point of present",
				Required:    true,
			},
		},
	}
}

func getLocationByDatacenter(arr []originshielding.OriginShieldingLocations, datacenter string) (int, error) {
	for _, el := range arr {
		if el.Datacenter == datacenter {
			return el.ID, nil
		}
	}
	return 0, fmt.Errorf("shielding location for datacenter %s not found", datacenter)
}

func dataOriginShieldingLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading origin shielding locations")

	datacenter := d.Get("datacenter").(string)
	config := m.(*Config)
	client := config.CDNClient

	result, err := client.OriginShielding().GetLocations(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Shielding locations: %v", *result)
	locationID, err := getLocationByDatacenter(*result, datacenter)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(locationID))
	err = d.Set("datacenter", datacenter)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish reading origin shielding locations")
	return nil

}
