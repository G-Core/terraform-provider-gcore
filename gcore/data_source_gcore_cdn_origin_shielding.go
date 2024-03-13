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
			"city": {
				Type:        schema.TypeString,
				Description: "Displayed shielding location point of present",
				Required:    true,
			},
		},
	}
}

func getLocationByCity(arr []originshielding.OriginShieldingLocations, city string) (int, error) {
	for _, el := range arr {
		if el.City == city {
			return el.ID, nil
		}
	}
	return 0, fmt.Errorf("shielding location for city %s not found", city)
}

func dataOriginShieldingLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading origin shielding locations")

	city := d.Get("city").(string)
	config := m.(*Config)
	client := config.CDNClient

	result, err := client.OriginShielding().GetLocations(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Shielding locations: %v", *result)
	locationID, err := getLocationByCity(*result, city)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(locationID))
	d.Set("city", city)

	log.Println("[DEBUG] Finish reading origin shielding locations")
	return nil

}
