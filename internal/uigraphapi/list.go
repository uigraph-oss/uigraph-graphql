package uigraphapi

import (
	"net/url"
	"strconv"
)

func listQuery(p ListParams) string {
	q := url.Values{}
	if p.FolderID != "" {
		q.Set("folderId", p.FolderID)
	}
	if p.TeamID != "" {
		q.Set("teamId", p.TeamID)
	}
	if p.ServiceID != "" {
		q.Set("serviceId", p.ServiceID)
	}
	if p.Search != "" {
		q.Set("search", p.Search)
	}
	if p.SortBy != "" {
		q.Set("sortBy", p.SortBy)
	}
	if p.SortDir != "" {
		q.Set("sortDir", p.SortDir)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}
