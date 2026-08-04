package convert

import (
	"strings"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func CloudConnectionToModel(c *uigraphapi.CloudConnection) *model.CloudConnection {
	if c == nil {
		return nil
	}
	return &model.CloudConnection{
		ID:             c.ID,
		OrgID:          c.OrgID,
		Provider:       model.CloudProvider(strings.ToUpper(c.Provider)),
		DisplayName:    c.DisplayName,
		Status:         model.CloudConnectionStatus(strings.ToUpper(c.Status)),
		StatusMessage:  c.StatusMessage,
		LastVerifiedAt: c.LastVerifiedAt,
		CreatedBy:      c.CreatedBy,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func CloudConnectionListToModel(rows []uigraphapi.CloudConnection) []*model.CloudConnection {
	out := make([]*model.CloudConnection, len(rows))
	for i := range rows {
		out[i] = CloudConnectionToModel(&rows[i])
	}
	return out
}

func ServiceCostTagRuleToModel(r *uigraphapi.ServiceCostTagRule) *model.ServiceCostTagRule {
	if r == nil {
		return nil
	}
	return &model.ServiceCostTagRule{
		ID:        r.ID,
		ServiceID: r.ServiceID,
		TagKey:    r.TagKey,
		TagValue:  r.TagValue,
		CreatedBy: r.CreatedBy,
		CreatedAt: r.CreatedAt,
	}
}

func ServiceCostTagRuleListToModel(rows []uigraphapi.ServiceCostTagRule) []*model.ServiceCostTagRule {
	out := make([]*model.ServiceCostTagRule, len(rows))
	for i := range rows {
		out[i] = ServiceCostTagRuleToModel(&rows[i])
	}
	return out
}

func InfraResourceToModel(r *uigraphapi.InfraResource) *model.InfraResource {
	if r == nil {
		return nil
	}
	return &model.InfraResource{
		ID:                 r.ID,
		ExternalResourceID: r.ExternalResourceID,
		Name:               r.Name,
		ResourceType:       r.ResourceType,
		Provider:           model.CloudProvider(strings.ToUpper(r.Provider)),
		Region:             r.Region,
		Environment:        r.Environment,
		Status:             model.InfraResourceStatus(strings.ToUpper(r.Status)),
		MonthlyCostUsd:     r.MonthlyCostUSD,
		Tags:               r.Tags,
		MatchedTags:        r.MatchedTags,
		LastSyncedAt:       r.LastSyncedAt,
	}
}

func InfraResourceListToModel(rows []uigraphapi.InfraResource) []*model.InfraResource {
	out := make([]*model.InfraResource, len(rows))
	for i := range rows {
		out[i] = InfraResourceToModel(&rows[i])
	}
	return out
}

func CostTrendPointToModel(p uigraphapi.CostTrendPoint) *model.CostTrendPoint {
	return &model.CostTrendPoint{
		Date:     p.Date,
		TotalUsd: p.TotalUSD,
		AWSUsd:   p.AWSUSD,
		AzureUsd: p.AzureUSD,
		GCPUsd:   p.GCPUSD,
	}
}

func CostTrendPointListToModel(rows []uigraphapi.CostTrendPoint) []*model.CostTrendPoint {
	out := make([]*model.CostTrendPoint, len(rows))
	for i := range rows {
		out[i] = CostTrendPointToModel(rows[i])
	}
	return out
}

func ResourceDailyCostListToModel(rows []uigraphapi.ResourceDailyCost) []*model.ResourceDailyCost {
	out := make([]*model.ResourceDailyCost, len(rows))
	for i, r := range rows {
		out[i] = &model.ResourceDailyCost{Date: r.Date, CostUsd: r.CostUSD}
	}
	return out
}

func ServiceCostSummaryToModel(s *uigraphapi.ServiceCostSummary) *model.ServiceCostSummary {
	if s == nil {
		return nil
	}
	return &model.ServiceCostSummary{
		TotalMonthlyCostUsd: s.TotalMonthlyCostUSD,
		MomChangePct:        s.MoMChangePct,
		ResourceCount:       s.ResourceCount,
		ProviderCount:       s.ProviderCount,
		TopCostDriverLabel:  s.TopCostDriverLabel,
		TopCostDriverUsd:    s.TopCostDriverUSD,
	}
}

func TestCloudConnectionResultToModel(r *uigraphapi.TestCloudConnectionResult) *model.TestCloudConnectionResult {
	if r == nil {
		return nil
	}
	var errPtr *string
	if r.Error != "" {
		errPtr = &r.Error
	}
	return &model.TestCloudConnectionResult{Ok: r.OK, Error: errPtr}
}

// CreateCloudConnectionInputToBody builds the REST request body for
// creating a connection. uigraph-api's Provider constants are lowercase
// ("aws"/"azure"/"gcp"); the GraphQL enum is uppercase.
func CreateCloudConnectionInputToBody(in model.CreateCloudConnectionInput) map[string]interface{} {
	body := map[string]interface{}{
		"provider":    strings.ToLower(string(in.Provider)),
		"displayName": in.DisplayName,
	}
	setIfPresent(body, "roleArn", in.RoleArn)
	setIfPresent(body, "externalId", in.ExternalID)
	setIfPresent(body, "serviceAccountJson", in.ServiceAccountJSON)
	setIfPresent(body, "billingDataset", in.BillingDataset)
	setIfPresent(body, "tenantId", in.TenantID)
	setIfPresent(body, "clientId", in.ClientID)
	setIfPresent(body, "clientSecret", in.ClientSecret)
	setIfPresent(body, "subscriptionId", in.SubscriptionID)
	return body
}

func setIfPresent(body map[string]interface{}, key string, val *string) {
	if val != nil {
		body[key] = *val
	}
}
