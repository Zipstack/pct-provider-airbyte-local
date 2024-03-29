package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-local/api"
)

// Resource implementation.
type sourceHubspotResource struct {
	Client *api.Client
}

type sourceHubspotResourceModel struct {
	Name                    string                       `pctsdk:"name"`
	SourceId                string                       `pctsdk:"source_id"`
	SourceDefinitionId      string                       `pctsdk:"source_definition_id"`
	WorkspaceId             string                       `pctsdk:"workspace_id"`
	ConnectionConfiguration sourceHubspotConnConfigModel `pctsdk:"connection_configuration"`
}

type sourceHubspotConnConfigModel struct {
	StartDate   string                 `pctsdk:"start_date"`
	Credentials HubspotCredConfigModel `pctsdk:"credentials"`
}

type HubspotCredConfigModel struct {
	CredentialsTitle string `pctsdk:"credentials_title"`
	RefreshToken     string `pctsdk:"refresh_token"`
	AccessToken      string `pctsdk:"access_token"`
	ClientSecret     string `pctsdk:"client_secret"`
	ClientId         string `pctsdk:"client_id"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceHubspotResource{}
)

// Helper function to return a resource service instance.
func NewSourceHubspotResource() schema.ResourceService {
	return &sourceHubspotResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceHubspotResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_hubspot",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceHubspotResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceHubspotResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source Hubspot resource for Airbyte",
		Attributes: map[string]schema.Attribute{
			"name": &schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			"source_id": &schema.StringAttribute{
				Description: "Source ID",
				Required:    false,
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
					"credentials": &schema.MapAttribute{
						Description: "credentials",
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"credentials_title": &schema.StringAttribute{
								Description: "credentials",
								Required:    true,
							},
							"access_token": &schema.StringAttribute{
								Description: "Access Token",
								Sensitive:   true,
								Required:    true,
								Optional:    true,
							},
							"refresh_token": &schema.StringAttribute{
								Description: "Refresh Token",
								Sensitive:   true,
								Required:    true,
								Optional:    true,
							},
							"client_secret": &schema.StringAttribute{
								Description: "Client Secret",
								Sensitive:   true,
								Required:    true,
								Optional:    true,
							},
							"client_id": &schema.StringAttribute{
								Description: "Client ID",
								Required:    true,
								Optional:    true,
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
func (r *sourceHubspotResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceHubspotResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceHubspot{}
	body.Name = plan.Name
	body.SourceDefinitionId = plan.SourceDefinitionId
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceHubspotConnConfig{}
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate

	body.ConnectionConfiguration.Credentials = api.HubspotCredConfigModel{}
	body.ConnectionConfiguration.Credentials.CredentialsTitle = plan.ConnectionConfiguration.Credentials.CredentialsTitle
	body.ConnectionConfiguration.Credentials.RefreshToken = plan.ConnectionConfiguration.Credentials.RefreshToken
	body.ConnectionConfiguration.Credentials.ClientSecret = plan.ConnectionConfiguration.Credentials.ClientSecret
	body.ConnectionConfiguration.Credentials.ClientId = plan.ConnectionConfiguration.Credentials.ClientId
	body.ConnectionConfiguration.Credentials.AccessToken = plan.ConnectionConfiguration.Credentials.AccessToken
	// Create new source
	source, err := r.Client.CreateHubspotSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceHubspotResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceHubspotConnConfigModel{}
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate

	state.ConnectionConfiguration.Credentials = HubspotCredConfigModel{}
	state.ConnectionConfiguration.Credentials.CredentialsTitle = source.ConnectionConfiguration.Credentials.CredentialsTitle
	state.ConnectionConfiguration.Credentials.RefreshToken = source.ConnectionConfiguration.Credentials.RefreshToken
	state.ConnectionConfiguration.Credentials.ClientSecret = source.ConnectionConfiguration.Credentials.ClientSecret
	state.ConnectionConfiguration.Credentials.ClientId = source.ConnectionConfiguration.Credentials.ClientId
	state.ConnectionConfiguration.Credentials.AccessToken = source.ConnectionConfiguration.Credentials.AccessToken

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
func (r *sourceHubspotResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceHubspotResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadHubspotSource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceDefinitionId = source.SourceDefinitionId
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceHubspotConnConfigModel{}
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate

		state.ConnectionConfiguration.Credentials = HubspotCredConfigModel{}
		state.ConnectionConfiguration.Credentials.CredentialsTitle = source.ConnectionConfiguration.Credentials.CredentialsTitle
		state.ConnectionConfiguration.Credentials.RefreshToken = source.ConnectionConfiguration.Credentials.RefreshToken
		state.ConnectionConfiguration.Credentials.ClientSecret = source.ConnectionConfiguration.Credentials.ClientSecret
		state.ConnectionConfiguration.Credentials.ClientId = source.ConnectionConfiguration.Credentials.ClientId
		state.ConnectionConfiguration.Credentials.AccessToken = source.ConnectionConfiguration.Credentials.AccessToken

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

func (r *sourceHubspotResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceHubspotResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceHubspot{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceHubspotConnConfig{}
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate

	body.ConnectionConfiguration.Credentials = api.HubspotCredConfigModel{}
	body.ConnectionConfiguration.Credentials.CredentialsTitle = plan.ConnectionConfiguration.Credentials.CredentialsTitle
	body.ConnectionConfiguration.Credentials.RefreshToken = plan.ConnectionConfiguration.Credentials.RefreshToken
	body.ConnectionConfiguration.Credentials.ClientSecret = plan.ConnectionConfiguration.Credentials.ClientSecret
	body.ConnectionConfiguration.Credentials.ClientId = plan.ConnectionConfiguration.Credentials.ClientId
	body.ConnectionConfiguration.Credentials.AccessToken = plan.ConnectionConfiguration.Credentials.AccessToken
	// Update existing source
	_, err = r.Client.UpdateHubspotSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadHubspotSource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceHubspotResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceHubspotConnConfigModel{}
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate

	state.ConnectionConfiguration.Credentials = HubspotCredConfigModel{}
	state.ConnectionConfiguration.Credentials.CredentialsTitle = source.ConnectionConfiguration.Credentials.CredentialsTitle
	state.ConnectionConfiguration.Credentials.RefreshToken = source.ConnectionConfiguration.Credentials.RefreshToken
	state.ConnectionConfiguration.Credentials.ClientSecret = source.ConnectionConfiguration.Credentials.ClientSecret
	state.ConnectionConfiguration.Credentials.ClientId = source.ConnectionConfiguration.Credentials.ClientId
	state.ConnectionConfiguration.Credentials.AccessToken = source.ConnectionConfiguration.Credentials.AccessToken

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
func (r *sourceHubspotResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteHubspotSource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
