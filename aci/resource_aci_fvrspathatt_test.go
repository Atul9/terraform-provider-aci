package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciStaticPath_Basic(t *testing.T) {
	var static_path models.StaticPath
	description := "static_path created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciStaticPathConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists("aci_epg_to_static_path.foostatic_path", &static_path),
					testAccCheckAciStaticPathAttributes(description, &static_path),
				),
			},
			{
				ResourceName:      "aci_epg_to_static_path",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciStaticPath_update(t *testing.T) {
	var static_path models.StaticPath
	description := "static_path created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciStaticPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciStaticPathConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists("aci_epg_to_static_path.foostatic_path", &static_path),
					testAccCheckAciStaticPathAttributes(description, &static_path),
				),
			},
			{
				Config: testAccCheckAciStaticPathConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciStaticPathExists("aci_epg_to_static_path.foostatic_path", &static_path),
					testAccCheckAciStaticPathAttributes(description, &static_path),
				),
			},
		},
	})
}

func testAccCheckAciStaticPathConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_epg_to_static_path" "foostatic_path" {
		  application_epg_dn  = "${aci_application_epg.example.id}"
		description = "%s"
		
		tdn  = "example"
		  annotation  = "example"
		  encap  = "example"
		  instr_imedcy  = "example"
		  mode  = "example"
		  primary_encap  = "example"
		}
	`, description)
}

func testAccCheckAciStaticPathExists(name string, static_path *models.StaticPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Static Path %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Static Path dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		static_pathFound := models.StaticPathFromContainer(cont)
		if static_pathFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Static Path %s not found", rs.Primary.ID)
		}
		*static_path = *static_pathFound
		return nil
	}
}

func testAccCheckAciStaticPathDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_epg_to_static_path" {
			cont, err := client.Get(rs.Primary.ID)
			static_path := models.StaticPathFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Static Path %s Still exists", static_path.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciStaticPathAttributes(description string, static_path *models.StaticPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != static_path.Description {
			return fmt.Errorf("Bad static_path Description %s", static_path.Description)
		}

		if "example" != static_path.TDn {
			return fmt.Errorf("Bad static_path t_dn %s", static_path.TDn)
		}

		if "example" != static_path.Annotation {
			return fmt.Errorf("Bad static_path annotation %s", static_path.Annotation)
		}

		if "example" != static_path.Encap {
			return fmt.Errorf("Bad static_path encap %s", static_path.Encap)
		}

		if "example" != static_path.InstrImedcy {
			return fmt.Errorf("Bad static_path instr_imedcy %s", static_path.InstrImedcy)
		}

		if "example" != static_path.Mode {
			return fmt.Errorf("Bad static_path mode %s", static_path.Mode)
		}

		if "example" != static_path.PrimaryEncap {
			return fmt.Errorf("Bad static_path primary_encap %s", static_path.PrimaryEncap)
		}

		return nil
	}
}
