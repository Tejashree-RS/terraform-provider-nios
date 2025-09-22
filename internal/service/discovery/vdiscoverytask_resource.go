package discovery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	niosclient "github.com/infobloxopen/infoblox-nios-go-client/client"

	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

var readableAttributesForVdiscoverytask = "accounts_list,allow_unsecured_connection,auto_consolidate_cloud_ea,auto_consolidate_managed_tenant,auto_consolidate_managed_vm,auto_create_dns_hostname_template,auto_create_dns_record,auto_create_dns_record_type,cdiscovery_file_token,comment,credentials_type,dns_view_private_ip,dns_view_public_ip,domain_name,driver_type,enable_filter,enabled,fqdn_or_ip,govcloud_enabled,identity_version,last_run,member,merge_data,multiple_accounts_sync_policy,name,network_filter,network_list,port,private_network_view,private_network_view_mapping_policy,protocol,public_network_view,public_network_view_mapping_policy,role_arn,scheduled_run,selected_regions,service_account_file,service_account_file_token,state,state_msg,sync_child_accounts,update_dns_view_private_ip,update_dns_view_public_ip,update_metadata,use_identity,username"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VdiscoverytaskResource{}
var _ resource.ResourceWithImportState = &VdiscoverytaskResource{}

func NewVdiscoverytaskResource() resource.Resource {
	return &VdiscoverytaskResource{}
}

// VdiscoverytaskResource defines the resource implementation.
type VdiscoverytaskResource struct {
	client *niosclient.APIClient
}

func (r *VdiscoverytaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + "discovery_vdiscoverytask"
}

func (r *VdiscoverytaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Vdiscoverytask.",
		Attributes:          VdiscoverytaskResourceSchemaAttributes,
	}
}

func (r *VdiscoverytaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*niosclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *niosclient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *VdiscoverytaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//var diags diag.Diagnostics
	var data VdiscoverytaskModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Add internal ID exists in the Extensible Attributes if not already present
	// data.ExtAttrs, diags = AddInternalIDToExtAttrs(ctx, data.ExtAttrs, diags)
	// if diags.HasError() {
	// 	return
	// }
	apiRes, _, err := r.client.DiscoveryAPI.
		VdiscoverytaskAPI.
		Create(ctx).
		Vdiscoverytask(*data.Expand(ctx, &resp.Diagnostics)).
		ReturnFieldsPlus(readableAttributesForVdiscoverytask).
		ReturnAsObject(1).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Vdiscoverytask, got error: %s", err))
		return
	}

	res := apiRes.CreateVdiscoverytaskResponseAsObject.GetResult()
	//res.ExtAttrs, data.ExtAttrsAll, diags = RemoveInheritedExtAttrs(ctx, data.ExtAttrs, *res.ExtAttrs)
	// if diags.HasError() {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while create Vdiscoverytask, got error: %s", err))
	// 	return
	// }

	data.Flatten(ctx, &res, &resp.Diagnostics)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VdiscoverytaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//var diags diag.Diagnostics
	var data VdiscoverytaskModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiRes, httpRes, err := r.client.DiscoveryAPI.
		VdiscoverytaskAPI.
		Read(ctx, utils.ExtractResourceRef(data.Ref.ValueString())).
		ReturnFieldsPlus(readableAttributesForVdiscoverytask).
		ReturnAsObject(1).
		Execute()

	// If the resource is not found, try searching using Extensible Attributes
	if err != nil {
		if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Vdiscoverytask, got error: %s", err))
		return
	}

	res := apiRes.GetVdiscoverytaskResponseObjectAsResult.GetResult()

	// apiTerraformId, ok := (*res.ExtAttrs)[terraformInternalIDEA]
	// if !ok {
	// 	apiTerraformId.Value = ""
	// }

	// stateExtAttrs := ExpandExtAttrs(ctx, data.ExtAttrsAll, &diags)
	// if stateExtAttrs == nil {
	// 	resp.Diagnostics.AddError(
	// 		"Missing Internal ID",
	// 		"Unable to read Vdiscoverytask because the internal ID (from extattrs_all) is missing or invalid.",
	// 	)
	// 	return
	// }

	// stateTerraformId := (*stateExtAttrs)[terraformInternalIDEA]
	// if apiTerraformId.Value != stateTerraformId.Value {
	// 	if r.ReadByExtAttrs(ctx, &data, resp) {
	// 		return
	// 	}
	// }

	// res.ExtAttrs, data.ExtAttrsAll, diags = RemoveInheritedExtAttrs(ctx, data.ExtAttrs, *res.ExtAttrs)
	// if diags.HasError() {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while reading Vdiscoverytask due inherited Extensible attributes, got error: %s", diags))
	// 	return
	// }

	data.Flatten(ctx, &res, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VdiscoverytaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var diags diag.Diagnostics
	var data VdiscoverytaskModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//planExtAttrs := data.ExtAttrs
	diags = req.State.GetAttribute(ctx, path.Root("ref"), &data.Ref)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// diags = req.State.GetAttribute(ctx, path.Root("extattrs_all"), &data.ExtAttrsAll)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	// // Add Inherited Extensible Attributes
	// data.ExtAttrs, diags = AddInheritedExtAttrs(ctx, data.ExtAttrs, data.ExtAttrsAll)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	apiRes, _, err := r.client.DiscoveryAPI.
		VdiscoverytaskAPI.
		Update(ctx, utils.ExtractResourceRef(data.Ref.ValueString())).
		Vdiscoverytask(*data.Expand(ctx, &resp.Diagnostics)).
		ReturnFieldsPlus(readableAttributesForVdiscoverytask).
		ReturnAsObject(1).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Vdiscoverytask, got error: %s", err))
		return
	}

	res := apiRes.UpdateVdiscoverytaskResponseAsObject.GetResult()

	// res.ExtAttrs, data.ExtAttrsAll, diags = RemoveInheritedExtAttrs(ctx, planExtAttrs, *res.ExtAttrs)
	// if diags.HasError() {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while update Vdiscoverytask due inherited Extensible attributes, got error: %s", diags))
	// 	return
	// }

	data.Flatten(ctx, &res, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VdiscoverytaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VdiscoverytaskModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.DiscoveryAPI.
		VdiscoverytaskAPI.
		Delete(ctx, utils.ExtractResourceRef(data.Ref.ValueString())).
		Execute()
	if err != nil {
		if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Vdiscoverytask, got error: %s", err))
		return
	}
}

func (r *VdiscoverytaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("ref"), req, resp)
}

// func (r *VdiscoverytaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	var diags diag.Diagnostics
// 	var data VdiscoverytaskModel

// 	resourceRef := utils.ExtractResourceRef(req.ID)

// 	apiRes, _, err := r.client.DiscoveryAPI.
// 		VdiscoverytaskAPI.
// 		Read(ctx, resourceRef).
// 		ReturnFieldsPlus(readableAttributesForVdiscoverytask).
// 		ReturnAsObject(1).
// 		Execute()
// 	if err != nil {
// 		resp.Diagnostics.AddError("Import Failed", fmt.Sprintf("Cannot read Vdiscoverytask for import, got error: %s", err))
// 		return
// 	}

// 	res := apiRes.GetVdiscoverytaskResponseObjectAsResult.GetResult()

// 	res.ExtAttrs, data.ExtAttrsAll, diags = RemoveInheritedExtAttrs(ctx, data.ExtAttrs, *res.ExtAttrs)
// 	if diags.HasError() {
// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while reading Vdiscoverytask for import due inherited Extensible attributes, got error: %s", diags))
// 		return
// 	}

// 	data.Flatten(ctx, &res, &resp.Diagnostics)

// 	planExtAttrs := data.ExtAttrs
// 	data.ExtAttrs, diags = AddInheritedExtAttrs(ctx, data.ExtAttrs, data.ExtAttrsAll)
// 	if diags.HasError() {
// 		resp.Diagnostics.Append(diags...)
// 		return
// 	}

// 	data.ExtAttrs, diags = AddInternalIDToExtAttrs(ctx, data.ExtAttrs, diags)
// 	if diags.HasError() {
// 		return
// 	}

// 	updateRes, _, err := r.client.DiscoveryAPI.
// 		VdiscoverytaskAPI.
// 		Update(ctx, resourceRef).
// 		Vdiscoverytask(*data.Expand(ctx, &resp.Diagnostics)).
// 		ReturnFieldsPlus(readableAttributesForVdiscoverytask).
// 		ReturnAsObject(1).
// 		Execute()
// 	if err != nil {
// 		resp.Diagnostics.AddError("Import Failed", fmt.Sprintf("Unable to update Vdiscoverytask for import, got error: %s", err))
// 		return
// 	}

// 	res = updateRes.UpdateVdiscoverytaskResponseAsObject.GetResult()

// 	res.ExtAttrs, data.ExtAttrsAll, diags = RemoveInheritedExtAttrs(ctx, planExtAttrs, *res.ExtAttrs)
// 	if diags.HasError() {
// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while update Vdiscoverytask due inherited Extensible attributes for import, got error: %s", diags))
// 		return
// 	}
// 	data.Flatten(ctx, &res, &resp.Diagnostics)

// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }
