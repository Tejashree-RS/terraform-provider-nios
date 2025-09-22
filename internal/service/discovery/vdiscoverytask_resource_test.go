package discovery_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/infobloxopen/infoblox-nios-go-client/discovery"
	"github.com/infobloxopen/terraform-provider-nios/internal/acctest"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

var readableAttributesForVdiscoverytask = "accounts_list,allow_unsecured_connection,auto_consolidate_cloud_ea,auto_consolidate_managed_tenant,auto_consolidate_managed_vm,auto_create_dns_hostname_template,auto_create_dns_record,auto_create_dns_record_type,cdiscovery_file_token,comment,credentials_type,dns_view_private_ip,dns_view_public_ip,domain_name,driver_type,enable_filter,enabled,fqdn_or_ip,govcloud_enabled,identity_version,last_run,member,merge_data,multiple_accounts_sync_policy,name,network_filter,network_list,port,private_network_view,private_network_view_mapping_policy,protocol,public_network_view,public_network_view_mapping_policy,role_arn,scheduled_run,selected_regions,service_account_file,service_account_file_token,state,state_msg,sync_child_accounts,update_dns_view_private_ip,update_dns_view_public_ip,update_metadata,use_identity,username"

func TestAccVdiscoverytaskResource_basic(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskBasicConfig(
					name,
					true, // auto_consolidate_cloud_ea (bool)
					true, // auto_consolidate_managed_tenant (bool)
					true, // auto_consolidate_managed_vm (bool)
					"AWS",
					"infoblox.172_28_83_29", // member (string)
					"AUTO_CREATE",           // private_network_view_mapping_policy (string)
					"AUTO_CREATE",           // public_network_view_mapping_policy (string)
					true,                    // merge_data (bool)
					false,                   // update_metadata (bool) - changed to false as per AWS example
					"us-east-1",             // selected_regions (string) - changed to AWS region
					"aws_access_key",        // username (string) - changed to AWS access key
					"aws_secret_key",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "driver_type", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "member", "infoblox.172_28_83_29"),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_cloud_ea", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_managed_tenant", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_managed_vm", "true"),
					resource.TestCheckResourceAttr(resourceName, "merge_data", "true"),
					resource.TestCheckResourceAttr(resourceName, "private_network_view_mapping_policy", "AUTO_CREATE"),
					resource.TestCheckResourceAttr(resourceName, "public_network_view_mapping_policy", "AUTO_CREATE"),
					resource.TestCheckResourceAttr(resourceName, "update_metadata", "false"),
					resource.TestCheckResourceAttr(resourceName, "selected_regions", "us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "username", "aws_access_key"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_disappears(t *testing.T) {
	resourceName := "nios_discovery_vdiscoverytask.test"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVdiscoverytaskDestroy(context.Background(), &v),
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskBasicConfig(
					name,                    // name (string)
					true,                    // auto_consolidate_cloud_ea (bool)
					true,                    // auto_consolidate_managed_tenant (bool)
					true,                    // auto_consolidate_managed_vm (bool)
					"AWS",                   // driver_type (string)
					"infoblox.172_28_83_29", // member (string)
					"AUTO_CREATE",           // private_network_view_mapping_policy (string)
					"AUTO_CREATE",           // public_network_view_mapping_policy (string)
					true,                    // merge_data (bool)
					false,                   // update_metadata (bool)
					"ap-northeast-1",        // selected_regions (string)
					"aws_access_key",        // username (string)
					"aws_secret_key",        // password (string)
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					testAccCheckVdiscoverytaskDisappears(context.Background(), &v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AllowUnsecuredConnection(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_allow_unsecured_connection"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAllowUnsecuredConnection(name, true, true, true, true, true, false, "VMWARE", "vcenter.example.com", "infoblox.172_28_83_29", "vmware_password", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1", "vc_admin", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "allow_unsecured_connection", "true"),
				),
			},
			{
				Config: testAccVdiscoverytaskAllowUnsecuredConnection(name, false, true, true, true, true, false, "VMWARE", "vcenter.example.com", "infoblox.172_28_83_29", "vmware_password", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1", "vc_admin", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "allow_unsecured_connection", "false"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AutoConsolidateCloudEa(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_auto_consolidate_cloud_ea"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAutoConsolidateCloudEa(name, true, true, true, true, false, "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_cloud_ea", "true"),
				),
			},
			{
				Config: testAccVdiscoverytaskAutoConsolidateCloudEa(name, false, true, true, true, false, "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_cloud_ea", "false"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AutoConsolidateManagedTenant(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_auto_consolidate_managed_tenant"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAutoConsolidateManagedTenant(name, true, true, true, true, false, "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_managed_tenant", "true"),
				),
			},
			{
				Config: testAccVdiscoverytaskAutoConsolidateManagedTenant(name, false, true, true, true, false, "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_managed_tenant", "false"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AutoConsolidateManagedVm(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_auto_consolidate_managed_vm"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAutoConsolidateManagedVm(name, true, true, true, true, false, "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_managed_vm", "true"),
				),
			},
			{
				Config: testAccVdiscoverytaskAutoConsolidateManagedVm(name, false, true, true, true, false, "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_consolidate_managed_vm", "false"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AutoCreateDnsHostnameTemplate(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_auto_create_dns_hostname_template"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAutoCreateDnsHostnameTemplate(name, "$${vm_name}.mycompany.com", true, true, true, true, "HOST_RECORD", "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "ap-northeast-1", "aws_access_key", "aws_secret_key", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_create_dns_hostname_template", "${vm_name}.mycompany.com"),
				),
			},
			{
				Config: testAccVdiscoverytaskAutoCreateDnsHostnameTemplate(name, "$${vm_name}.updated.com", true, true, true, true, "HOST_RECORD", "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "ap-northeast-1", "aws_access_key", "aws_secret_key", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_create_dns_hostname_template", "${vm_name}.updated.com"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AutoCreateDnsRecord(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_auto_create_dns_record"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAutoCreateDnsRecord(name, true, true, true, true, "$${vm_name}.mycompany.com", "HOST_RECORD", "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "ap-northeast-1", "aws_access_key", "aws_secret_key", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_create_dns_record", "true"),
				),
			},
			{
				Config: testAccVdiscoverytaskAutoCreateDnsRecord(name, false, true, true, true, "$${vm_name}.mycompany.com", "HOST_RECORD", "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "ap-northeast-1", "aws_access_key", "aws_secret_key", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_create_dns_record", "false"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_AutoCreateDnsRecordType(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_auto_create_dns_record_type"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskAutoCreateDnsRecordType(name, "HOST_RECORD", true, true, true, true, "$${vm_name}.mycompany.com", "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "ap-northeast-1", "aws_access_key", "aws_secret_key", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_create_dns_record_type", "HOST_RECORD"),
				),
			},
			{
				Config: testAccVdiscoverytaskAutoCreateDnsRecordType(name, "A_PTR_RECORD", true, true, true, true, "$${vm_name}.mycompany.com", "AWS", "infoblox.172_28_83_29", "AUTO_CREATE", "AUTO_CREATE", "ap-northeast-1", "aws_access_key", "aws_secret_key", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_create_dns_record_type", "A_PTR_RECORD"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_CdiscoveryFileToken(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_cdiscovery_file_token"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskCdiscoveryFileToken("CDISCOVERY_FILE_TOKEN_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "cdiscovery_file_token", "CDISCOVERY_FILE_TOKEN_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskCdiscoveryFileToken("CDISCOVERY_FILE_TOKEN_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "cdiscovery_file_token", "CDISCOVERY_FILE_TOKEN_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Comment(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_comment"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskComment(name, "This is a test comment", true, true, true, "infoblox.172_28_83_29", "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "ap-northeast-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "comment", "This is a test comment"),
				),
			},
			{
				Config: testAccVdiscoverytaskComment(name, "This is an updated comment", true, true, true, "infoblox.172_28_83_29", "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "ap-northeast-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "comment", "This is an updated comment"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_CredentialsType(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_credentials_type"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskCredentialsType("CREDENTIALS_TYPE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "credentials_type", "CREDENTIALS_TYPE_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskCredentialsType("CREDENTIALS_TYPE_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "credentials_type", "CREDENTIALS_TYPE_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_DnsViewPrivateIp(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_dns_view_private_ip"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskDnsViewPrivateIp(name, "default", true, "$${vm_name}.domain.com", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "DIRECT", "default", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "dns_view_private_ip", "default"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskDnsViewPrivateIp(name, "custom_dns_view", true, "$${vm_name}.domain.com", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "DIRECT", "default", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "dns_view_private_ip", "custom_dns_view"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_DnsViewPublicIp(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_dns_view_public_ip"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskDnsViewPublicIp(name, "default", true, "$${vm_name}.domain.com", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "DIRECT", "default", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "dns_view_public_ip", "default"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskDnsViewPublicIp(name, "custom_dns_view", true, "$${vm_name}.domain.com", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "DIRECT", "default", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "dns_view_public_ip", "custom_dns_view"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_DomainName(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_domain_name"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskDomainName(name, "default", "openstack.example.com", "infoblox.172_28_83_29", "KEYSTONE_V3", "openstack_user", "openstack_password", true, true, true, true, "OPENSTACK", "AUTO_CREATE", "AUTO_CREATE", true, false, 443, "HTTPS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "default"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskDomainName(name, "custom", "openstack.example.com", "infoblox.172_28_83_29", "KEYSTONE_V3", "openstack_user", "openstack_password", true, true, true, true, "OPENSTACK", "AUTO_CREATE", "AUTO_CREATE", true, false, 443, "HTTPS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "custom"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_DriverType(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_driver_type"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create OPENSTACK and Read
			{
				Config: testAccVdiscoverytaskDriverType(name, "OPENSTACK", "openstack.example.com", "infoblox.172_28_83_29", "openstack_user", "openstack_password", true, true, true, true, "AUTO_CREATE", "AUTO_CREATE", true, false, 80, "HTTP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "driver_type", "OPENSTACK"),
				),
			},
			// Update to VMWARE and Read
			{
				Config: testAccVdiscoverytaskDriverType(name, "VMWARE", "vcenter.example.com", "infoblox.172_28_83_29", "vc_admin", "vmware_password", false, true, true, true, "AUTO_CREATE", "AUTO_CREATE", true, false, 80, "HTTP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "driver_type", "VMWARE"),
					resource.TestCheckResourceAttr(resourceName, "fqdn_or_ip", "vcenter.example.com"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_EnableFilter(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_enable_filter"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")
	networklist := []string{"10.0.0.0/8", "20.0.0.0/16"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskEnableFilter(name, true, networklist, "INCLUDE", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enable_filter", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskEnableFilter(name, false, networklist, "INCLUDE", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enable_filter", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Enabled(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_enabled"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskEnabled("ENABLED_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "ENABLED_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskEnabled("ENABLED_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "ENABLED_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_FqdnOrIp(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_fqdn_or_ip"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskFqdnOrIp(name, "vcenter.example.com", "vc_admin", "vmware_password", "infoblox.172_28_83_29", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "fqdn_or_ip", "vcenter.example.com"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskFqdnOrIp(name, "vcenter2.example.com", "vc_admin", "vmware_password", "infoblox.172_28_83_29", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "fqdn_or_ip", "vcenter2.example.com"),
				),
			},
			//Update and Read
			{
				Config: testAccVdiscoverytaskFqdnOrIp(name, "15.0.0.1", "vc_admin", "vmware_password", "infoblox.172_28_83_29", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "fqdn_or_ip", "15.0.0.1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// "Cannot change the job type(GovCloud) during update."
func TestAccVdiscoverytaskResource_GovcloudEnabled(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_govcloud_enabled"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskGovcloudEnabled(name, true, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-gov-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "govcloud_enabled", "true"),
				),
			},
			// Update and Read
			// {
			// 	Config: testAccVdiscoverytaskGovcloudEnabled(name, false, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-gov-east-1"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
			// 		resource.TestCheckResourceAttr(resourceName, "govcloud_enabled", "false"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_IdentityVersion(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_identity_version"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskIdentityVersionV2(name, "KEYSTONE_V2", "openstack.example.com", "infoblox.172_28_83_29", "openstack_user", "openstack_password", true, true, true, true, "OPENSTACK", "AUTO_CREATE", "AUTO_CREATE", true, false, 80, "HTTP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "identity_version", "KEYSTONE_V2"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskIdentityVersionV3(name, "KEYSTONE_V3", "default", "openstack.example.com", "infoblox.172_28_83_29", "openstack_user", "openstack_password", true, true, true, true, "OPENSTACK", "AUTO_CREATE", "AUTO_CREATE", true, false, 80, "HTTP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "identity_version", "KEYSTONE_V3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Member(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_member"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskMember(name, "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "member", "infoblox.172_28_83_29"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskMember(name, "infoblox.172_28_82_115", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "member", "infoblox.172_28_82_115"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_MergeData(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_merge_data"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskMergeData(name, true, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "merge_data", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskMergeData(name, false, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "merge_data", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_MultipleAccountsSyncPolicy(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_multiple_accounts_sync_policy"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskMultipleAccountsSyncPolicy(name, "DISCOVER", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_140", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, true, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "multiple_accounts_sync_policy", "DISCOVER"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskMultipleAccountsSyncPolicy(name, "UPLOAD", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_140", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, true, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "multiple_accounts_sync_policy", "UPLOAD"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Name(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_name"
	var v discovery.Vdiscoverytask

	name1 := acctest.RandomNameWithPrefix("vdiscoverytask-")
	name2 := acctest.RandomNameWithPrefix("updated-vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskName(name1, "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", name1),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskName(name2, "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1", "aws_access_key", "aws_secret_key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_NetworkFilter(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_network_filter"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")
	networkList := []string{"10.0.0.0/8", "25.0.0.0/16"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read - INCLUDE
			{
				Config: testAccVdiscoverytaskNetworkFilter(name, true, "INCLUDE", networkList, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "network_filter", "INCLUDE"),
				),
			},
			// Update and Read - EXCLUDE
			{
				Config: testAccVdiscoverytaskNetworkFilter(name, true, "EXCLUDE", networkList, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "network_filter", "EXCLUDE"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_NetworkList(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_network_list"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	// Define two different network lists for create and update scenarios
	networkList1 := []string{"10.0.0.0/8", "192.168.0.0/16"}
	networkList2 := []string{"172.16.0.0/12", "203.0.113.0/24", "198.51.100.0/24"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskNetworkList(name, networkList1, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "network_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_list.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "network_list.1", "192.168.0.0/16"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskNetworkList(name, networkList2, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "network_list.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_list.0", "172.16.0.0/12"),
					resource.TestCheckResourceAttr(resourceName, "network_list.1", "203.0.113.0/24"),
					resource.TestCheckResourceAttr(resourceName, "network_list.2", "198.51.100.0/24"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Password(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_password"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")
	password1 := "aws_secret_key1"
	password2 := "aws_secret_key2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskPassword(name, "aws_access_key", password1, "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "password", password1),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskPassword(name, "aws_access_key", password2, "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "password", password2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Port(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_port"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskPort(name, 443, "vc_admin", "vmware_password", "infoblox.172_28_83_29", "vcenter.example.com", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskPort(name, 8080, "vc_admin", "vmware_password", "infoblox.172_28_83_29", "vcenter.example.com", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "HTTPS", "AUTO_CREATE", "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_PrivateNetworkView(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_private_network_view"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskPrivateNetworkView(name, "default", "DIRECT", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "private_network_view", "default"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskPrivateNetworkView(name, "custom_private_view", "DIRECT", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "private_network_view", "custom_private_view"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccVdiscoverytaskResource_PrivateNetworkViewMappingPolicy(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_private_network_view_mapping_policy"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskPrivateNetworkViewMappingPolicyAutoCreate(name, "AUTO_CREATE", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "private_network_view_mapping_policy", "AUTO_CREATE"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskPrivateNetworkViewMappingPolicyDirect(name, "DIRECT", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "private_network_view_mapping_policy", "DIRECT"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_Protocol(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_protocol"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskProtocol(name, "HTTPS", "vc_admin", "vmware_password", "infoblox.172_28_83_29", "vcenter.example.com", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskProtocol(name, "HTTP", "vc_admin", "vmware_password", "infoblox.172_28_83_29", "vcenter.example.com", true, true, true, true, false, "VMWARE", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_PublicNetworkView(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_public_network_view"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskPublicNetworkView(name, "default", "DIRECT", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "public_network_view", "default"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskPublicNetworkView(name, "custom_public_view", "DIRECT", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "public_network_view", "custom_public_view"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_PublicNetworkViewMappingPolicy(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_public_network_view_mapping_policy"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read - AUTO_CREATE (no public_network_view)
			{
				Config: testAccVdiscoverytaskPublicNetworkViewMappingPolicyAutoCreate(name, "AUTO_CREATE", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "public_network_view_mapping_policy", "AUTO_CREATE"),
				),
			},
			// Update and Read - DIRECT (with public_network_view)
			{
				Config: testAccVdiscoverytaskPublicNetworkViewMappingPolicyDirect(name, "DIRECT", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "public_network_view_mapping_policy", "DIRECT"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_RoleArn(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_role_arn"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskRoleArn(name, "arn:aws:iam::123456789012:role/InfobloxDiscoveryRole", "DISCOVER", true, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_140", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, true, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "role_arn", "arn:aws:iam::123456789012:role/InfobloxDiscoveryRole"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskRoleArn(name, "arn:aws:iam::123456789012:role/UpdatedInfobloxRole", "DISCOVER", true, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_140", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, true, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "role_arn", "arn:aws:iam::123456789012:role/UpdatedInfobloxRole"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_ScheduledRun(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_scheduled_run"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskScheduledRun("SCHEDULED_RUN_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "scheduled_run", "SCHEDULED_RUN_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskScheduledRun("SCHEDULED_RUN_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "scheduled_run", "SCHEDULED_RUN_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_SelectedRegions(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_selected_regions"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskSelectedRegions(name, "us-east-1", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "selected_regions", "us-east-1"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskSelectedRegions(name, "us-west-1", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "selected_regions", "us-west-1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_ServiceAccountFile(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_service_account_file"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskServiceAccountFile("SERVICE_ACCOUNT_FILE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "service_account_file", "SERVICE_ACCOUNT_FILE_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskServiceAccountFile("SERVICE_ACCOUNT_FILE_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "service_account_file", "SERVICE_ACCOUNT_FILE_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_ServiceAccountFileToken(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_service_account_file_token"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskServiceAccountFileToken("SERVICE_ACCOUNT_FILE_TOKEN_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "service_account_file_token", "SERVICE_ACCOUNT_FILE_TOKEN_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskServiceAccountFileToken("SERVICE_ACCOUNT_FILE_TOKEN_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "service_account_file_token", "SERVICE_ACCOUNT_FILE_TOKEN_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_SyncChildAccounts(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_sync_child_accounts"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskSyncChildAccounts(name, true, "arn:aws:iam::123456789012:role/InfobloxDiscoveryRole", "DISCOVER", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_140", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, true, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "sync_child_accounts", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskSyncChildAccounts(name, false, "arn:aws:iam::123456789012:role/UpdatedInfobloxRole", "DISCOVER", "aws_access_key", "aws_secret_key", "infoblox.172_28_83_140", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, true, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "sync_child_accounts", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// func TestAccVdiscoverytaskResource_UpdateDnsViewPrivateIp(t *testing.T) {
// 	var resourceName = "nios_discovery_vdiscoverytask.test_update_dns_view_private_ip"
// 	var v discovery.Vdiscoverytask
// 	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

//		resource.ParallelTest(t, resource.TestCase{
//			PreCheck:                 func() { acctest.PreCheck(t) },
//			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
//			Steps: []resource.TestStep{
//				// Create and Read
//				{
//					Config: testAccVdiscoverytaskUpdateDnsViewPrivateIp(name, true, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "DIRECT", "AUTO_CREATE", true, false, "us-east-1"),
//					Check: resource.ComposeTestCheckFunc(
//						testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
//						resource.TestCheckResourceAttr(resourceName, "update_dns_view_private_ip", "true"),
//					),
//				},
//				// Update and Read
//				{
//					Config: testAccVdiscoverytaskUpdateDnsViewPrivateIp(name, false, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
//					Check: resource.ComposeTestCheckFunc(
//						testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
//						resource.TestCheckResourceAttr(resourceName, "update_dns_view_private_ip", "false"),
//					),
//				},
//				// Delete testing automatically occurs in TestCase
//			},
//		})
//	}
func TestAccVdiscoverytaskResource_UpdateDnsViewPrivateIp(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_update_dns_view_private_ip"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskUpdateDnsViewPrivateIp(
					name,
					true, // update_dns_view_private_ip
					"aws_access_key",
					"aws_secret_key",
					"infoblox.172_28_83_29",
					true, true, true,
					"AWS",
					"DIRECT", // private_network_view_mapping_policy
					"AUTO_CREATE",
					true, false,
					"us-east-1",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "update_dns_view_private_ip", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskUpdateDnsViewPrivateIp(
					name,
					false, // update_dns_view_private_ip
					"aws_access_key",
					"aws_secret_key",
					"infoblox.172_28_83_29",
					true, true, true,
					"AWS",
					"AUTO_CREATE", // private_network_view_mapping_policy
					"AUTO_CREATE",
					true, false,
					"us-east-1",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "update_dns_view_private_ip", "false"),
				),
			},
		},
	})
}

func TestAccVdiscoverytaskResource_UpdateDnsViewPublicIp(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_update_dns_view_public_ip"
	var v discovery.Vdiscoverytask

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskUpdateDnsViewPublicIp("UPDATE_DNS_VIEW_PUBLIC_IP_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "update_dns_view_public_ip", "UPDATE_DNS_VIEW_PUBLIC_IP_REPLACE_ME"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskUpdateDnsViewPublicIp("UPDATE_DNS_VIEW_PUBLIC_IP_UPDATE_REPLACE_ME"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "update_dns_view_public_ip", "UPDATE_DNS_VIEW_PUBLIC_IP_UPDATE_REPLACE_ME"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVdiscoverytaskResource_UpdateMetadata(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_update_metadata"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskUpdateMetadata(name, true, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "update_metadata", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskUpdateMetadata(name, false, "aws_access_key", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", "us-east-1", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "update_metadata", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccVdiscoverytaskResource_UseIdentity(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_use_identity"
	var v discovery.Vdiscoverytask
	name := acctest.RandomNameWithPrefix("vdiscoverytask-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVdiscoverytaskUseIdentity(name, true, "openstack.example.com", 80, "HTTP", "infoblox.172_28_83_29", "KEYSTONE_V2", "openstack_user", "openstack_password", true, true, true, "OPENSTACK", "AUTO_CREATE", "AUTO_CREATE", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "use_identity", "true"),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskUseIdentity(name, false, "openstack.example.com", 80, "HTTP", "infoblox.172_28_83_29", "KEYSTONE_V2", "openstack_user", "openstack_password", true, true, true, "OPENSTACK", "AUTO_CREATE", "AUTO_CREATE", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "use_identity", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccVdiscoverytaskResource_Username(t *testing.T) {
	var resourceName = "nios_discovery_vdiscoverytask.test_username"
	var v discovery.Vdiscoverytask

	name := acctest.RandomNameWithPrefix("vdiscoverytask-")
	username1 := "User1"
	username2 := "User2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVdiscoverytaskUsername(name, "USERNAME_REPLACE_ME", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "username", "USERNAME_REPLACE_ME"),
				),
			},
			// Create and Read
			{
				Config: testAccVdiscoverytaskUsername(name, "User1", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "username", username1),
				),
			},
			// Update and Read
			{
				Config: testAccVdiscoverytaskUsername(name, "User2", "aws_secret_key", "infoblox.172_28_83_29", true, true, true, "AWS", "AUTO_CREATE", "AUTO_CREATE", true, false, "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdiscoverytaskExists(context.Background(), resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "username", username2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckVdiscoverytaskExists(ctx context.Context, resourceName string, v *discovery.Vdiscoverytask) resource.TestCheckFunc {
	// Verify the resource exists in the cloud
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		apiRes, _, err := acctest.NIOSClient.DiscoveryAPI.
			VdiscoverytaskAPI.
			Read(ctx, utils.ExtractResourceRef(rs.Primary.Attributes["ref"])).
			ReturnFieldsPlus(readableAttributesForVdiscoverytask).
			ReturnAsObject(1).
			Execute()
		if err != nil {
			return err
		}
		if !apiRes.GetVdiscoverytaskResponseObjectAsResult.HasResult() {
			return fmt.Errorf("expected result to be returned: %s", resourceName)
		}
		*v = apiRes.GetVdiscoverytaskResponseObjectAsResult.GetResult()
		return nil
	}
}

func testAccCheckVdiscoverytaskDestroy(ctx context.Context, v *discovery.Vdiscoverytask) resource.TestCheckFunc {
	// Verify the resource was destroyed
	return func(state *terraform.State) error {
		_, httpRes, err := acctest.NIOSClient.DiscoveryAPI.
			VdiscoverytaskAPI.
			Read(ctx, utils.ExtractResourceRef(*v.Ref)).
			ReturnAsObject(1).
			ReturnFieldsPlus(readableAttributesForVdiscoverytask).
			Execute()
		if err != nil {
			if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
				// resource was deleted
				return nil
			}
			return err
		}
		return errors.New("expected to be deleted")
	}
}

func testAccCheckVdiscoverytaskDisappears(ctx context.Context, v *discovery.Vdiscoverytask) resource.TestCheckFunc {
	// Delete the resource externally to verify disappears test
	return func(state *terraform.State) error {
		_, err := acctest.NIOSClient.DiscoveryAPI.
			VdiscoverytaskAPI.
			Delete(ctx, utils.ExtractResourceRef(*v.Ref)).
			Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

// With Respect to AWS
func testAccVdiscoverytaskBasicConfig(name string, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy string, merge_data, update_metadata bool, selected_regions, username, password string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test" {
    name = %q
    driver_type = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    update_metadata = %t
    selected_regions = %q
    username = %q
    password = %q
}
`, name, driver_type, member, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, private_network_view_mapping_policy, public_network_view_mapping_policy, update_metadata, selected_regions, username, password)
}

func testAccVdiscoverytaskAllowUnsecuredConnection(
	name string,
	allow_unsecured_connection bool,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata bool,
	driver_type, fqdn_or_ip, member, password, private_network_view_mapping_policy, protocol, public_network_view_mapping_policy, selected_regions, username string,
	port int,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_allow_unsecured_connection" {
    name = %q
    allow_unsecured_connection = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    fqdn_or_ip = %q
    member = %q
    password = %q
    private_network_view_mapping_policy = %q
    protocol = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    port = %d
}
`, name, allow_unsecured_connection, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata, driver_type, fqdn_or_ip, member, password, private_network_view_mapping_policy, protocol, public_network_view_mapping_policy, selected_regions, username, port)
}

func testAccVdiscoverytaskAutoConsolidateCloudEa(
	name string,
	auto_consolidate_cloud_ea bool,
	auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata bool,
	driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_auto_consolidate_cloud_ea" {
    name = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    member = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    password = %q
}
`, name, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password)
}

func testAccVdiscoverytaskAutoConsolidateManagedTenant(
	name string,
	auto_consolidate_managed_tenant bool,
	auto_consolidate_cloud_ea, auto_consolidate_managed_vm, merge_data, update_metadata bool,
	driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_auto_consolidate_managed_tenant" {
    name = %q
    auto_consolidate_managed_tenant = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    member = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    password = %q
}
`, name, auto_consolidate_managed_tenant, auto_consolidate_cloud_ea, auto_consolidate_managed_vm, merge_data, update_metadata, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password)
}

func testAccVdiscoverytaskAutoConsolidateManagedVm(
	name string,
	auto_consolidate_managed_vm bool,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, merge_data, update_metadata bool,
	driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_auto_consolidate_managed_vm" {
    name = %q
    auto_consolidate_managed_vm = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    member = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    password = %q
}
`, name, auto_consolidate_managed_vm, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, merge_data, update_metadata, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password)
}

func testAccVdiscoverytaskAutoCreateDnsHostnameTemplate(
	name string,
	auto_create_dns_hostname_template string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, auto_create_dns_record bool,
	auto_create_dns_record_type, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password string,
	merge_data, update_metadata bool,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_auto_create_dns_hostname_template" {
    name = %q
    auto_create_dns_hostname_template = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    auto_create_dns_record = %t
    auto_create_dns_record_type = %q
    driver_type = %q
    member = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    password = %q
    merge_data = %t
    update_metadata = %t
}
`, name, auto_create_dns_hostname_template, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, auto_create_dns_record, auto_create_dns_record_type, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password, merge_data, update_metadata)
}

func testAccVdiscoverytaskAutoCreateDnsRecord(
	name string,
	auto_create_dns_record bool,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	auto_create_dns_hostname_template, auto_create_dns_record_type, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password string,
	merge_data, update_metadata bool,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_auto_create_dns_record" {
    name = %q
    auto_create_dns_record = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    auto_create_dns_hostname_template = %q
    auto_create_dns_record_type = %q
    driver_type = %q
    member = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    password = %q
    merge_data = %t
    update_metadata = %t
}
`, name, auto_create_dns_record, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, auto_create_dns_hostname_template, auto_create_dns_record_type, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password, merge_data, update_metadata)
}

func testAccVdiscoverytaskAutoCreateDnsRecordType(
	name string,
	auto_create_dns_record_type string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, auto_create_dns_record bool,
	auto_create_dns_hostname_template, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password string,
	merge_data, update_metadata bool,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_auto_create_dns_record_type" {
    name = %q
    auto_create_dns_record_type = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    auto_create_dns_record = %t
    auto_create_dns_hostname_template = %q
    driver_type = %q
    member = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    username = %q
    password = %q
    merge_data = %t
    update_metadata = %t
}
`, name, auto_create_dns_record_type, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, auto_create_dns_record, auto_create_dns_hostname_template, driver_type, member, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions, username, password, merge_data, update_metadata)
}

func testAccVdiscoverytaskCdiscoveryFileToken(cdiscoveryFileToken string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_cdiscovery_file_token" {
    cdiscovery_file_token = %q
}
`, cdiscoveryFileToken)
}

func testAccVdiscoverytaskComment(name, comment string, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool, member, driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string, merge_data, update_metadata bool, selected_regions, username, password string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_comment" {
	name = %q
	comment = %q
    driver_type = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    update_metadata = %t
    selected_regions = %q
    username = %q
    password = %q
}
`, name, comment, driver_type, member, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, private_network_view_mapping_policy, public_network_view_mapping_policy, update_metadata, selected_regions, username, password)
}

func testAccVdiscoverytaskCredentialsType(credentialsType string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_credentials_type" {
    credentials_type = %q
}
`, credentialsType)
}

func testAccVdiscoverytaskDnsViewPrivateIp(name, dns_view_private_ip string, update_dns_view_private_ip bool, auto_create_dns_hostname_template, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, private_network_view, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_dns_view_private_ip" {
    name = %q
    dns_view_private_ip = %q
    update_dns_view_private_ip = %t
    auto_create_dns_hostname_template = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
	private_network_view = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, dns_view_private_ip, update_dns_view_private_ip, auto_create_dns_hostname_template, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, private_network_view, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskDnsViewPublicIp(name, dns_view_public_ip string, update_dns_view_public_ip bool, auto_create_dns_hostname_template, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, public_network_view string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_dns_view_public_ip" {
    name = %q
    dns_view_public_ip = %q
    update_dns_view_public_ip = %t
    auto_create_dns_hostname_template = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
	public_network_view = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, dns_view_public_ip, update_dns_view_public_ip, auto_create_dns_hostname_template, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, public_network_view,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskDomainName(name, domain_name, fqdn_or_ip, member, identity_version, username, password string,
	useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool, port int, protocol string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_domain_name" {
    name = %q
    domain_name = %q
    fqdn_or_ip = %q
    member = %q
    identity_version = %q
    username = %q
    password = %q
    use_identity = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    port = %d
    protocol = %q
}
`, name, domain_name, fqdn_or_ip, member, identity_version, username, password, useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, merge_data, update_metadata, port, protocol)
}

func testAccVdiscoverytaskDriverType(name, driver_type, fqdn_or_ip, member, username, password string,
	useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool, port int, protocol string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_driver_type" {
    name = %q
    driver_type = %q
    fqdn_or_ip = %q
    member = %q
    identity_version = "KEYSTONE_V2"
    username = %q
    password = %q
    use_identity = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    port = %d
    protocol = %q
}
`, name, driver_type, fqdn_or_ip, member, username, password, useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, private_network_view_mapping_policy, public_network_view_mapping_policy, merge_data, update_metadata, port, protocol)
}

func testAccVdiscoverytaskEnableFilter(name string, enable_filter bool, network_list []string, network_filter string, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	networkListHCL := utils.ConvertStringSliceToHCL(network_list)
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_enable_filter" {
    name = %q
    enable_filter = %t
    network_list = %s
    network_filter = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, enable_filter, networkListHCL, network_filter, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskEnabled(enabled string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_enabled" {
    enabled = %q
}
`, enabled)
}

func testAccVdiscoverytaskFqdnOrIp(
	name, fqdn_or_ip, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata bool,
	driver_type, private_network_view_mapping_policy, protocol, public_network_view_mapping_policy, selected_regions string,
	port int,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_fqdn_or_ip" {
    name = %q
    fqdn_or_ip = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    protocol = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    port = %d
}
`, name, fqdn_or_ip, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata,
		driver_type, private_network_view_mapping_policy, protocol, public_network_view_mapping_policy, selected_regions,
		port)
}

func testAccVdiscoverytaskGovcloudEnabled(name string, govcloud_enabled bool, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_govcloud_enabled" {
    name = %q
    govcloud_enabled = %t
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, govcloud_enabled, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskIdentityVersionV2(name, identity_version, fqdn_or_ip, member, username, password string,
	useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool, port int, protocol string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_identity_version" {
    name = %q
    identity_version = %q
    fqdn_or_ip = %q
    member = %q
    username = %q
    password = %q
    use_identity = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    port = %d
    protocol = %q
	domain_name = "default"
}
`, name, identity_version, fqdn_or_ip, member, username, password, useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, merge_data, update_metadata, port, protocol)
}

func testAccVdiscoverytaskIdentityVersionV3(name, identity_version, domain_name, fqdn_or_ip, member, username, password string,
	useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool, port int, protocol string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_identity_version" {
    name = %q
    identity_version = %q
    domain_name = %q
    fqdn_or_ip = %q
    member = %q
    username = %q
    password = %q
    use_identity = %t
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    port = %d
    protocol = %q
}
`, name, identity_version, domain_name, fqdn_or_ip, member, username, password, useIdentity, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, merge_data, update_metadata, port, protocol)
}

func testAccVdiscoverytaskMember(
	name, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions, username, password string,
) string {
	return fmt.Sprintf(`
	resource "nios_discovery_vdiscoverytask" "test_member" {
    name = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    merge_data = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    update_metadata = %t
    username = %q
    password = %q
}
`, name, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, merge_data,
		private_network_view_mapping_policy, public_network_view_mapping_policy,
		selected_regions, update_metadata,
		username, password)
}

func testAccVdiscoverytaskMergeData(name string, merge_data bool, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	selected_regions string, update_metadata bool) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_merge_data" {
    name = %q
	merge_data = %t
	username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    update_metadata = %t
}
`, name, merge_data, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		selected_regions, update_metadata)
}

func testAccVdiscoverytaskMultipleAccountsSyncPolicy(name, multiple_accounts_sync_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_multiple_accounts_sync_policy" {
    name = %q
    multiple_accounts_sync_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, multiple_accounts_sync_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskName(
	name, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions, username, password string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_name" {
    name = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    merge_data = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    update_metadata = %t
    username = %q
    password = %q
}
`, name, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, merge_data,
		private_network_view_mapping_policy, public_network_view_mapping_policy,
		selected_regions, update_metadata,
		username, password)
}

func testAccVdiscoverytaskNetworkFilter(name string, enable_filter bool, network_filter string, network_list []string, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	networkListHCL := utils.ConvertStringSliceToHCL(network_list)
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_network_filter" {
    name = %q
    enable_filter = %t
    network_filter = %q
    network_list = %s
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, enable_filter, network_filter, networkListHCL, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskNetworkList(name string, network_list []string, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	networkListHCL := utils.ConvertStringSliceToHCL(network_list)
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_network_list" {
    name = %q
    network_list = %s
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, networkListHCL, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskPassword(name, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_password" {
    name = %q
	username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    merge_data = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    update_metadata = %t
}
`, name, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, merge_data,
		private_network_view_mapping_policy, public_network_view_mapping_policy,
		selected_regions, update_metadata)
}

func testAccVdiscoverytaskPort(
	name string, port int, username, password, member, fqdn_or_ip string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata bool,
	driver_type, private_network_view_mapping_policy, protocol, public_network_view_mapping_policy, selected_regions string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_port" {
    name = %q
	port = %d
    username = %q
    password = %q
    member = %q
    fqdn_or_ip = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    protocol = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
}
`, name, port, username, password, member, fqdn_or_ip,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata,
		driver_type, private_network_view_mapping_policy, protocol, public_network_view_mapping_policy, selected_regions)
}

func testAccVdiscoverytaskPrivateNetworkView(name, private_network_view, private_network_view_mapping_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_private_network_view" {
    name = %q
    private_network_view = %q
    private_network_view_mapping_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, private_network_view, private_network_view_mapping_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

// Helper for AUTO_CREATE policy (without private_network_view)
func testAccVdiscoverytaskPrivateNetworkViewMappingPolicyAutoCreate(name, private_network_view_mapping_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_private_network_view_mapping_policy" {
    name = %q
    private_network_view_mapping_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, private_network_view_mapping_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

// Helper for DIRECT policy (with private_network_view)
func testAccVdiscoverytaskPrivateNetworkViewMappingPolicyDirect(name, private_network_view_mapping_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_private_network_view_mapping_policy" {
    name = %q
    private_network_view_mapping_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
    private_network_view = "default"
}
`, name, private_network_view_mapping_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskProtocol(
	name, protocol, username, password, member, fqdn_or_ip string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions string,
	port int,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_protocol" {
    name = %q
    protocol = %q
    username = %q
    password = %q
    member = %q
    fqdn_or_ip = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    merge_data = %t
    update_metadata = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    port = %d
}
`, name, protocol, username, password, member, fqdn_or_ip,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, merge_data, update_metadata,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy, selected_regions,
		port)
}

func testAccVdiscoverytaskPublicNetworkView(name, public_network_view, public_network_view_mapping_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_public_network_view" {
    name = %q
    public_network_view = %q
    public_network_view_mapping_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, public_network_view, public_network_view_mapping_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

// func testAccVdiscoverytaskPublicNetworkViewMappingPolicy(name, public_network_view_mapping_policy, public_network_view, username, password, member string,
// 	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
// 	driver_type, private_network_view_mapping_policy string,
// 	merge_data, update_metadata bool,
// 	selected_regions string) string {
// 	return fmt.Sprintf(`
// resource "nios_discovery_vdiscoverytask" "test_public_network_view_mapping_policy" {
//     name = %q
//     public_network_view_mapping_policy = %q
// 	public_network_view = %q
//     username = %q
//     password = %q
//     member = %q
//     auto_consolidate_cloud_ea = %t
//     auto_consolidate_managed_tenant = %t
//     auto_consolidate_managed_vm = %t
//     driver_type = %q
//     private_network_view_mapping_policy = %q
//     merge_data = %t
//     update_metadata = %t
//     selected_regions = %q
//     public_network_view = "custom_public_view"

// }
// `, name, public_network_view_mapping_policy, public_network_view, username, password, member,
// 		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
// 		driver_type, private_network_view_mapping_policy,
// 		merge_data, update_metadata,
// 		selected_regions)
// }

// Helper for AUTO_CREATE policy (without public_network_view)
func testAccVdiscoverytaskPublicNetworkViewMappingPolicyAutoCreate(name, public_network_view_mapping_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_public_network_view_mapping_policy" {
    name = %q
    public_network_view_mapping_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, public_network_view_mapping_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

// Helper for DIRECT policy (with public_network_view)
func testAccVdiscoverytaskPublicNetworkViewMappingPolicyDirect(name, public_network_view_mapping_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_public_network_view_mapping_policy" {
    name = %q
    public_network_view_mapping_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
    public_network_view = "custom_public_view"
}
`, name, public_network_view_mapping_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskRoleArn(name, role_arn, multiple_accounts_sync_policy string, sync_child_accounts bool, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_role_arn" {
    name = %q
    role_arn = %q
    multiple_accounts_sync_policy = %q
    sync_child_accounts = %t
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, role_arn, multiple_accounts_sync_policy, sync_child_accounts, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}

func testAccVdiscoverytaskScheduledRun(scheduledRun string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_scheduled_run" {
    scheduled_run = %q
}
`, scheduledRun)
}

func testAccVdiscoverytaskSelectedRegions(name, selected_regions, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_selected_regions" {
    name = %q
    selected_regions = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
}
`, name, selected_regions, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata)
}

func testAccVdiscoverytaskServiceAccountFile(serviceAccountFile string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_service_account_file" {
    service_account_file = %q
}
`, serviceAccountFile)
}

func testAccVdiscoverytaskServiceAccountFileToken(serviceAccountFileToken string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_service_account_file_token" {
    service_account_file_token = %q
}
`, serviceAccountFileToken)
}

func testAccVdiscoverytaskSyncChildAccounts(name string, sync_child_accounts bool, role_arn, multiple_accounts_sync_policy, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_sync_child_accounts" {
    name = %q
    sync_child_accounts = %t
    role_arn = %q
    multiple_accounts_sync_policy = %q
    username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    merge_data = %t
    update_metadata = %t
    selected_regions = %q
}
`, name, sync_child_accounts, role_arn, multiple_accounts_sync_policy, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		merge_data, update_metadata,
		selected_regions)
}
func testAccVdiscoverytaskUpdateDnsViewPrivateIp(
	name string,
	update_dns_view_private_ip bool,
	username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string,
) string {

	var dnsViewBlock, networkViewBlock string

	// Only include these fields when DNS update is enabled AND mapping is DIRECT
	if update_dns_view_private_ip && private_network_view_mapping_policy == "DIRECT" {
		dnsViewBlock = `  dns_view_private_ip = "default"`
		networkViewBlock = `  private_network_view = "default"`
	}

	config := fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_update_dns_view_private_ip" {
  name = %q
  update_dns_view_private_ip = %t
%s
  username = %q
  password = %q
  member = %q
  auto_consolidate_cloud_ea = %t
  auto_consolidate_managed_tenant = %t
  auto_consolidate_managed_vm = %t
  driver_type = %q
  private_network_view_mapping_policy = %q
%s
  public_network_view_mapping_policy = %q
  merge_data = %t
  update_metadata = %t
  selected_regions = %q
}
`, name, update_dns_view_private_ip, dnsViewBlock, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, networkViewBlock,
		public_network_view_mapping_policy, merge_data, update_metadata,
		selected_regions)

	// Debug print to verify generated config
	fmt.Println("Generated Terraform config:\n", config)

	return config
}

// func testAccVdiscoverytaskUpdateDnsViewPrivateIp(name string, update_dns_view_private_ip bool, username, password, member string,
// 	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
// 	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
// 	merge_data, update_metadata bool,
// 	selected_regions string) string {
// 	return fmt.Sprintf(`
// resource "nios_discovery_vdiscoverytask" "test_update_dns_view_private_ip" {
//     name = %q
//     update_dns_view_private_ip = %t
// 	dns_view_private_ip = "default"
//     username = %q
//     password = %q
//     member = %q
//     auto_consolidate_cloud_ea = %t
//     auto_consolidate_managed_tenant = %t
//     auto_consolidate_managed_vm = %t
//     driver_type = %q
//     private_network_view_mapping_policy = %q
// 	private_network_view = "default"
//     public_network_view_mapping_policy = %q
//     merge_data = %t
//     update_metadata = %t
//     selected_regions = %q
// }
// `, name, update_dns_view_private_ip, username, password, member,
// 		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
// 		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
// 		merge_data, update_metadata,
// 		selected_regions)
// }

func testAccVdiscoverytaskUpdateDnsViewPublicIp(updateDnsViewPublicIp string) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_update_dns_view_public_ip" {
    update_dns_view_public_ip = %q
}
`, updateDnsViewPublicIp)
}

func testAccVdiscoverytaskUpdateMetadata(name string, update_metadata bool, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	selected_regions string, merge_data bool) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_update_metadata" {
    name = %q
	update_metadata = %t
	username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    merge_data = %t
}
`, name, update_metadata, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy,
		selected_regions, merge_data)
}

func testAccVdiscoverytaskUseIdentity(name string, useIdentity bool, fqdn_or_ip string, port int, protocol, member, identity_version, username, password string, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool, driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string, merge_data, update_metadata bool) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_use_identity" {
	name = %q
    use_identity = %t
    fqdn_or_ip = %q
    port = %d
    protocol = %q
	member = %q
    identity_version = %q
    username = %q
    password = %q
	auto_consolidate_cloud_ea = %t
	auto_consolidate_managed_tenant = %t
	auto_consolidate_managed_vm = %t
	driver_type = %q
	merge_data = %t
	private_network_view_mapping_policy = %q
	public_network_view_mapping_policy = %q
	update_metadata = %t
}
`, name, useIdentity, fqdn_or_ip, port, protocol, member, identity_version, username, password, auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm, driver_type, merge_data, private_network_view_mapping_policy, public_network_view_mapping_policy, update_metadata)
}

func testAccVdiscoverytaskUsername(
	name, username, password, member string,
	auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm bool,
	driver_type, private_network_view_mapping_policy, public_network_view_mapping_policy string,
	merge_data, update_metadata bool,
	selected_regions string,
) string {
	return fmt.Sprintf(`
resource "nios_discovery_vdiscoverytask" "test_username" {
    name = %q
	username = %q
    password = %q
    member = %q
    auto_consolidate_cloud_ea = %t
    auto_consolidate_managed_tenant = %t
    auto_consolidate_managed_vm = %t
    driver_type = %q
    merge_data = %t
    private_network_view_mapping_policy = %q
    public_network_view_mapping_policy = %q
    selected_regions = %q
    update_metadata = %t
}
`, name, username, password, member,
		auto_consolidate_cloud_ea, auto_consolidate_managed_tenant, auto_consolidate_managed_vm,
		driver_type, merge_data,
		private_network_view_mapping_policy, public_network_view_mapping_policy,
		selected_regions, update_metadata,
	)
}
