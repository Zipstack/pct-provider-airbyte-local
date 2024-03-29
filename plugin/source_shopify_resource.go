package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-local/api"
)

// Resource implementation.
type sourceShopifyResource struct {
	Client *api.Client
}

type sourceShopifyResourceModel struct {
	Name                    string                       `pctsdk:"name"`
	SourceId                string                       `pctsdk:"source_id"`
	SourceDefinitionId      string                       `pctsdk:"source_definition_id"`
	WorkspaceId             string                       `pctsdk:"workspace_id"`
	ConnectionConfiguration sourceShopifyConnConfigModel `pctsdk:"connection_configuration"`
}

type sourceShopifyConnConfigModel struct {
	StartDate   string                 `pctsdk:"start_date"`
	Shop        string                 `pctsdk:"shop"`
	Credentials shopifyCredConfigModel `pctsdk:"credentials"`
}

type shopifyCredConfigModel struct {
	AuthMethod   string `pctsdk:"auth_method"`
	ApiPassword  string `pctsdk:"api_password"`
	ClientSecret string `pctsdk:"client_secret"`
	AccessToken  string `pctsdk:"access_token"`
	ClientId     string `pctsdk:"client_id"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceShopifyResource{}
)

// Helper function to return a resource service instance.
func NewSourceShopifyResource() schema.ResourceService {
	return &sourceShopifyResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceShopifyResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_shopify",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceShopifyResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
	if req.ResourceData == "" {
		return schema.ErrorResponse(fmt.Errorf("no data provided to configure resource"))
	}

	var creds map[string]string
	err := fwhelpers.Decode(req.ResourceData, &creds)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	client, err := api.NewClient(
		creds["host"], creds["username"], creds["password"],
	)
	if err != nil {
		return schema.ErrorResponse(fmt.Errorf("malformed data provided to configure resource"))
	}

	r.Client = client

	return &schema.ServiceResponse{}
}

// Schema defines the schema for the resource.
func (r *sourceShopifyResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source Shopify resource for Airbyte",
		Attributes: map[string]schema.Attribute{
			"name": &schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			"source_id": &schema.StringAttribute{
				Description: "Source ID",
				Computed:    true,
			},
			"source_definition_id": &schema.StringAttribute{
				Description: "Definition ID",
				Required:    true,
			},
			"workspace_id": &schema.StringAttribute{
				Description: "Workspace ID",
				Required:    true,
			},
			"connection_configuration": &schema.MapAttribute{
				Description: "Connection configuration",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"start_date": &schema.StringAttribute{
						Description: "Start Date",
						Required:    true,
					},
					"shop": &schema.StringAttribute{
						Description: "shop",
						Required:    true,
					},
					"credentials": &schema.MapAttribute{
						Description: "Connection configuration",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"auth_method": &schema.StringAttribute{
								Description: "Auth Method",
								Required:    true,
							},
							"api_password": &schema.StringAttribute{
								Description: "API Password",
								Required:    true,
								Sensitive:   true,
							},
							"client_secret": &schema.StringAttribute{
								Description: "Client Secret",
								Required:    true,
								Sensitive:   true,
							},
							"access_token": &schema.StringAttribute{
								Description: "Access Token",
								Required:    true,
								Sensitive:   true,
							},
							"client_id": &schema.StringAttribute{
								Description: "Client ID",
								Required:    true,
								Sensitive:   true,
							},
						},
					},
				},
			},
		},
	}

	sEnc, err := fwhelpers.Encode(s)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		SchemaContents: sEnc,
	}
}

// Create a new resource
func (r *sourceShopifyResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceShopifyResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceShopify{}
	body.Name = plan.Name
	body.SourceDefinitionId = plan.SourceDefinitionId
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceShopifyConnConfig{}
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.Shop = plan.ConnectionConfiguration.Shop
	body.ConnectionConfiguration.Credentials = api.ShopifyCredConfigModel{}
	body.ConnectionConfiguration.Credentials.AuthMethod = plan.ConnectionConfiguration.Credentials.AuthMethod
	body.ConnectionConfiguration.Credentials.ApiPassword = plan.ConnectionConfiguration.Credentials.ApiPassword
	body.ConnectionConfiguration.Credentials.ClientSecret = plan.ConnectionConfiguration.Credentials.ClientSecret
	body.ConnectionConfiguration.Credentials.ClientId = plan.ConnectionConfiguration.Credentials.ClientId
	// Create new source
	source, err := r.Client.CreateShopifySource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceShopifyResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceShopifyConnConfigModel{}
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.Shop = source.ConnectionConfiguration.Shop
	state.ConnectionConfiguration.Credentials = shopifyCredConfigModel{}
	state.ConnectionConfiguration.Credentials.AuthMethod = source.ConnectionConfiguration.Credentials.AuthMethod
	state.ConnectionConfiguration.Credentials.ApiPassword = source.ConnectionConfiguration.Credentials.ApiPassword
	state.ConnectionConfiguration.Credentials.ClientSecret = source.ConnectionConfiguration.Credentials.ClientSecret
	state.ConnectionConfiguration.Credentials.ClientId = source.ConnectionConfiguration.Credentials.ClientId

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.SourceId,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Read resource information
func (r *sourceShopifyResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceShopifyResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadShopifySource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceDefinitionId = source.SourceDefinitionId
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceShopifyConnConfigModel{}
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
		state.ConnectionConfiguration.Shop = source.ConnectionConfiguration.Shop
		state.ConnectionConfiguration.Credentials = shopifyCredConfigModel{}
		state.ConnectionConfiguration.Credentials.AuthMethod = source.ConnectionConfiguration.Credentials.AuthMethod
		state.ConnectionConfiguration.Credentials.ApiPassword = source.ConnectionConfiguration.Credentials.ApiPassword
		state.ConnectionConfiguration.Credentials.ClientSecret = source.ConnectionConfiguration.Credentials.ClientSecret
		state.ConnectionConfiguration.Credentials.ClientId = source.ConnectionConfiguration.Credentials.ClientId

		res.StateID = state.SourceId
	} else {
		// No previous state exists.
		res.StateID = ""
	}

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}
	res.StateContents = stateEnc

	return &res
}

func (r *sourceShopifyResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceShopifyResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceShopify{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceShopifyConnConfig{}
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.Shop = plan.ConnectionConfiguration.Shop
	body.ConnectionConfiguration.Credentials = api.ShopifyCredConfigModel{}
	body.ConnectionConfiguration.Credentials.AuthMethod = plan.ConnectionConfiguration.Credentials.AuthMethod
	body.ConnectionConfiguration.Credentials.ApiPassword = plan.ConnectionConfiguration.Credentials.ApiPassword
	body.ConnectionConfiguration.Credentials.ClientSecret = plan.ConnectionConfiguration.Credentials.ClientSecret
	body.ConnectionConfiguration.Credentials.ClientId = plan.ConnectionConfiguration.Credentials.ClientId
	// Update existing source
	_, err = r.Client.UpdateShopifySource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadShopifySource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceShopifyResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceShopifyConnConfigModel{}
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.Shop = source.ConnectionConfiguration.Shop
	state.ConnectionConfiguration.Credentials = shopifyCredConfigModel{}
	state.ConnectionConfiguration.Credentials.AuthMethod = source.ConnectionConfiguration.Credentials.AuthMethod
	state.ConnectionConfiguration.Credentials.ApiPassword = source.ConnectionConfiguration.Credentials.ApiPassword
	state.ConnectionConfiguration.Credentials.ClientSecret = source.ConnectionConfiguration.Credentials.ClientSecret
	state.ConnectionConfiguration.Credentials.ClientId = source.ConnectionConfiguration.Credentials.ClientId

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.SourceId,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Delete deletes the resource and removes the state on success.
func (r *sourceShopifyResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteShopifySource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
