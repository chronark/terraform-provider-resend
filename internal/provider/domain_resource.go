// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/resendlabs/resend-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DomainResource{}
var _ resource.ResourceWithImportState = &DomainResource{}

func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// DomainResource defines the resource implementation.
type DomainResource struct {
	client *resend.Client
}

type Record struct {
	Record   types.String `tfsdk:"record"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Value    types.String `tfsdk:"value"`
	Ttl      types.String `tfsdk:"ttl"`
	Status   types.String `tfsdk:"status"`
	Priority types.Int64  `tfsdk:"priority"`
}

// DomainResourceModel describes the resource data model.
type DomainResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Region      types.String `tfsdk:"region"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Status      types.String `tfsdk:"status"`
	DnsProvider types.String `tfsdk:"dns_provider"`
	// Records     basetypes.ListValue `tfsdk:"records"`
}

func (r *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Add a new Domain.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the domain within Resend.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the domain you want to create",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region where emails will be sent from. Possible values: `us-east-1` | `eu-west-1` | `sa-east-1`",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the domain was created",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the domain. TODO: find out possible values",
				Computed:            true,
			},
			"dns_provider": schema.StringAttribute{
				MarkdownDescription: "The DNS provider used to configure the domain.",
				Computed:            true,
			},

			// "records": schema.ListNestedAttribute{
			// 		MarkdownDescription: "The DNS records used to configure the domain.",
			// Computed: true,
			// Optional: true,
			// 		NestedObject: schema.NestedAttributeObject{
			// 			Attributes: map[string]schema.Attribute{
			// 				"record": schema.StringAttribute{
			// 					MarkdownDescription: "The Record Type.",
			// 					Computed:            true,
			// 				},
			// 				"name": schema.StringAttribute{
			// 					MarkdownDescription: "The name of the record.",
			// 					Computed:            true,
			// 				},
			// 				"type": schema.StringAttribute{
			// 					MarkdownDescription: "The type of the record.",
			// 					Computed:            true,
			// 				},
			// 				"ttl": schema.StringAttribute{
			// 					MarkdownDescription: "The TTL of the record.",
			// 					Computed:            true,
			// 				},
			// 				"status": schema.StringAttribute{
			// 					MarkdownDescription: "The status of the record.",
			// 					Computed:            true,
			// 				},
			// 				"value": schema.StringAttribute{
			// 					MarkdownDescription: "The value of the record.",
			// 					Computed:            true,
			// 				},
			// 				"priority": schema.NumberAttribute{
			// 					MarkdownDescription: "The priority of the record.",
			// 					Computed:            true,
			// 					Optional:            true,
			// 				},
			// 			},
			// 		},
			// },
		},
	}
}

func (r *DomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*resend.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *resend.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.client.Domains.Create(&resend.CreateDomainRequest{
		Name:   data.Name.ValueString(),
		Region: data.Region.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create domain, got error: %s", err))
		return
	}
	data.Id = types.StringValue(domain.Id)
	data.CreatedAt = types.StringValue(domain.CreatedAt)
	data.Status = types.StringValue(domain.Status)
	data.DnsProvider = types.StringValue(domain.DnsProvider)
	data.Region = types.StringValue(domain.Region)
	data.Status = types.StringValue(domain.Status)
	// data.Records = basetypes.ListValue{}

	// for i, r := range domain.Records {
	// 	data.Records[i] = record{
	// 		Record: types.StringValue(r.Record),
	// 		Name:   types.StringValue(r.Name),
	// 		Type:   types.StringValue(r.Type),
	// 		Ttl:    types.StringValue(r.Ttl),
	// 		Status: types.StringValue(r.Status),
	// 		Value:  types.StringValue(r.Value),
	// 	}
	// 	p, err := r.Priority.Int64()
	// 	if err == nil {
	// 		// I guess we can ignore this?
	// 		data.Records[i].Priority = types.Int64Value(p)
	// 	}
	// }

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.client.Domains.Get(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read domain, got error: %s", err))
		return
	}
	data.Name = types.StringValue(domain.Name)
	data.Region = types.StringValue(domain.Region)
	data.CreatedAt = types.StringValue(domain.CreatedAt)
	data.Status = types.StringValue(domain.Status)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update domain, got error: %s", "not implemented"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Domains.Remove(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete domain, got error: %s", err))
		return
	}

}

func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
