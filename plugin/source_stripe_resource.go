package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte/api"
)

// Resource implementation.
type sourceStripeResource struct {
	Client *api.Client
}

type sourceStripeResourceModel struct {
	Name                    string                      `cty:"name"`
	SourceId                string                      `cty:"source_id"`
	SourceDefinitionId      string                      `cty:"source_definition_id"`
	WorkspaceId             string                      `cty:"workspace_id"`
	ConnectionConfiguration sourceStripeConnConfigModel `cty:"connection_configuration"`
}

type sourceStripeConnConfigModel struct {
	StartDate          string `cty:"start_date"`
	LookbackWindowDays int    `cty:"lookback_window_days"`
	SliceRange         int    `cty:"slice_range"`
	ClientSecret       string `cty:"client_secret"`
	AccountId          string `cty:"account_id"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceStripeResource{}
)

// Helper function to return a resource service instance.
func NewsourceStripeResource() schema.ResourceService {
	return &sourceStripeResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceStripeResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_stripe",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceStripeResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceStripeResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source pipedrive resource for Airbyte",
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
				//Sensitive:   true,
				Attributes: map[string]schema.Attribute{
					"start_date": &schema.StringAttribute{
						Description: "Start Date",
						Required:    true,
					},
					"slice_range": &schema.IntAttribute{
						Description: "Slice Range",
						Required:    false,
					},
					"lookback_window_days": &schema.IntAttribute{
						Description: "lookback window days",
						Required:    false,
					},
					"client_secret": &schema.StringAttribute{
						Description: "Client Secret",
						Required:    true,
					},
					"account_id": &schema.StringAttribute{
						Description: "Account Id",
						Required:    true,
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
func (r *sourceStripeResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceStripeResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceStripe{}
	body.Name = plan.Name
	body.SourceDefinitionId = plan.SourceDefinitionId
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceStripeConnConfig{}
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.ClientSecret = plan.ConnectionConfiguration.ClientSecret
	body.ConnectionConfiguration.AccountId = plan.ConnectionConfiguration.AccountId
	body.ConnectionConfiguration.LookbackWindowDays = plan.ConnectionConfiguration.LookbackWindowDays
	body.ConnectionConfiguration.SliceRange = plan.ConnectionConfiguration.SliceRange

	// Create new source
	source, err := r.Client.CreateStripeSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceStripeResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceStripeConnConfigModel{}
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.ClientSecret = source.ConnectionConfiguration.ClientSecret
	state.ConnectionConfiguration.AccountId = source.ConnectionConfiguration.AccountId
	state.ConnectionConfiguration.LookbackWindowDays = source.ConnectionConfiguration.LookbackWindowDays
	state.ConnectionConfiguration.SliceRange = source.ConnectionConfiguration.SliceRange

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
func (r *sourceStripeResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceStripeResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadStripeSource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceDefinitionId = source.SourceDefinitionId
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceStripeConnConfigModel{}
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
		state.ConnectionConfiguration.ClientSecret = source.ConnectionConfiguration.ClientSecret
		state.ConnectionConfiguration.AccountId = source.ConnectionConfiguration.AccountId
		state.ConnectionConfiguration.LookbackWindowDays = source.ConnectionConfiguration.LookbackWindowDays
		state.ConnectionConfiguration.SliceRange = source.ConnectionConfiguration.SliceRange

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

func (r *sourceStripeResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceStripeResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceStripe{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceStripeConnConfig{}
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.ClientSecret = plan.ConnectionConfiguration.ClientSecret
	body.ConnectionConfiguration.AccountId = plan.ConnectionConfiguration.AccountId
	body.ConnectionConfiguration.LookbackWindowDays = plan.ConnectionConfiguration.LookbackWindowDays
	body.ConnectionConfiguration.SliceRange = plan.ConnectionConfiguration.SliceRange

	// Update existing source
	_, err = r.Client.UpdateStripeSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadStripeSource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceStripeResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceStripeConnConfigModel{}
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.ClientSecret = source.ConnectionConfiguration.ClientSecret
	state.ConnectionConfiguration.AccountId = source.ConnectionConfiguration.AccountId
	state.ConnectionConfiguration.LookbackWindowDays = source.ConnectionConfiguration.LookbackWindowDays
	state.ConnectionConfiguration.SliceRange = source.ConnectionConfiguration.SliceRange

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
func (r *sourceStripeResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteStripeSource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
