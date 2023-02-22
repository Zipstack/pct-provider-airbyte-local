package plugin

import (
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte/api"
)

// Resource implementation.
type sourceFakerResource struct {
	Client *api.Client
}

type sourceFakerResourceModel struct {
	Name                    string                     `cty:"name"`
	SourceDefinitionId      string                     `cty:"source_definition_id"`
	SourceId                string                     `cty:"source_id"`
	WorkspaceId             string                     `cty:"workspace_id"`
	ConnectionConfiguration sourceFakerConnConfigModel `cty:"connection_configuration"`
}

type sourceFakerConnConfigModel struct {
	Seed            int64 `cty:"seed"`
	Count           int64 `cty:"count"`
	RecordsPerSync  int64 `cty:"records_per_sync"`
	RecordsPerSlice int64 `cty:"records_per_slice"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceFakerResource{}
)

// Helper function to return a resource service instance.
func NewSourceFakerResource() schema.ResourceService {
	return &sourceFakerResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceFakerResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_faker",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceFakerResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
	if req.ResourceData == "" {
		return schema.ErrorResponse(fmt.Errorf("no data provided to configure resource"))
	}

	var client *api.Client
	err := fwhelpers.Decode(req.ResourceData, &client)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	r.Client = client

	return &schema.ServiceResponse{}
}

// Schema defines the schema for the resource.
func (r *sourceFakerResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source Faker resource for Airbyte",
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
					"seed": &schema.IntAttribute{
						Description: "Seed",
						Required:    true,
					},
					"count": &schema.IntAttribute{
						Description: "Count",
						Required:    true,
					},
					"records_per_sync": &schema.IntAttribute{
						Description: "Records per sync",
						Required:    true,
					},
					"records_per_slice": &schema.IntAttribute{
						Description: "Records per slice",
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
func (r *sourceFakerResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var plan sourceFakerResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceFaker{}
	body.Name = plan.Name
	body.SourceDefinitionId = plan.SourceDefinitionId
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceFakerConnConfig{}
	body.ConnectionConfiguration.Seed = plan.ConnectionConfiguration.Seed
	body.ConnectionConfiguration.Count = plan.ConnectionConfiguration.Count
	body.ConnectionConfiguration.RecordsPerSync = plan.ConnectionConfiguration.RecordsPerSync
	body.ConnectionConfiguration.RecordsPerSlice = plan.ConnectionConfiguration.RecordsPerSlice

	// Create new source
	source, err := r.Client.CreateSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Map response body to schema
	state := sourceFakerResourceModel{}
	state.Name = source.Name
	state.SourceDefinitionId = source.SourceDefinitionId
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceFakerConnConfigModel{}
	state.ConnectionConfiguration.Seed = source.ConnectionConfiguration.Seed
	state.ConnectionConfiguration.Count = source.ConnectionConfiguration.Count
	state.ConnectionConfiguration.RecordsPerSync = source.ConnectionConfiguration.RecordsPerSync
	state.ConnectionConfiguration.RecordsPerSlice = source.ConnectionConfiguration.RecordsPerSlice

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID: state.SourceId,
		// StateLastUpdated: time.Now().Format(time.RFC850),
		StateContents: stateEnc,
	}

}

// Read resource information
func (r *sourceFakerResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceFakerResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if state.SourceId != "" {
		source, err := r.Client.ReadSource(state.SourceId)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Overwrite state with refreshed value
		state.Name = source.Name
		state.SourceDefinitionId = source.SourceDefinitionId
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceFakerConnConfigModel{}
		state.ConnectionConfiguration.Seed = source.ConnectionConfiguration.Seed
		state.ConnectionConfiguration.Count = source.ConnectionConfiguration.Count
		state.ConnectionConfiguration.RecordsPerSync = source.ConnectionConfiguration.RecordsPerSync
		state.ConnectionConfiguration.RecordsPerSlice = source.ConnectionConfiguration.RecordsPerSlice

		res.StateID = state.SourceId
		// res.StateLastUpdated = time.Now().Format(time.RFC850)
	}

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}
	res.StateContents = stateEnc

	return &res
}

func (r *sourceFakerResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{}

	// // Retrieve values from plan
	// var plan sourceResourceModel
	// err := fwhelpers.UnpackModel(req.StateContents, &plan)
	// if err != nil {
	// 	return schema.ErrorResponse(err)
	// }

	// // Generate API request body from plan
	// // plan??

	// // Update existing source
	// source, err := r.Client.UpdateSource(plan)
	// if err != nil {
	// 	return schema.ErrorResponse(err)
	// }

	// // Fetch updated items from GetSource as UpdateSource items are not
	// // populated.
	// source, err := r.Client.GetSource(state.ID)
	// if err != nil {
	// 	return schema.ErrorResponse(err)
	// }

	// // Update resource state with updated items and timestamp
	// state.Items = []sourceResourceModel{}
	// // fill state.Items with source
	// // plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// // Set refreshed state
	// stateEnc, err := fwhelpers.Encode(state)
	// if err != nil {
	// 	return schema.ErrorResponse(err)
	// }

	// return &schema.ServiceResponse{
	// 	StateContents: stateEnc,
	// }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sourceFakerResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{}

	// // Retrieve values from state
	// var state sourceResourceModel
	// err := fwhelpers.UnpackModel(req.StateContents, &state)
	// if err != nil {
	// 	return schema.ErrorResponse(err)
	// }

	// // Delete existing source
	// source, err := r.Client.DeleteSource(state.ID)
	// if err != nil {
	// 	return schema.ErrorResponse(err)
	// }

	// return &schema.ServiceResponse{}
}