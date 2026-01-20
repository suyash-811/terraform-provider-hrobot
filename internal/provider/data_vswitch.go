// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-hrobot/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &vSwitchDataSource{}
var _ datasource.DataSourceWithConfigure = &vSwitchDataSource{}

func NewVSwitchDataSource() datasource.DataSource {
	return &vSwitchDataSource{}
}

type vSwitchDataSource struct {
	client *client.HetznerRobotClient
}

func (d *vSwitchDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vswitch"
}

func (d *vSwitchDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches vswitch information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "vSwitch ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "vSwitch name.",
				Computed:    true,
			},
			"vlan": schema.Int64Attribute{
				Description: "VLAN assigned to vSwitch.",
				Computed:    true,
			},
			"cancelled": schema.BoolAttribute{
				Description: "vSwitch cancellation status.",
				Computed:    true,
			},
			"server": schema.ListNestedAttribute{
				Description: "List of servers connected to vSwitch.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"server_ip": schema.StringAttribute{
							Description: "IPv4 of server.",
							Computed:    true,
						},
						"server_ipv6_net": schema.StringAttribute{
							Description: "IPv6 of server.",
							Computed:    true,
						},
						"server_number": schema.Int64Attribute{
							Description: "Server ID",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Server status",
							Computed:    true,
						},
					},
				},
			},
			"subnet": schema.ListNestedAttribute{
				Description: "List of additional subnets on vSwitch.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Description: "IP on subnet",
							Computed:    true,
						},
						"mask": schema.StringAttribute{
							Description: "Subnet mask",
							Computed:    true,
						},
						"gateway": schema.StringAttribute{
							Description: "Subnet gateway",
							Computed:    true,
						},
					},
				},
			},
			"cloud_network": schema.ListNestedAttribute{
				Description: "List of cloud networks attached to vSwitch",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Cloud network ID",
							Computed:    true,
						},
						"ip": schema.StringAttribute{
							Description: "vSwitch subnet address on cloud network.",
							Computed:    true,
						},
						"mask": schema.Int64Attribute{
							Description: "vSwitch subnet mask on cloud network.",
							Computed:    true,
						},
						"gateway": schema.StringAttribute{
							Description: "Gateway IP on vSwitch subnet.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *vSwitchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var vSwitchID types.Int64

	diags := req.Config.GetAttribute(ctx, path.Root("id"), &vSwitchID)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vSwitch, err := d.client.GetVSwitch(ctx, vSwitchID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read vswitch with id %d, got error: %s", vSwitchID.ValueInt64(), err),
		)
		return
	}

	state := vSwitchDataSourceModel{
		ID:        types.Int64Value(vSwitch.ID),
		Name:      types.StringValue(vSwitch.Name),
		VLAN:      types.Int64Value(vSwitch.VLAN),
		Cancelled: types.BoolValue(vSwitch.Cancelled),
	}

	for _, server := range vSwitch.Server {
		serverState := vSwitchServerModel{
			ServerIP:      types.StringValue(server.ServerIP),
			ServerIPv6Net: types.StringValue(server.ServerIPv6Net),
			ServerNumber:  types.Int64Value(server.ServerNumber),
			Status:        types.StringValue(server.Status),
		}
		state.Server = append(state.Server, serverState)
	}

	for _, subnet := range vSwitch.Subnet {
		subnetState := vSwitchSubnetModel{
			IP:      types.StringValue(subnet.IP),
			Mask:    types.Int64Value(subnet.Mask),
			Gateway: types.StringValue(subnet.Gateway),
		}
		state.Subnet = append(state.Subnet, subnetState)
	}

	for _, cloudNet := range vSwitch.CloudNetwork {
		cloudNetState := vSwitchCloudNetworkModel{
			ID:      types.Int64Value(cloudNet.ID),
			IP:      types.StringValue(cloudNet.IP),
			Mask:    types.Int64Value(cloudNet.Mask),
			Gateway: types.StringValue(cloudNet.Gateway),
		}
		state.CloudNetwork = append(state.CloudNetwork, cloudNetState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (d *vSwitchDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.HetznerRobotClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

type vSwitchDataSourceModel struct {
	ID           types.Int64                `tfsdk:"id"`
	Name         types.String               `tfsdk:"name"`
	VLAN         types.Int64                `tfsdk:"vlan"`
	Cancelled    types.Bool                 `tfsdk:"cancelled"`
	Server       []vSwitchServerModel       `tfsdk:"server"`
	Subnet       []vSwitchSubnetModel       `tfsdk:"subnet"`
	CloudNetwork []vSwitchCloudNetworkModel `tfsdk:"cloud_network"`
}

type vSwitchServerModel struct {
	ServerIP      types.String `tfsdk:"server_ip"`
	ServerIPv6Net types.String `tfsdk:"server_ipv6_net"`
	ServerNumber  types.Int64  `tfsdk:"server_number"`
	Status        types.String `tfsdk:"status"`
}

type vSwitchSubnetModel struct {
	IP      types.String `tfsdk:"ip"`
	Mask    types.Int64  `tfsdk:"mask"`
	Gateway types.String `tfsdk:"gateway"`
}

type vSwitchCloudNetworkModel struct {
	ID      types.Int64  `tfsdk:"id"`
	IP      types.String `tfsdk:"ip"`
	Mask    types.Int64  `tfsdk:"mask"`
	Gateway types.String `tfsdk:"gateway"`
}
