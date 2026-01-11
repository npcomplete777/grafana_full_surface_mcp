package tools

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aaronjacobs/grafana-mcp-server/internal/grafana"
	"github.com/aaronjacobs/grafana-mcp-server/internal/mcp"
)

// Registry holds all tool definitions and handlers
type Registry struct {
	client *grafana.Client
	tools  map[string]ToolHandler
}

// ToolHandler processes a tool call
type ToolHandler func(args map[string]interface{}) (*mcp.CallToolResult, error)

// NewRegistry creates a new tool registry
func NewRegistry(client *grafana.Client) *Registry {
	r := &Registry{
		client: client,
		tools:  make(map[string]ToolHandler),
	}
	r.registerAll()
	return r
}

// GetTools returns all tool definitions
func (r *Registry) GetTools() []mcp.Tool {
	return []mcp.Tool{
		// Health
		r.grafanaHealthTool(),
		
		// Dashboard tools
		r.grafanaSearchDashboardsTool(),
		r.grafanaGetDashboardTool(),
		r.grafanaCreateDashboardTool(),
		r.grafanaUpdateDashboardTool(),
		r.grafanaDeleteDashboardTool(),
		
		// Datasource tools
		r.grafanaListDatasourcesTool(),
		r.grafanaGetDatasourceTool(),
		r.grafanaCreateDatasourceTool(),
		r.grafanaUpdateDatasourceTool(),
		r.grafanaDeleteDatasourceTool(),
		
		// Folder tools
		r.grafanaListFoldersTool(),
		r.grafanaGetFolderTool(),
		r.grafanaCreateFolderTool(),
		r.grafanaUpdateFolderTool(),
		r.grafanaDeleteFolderTool(),
		
		// Alert tools
		r.grafanaListAlertRulesTool(),
		r.grafanaGetAlertRuleTool(),
		r.grafanaCreateAlertRuleTool(),
		r.grafanaUpdateAlertRuleTool(),
		r.grafanaDeleteAlertRuleTool(),
		
		// Annotation tools
		r.grafanaListAnnotationsTool(),
		r.grafanaCreateAnnotationTool(),
		r.grafanaUpdateAnnotationTool(),
		r.grafanaDeleteAnnotationTool(),
		
		// Query tools
		r.grafanaQueryTool(),
		
		// Organization tools
		r.grafanaGetOrgTool(),
		r.grafanaListOrgUsersTool(),
		
		// User tools
		r.grafanaGetCurrentUserTool(),
		
		// Team tools
		r.grafanaListTeamsTool(),
		r.grafanaGetTeamTool(),
		r.grafanaCreateTeamTool(),
		r.grafanaDeleteTeamTool(),
	}
}

// CallTool executes a tool by name
func (r *Registry) CallTool(name string, args map[string]interface{}) (*mcp.CallToolResult, error) {
	handler, ok := r.tools[name]
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.ContentBlock{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", name)}},
		}, nil
	}
	return handler(args)
}

func (r *Registry) registerAll() {
	// Health
	r.tools["grafana_health"] = r.handleHealth
	
	// Dashboards
	r.tools["grafana_search_dashboards"] = r.handleSearchDashboards
	r.tools["grafana_get_dashboard"] = r.handleGetDashboard
	r.tools["grafana_create_dashboard"] = r.handleCreateDashboard
	r.tools["grafana_update_dashboard"] = r.handleUpdateDashboard
	r.tools["grafana_delete_dashboard"] = r.handleDeleteDashboard
	
	// Datasources
	r.tools["grafana_list_datasources"] = r.handleListDatasources
	r.tools["grafana_get_datasource"] = r.handleGetDatasource
	r.tools["grafana_create_datasource"] = r.handleCreateDatasource
	r.tools["grafana_update_datasource"] = r.handleUpdateDatasource
	r.tools["grafana_delete_datasource"] = r.handleDeleteDatasource
	
	// Folders
	r.tools["grafana_list_folders"] = r.handleListFolders
	r.tools["grafana_get_folder"] = r.handleGetFolder
	r.tools["grafana_create_folder"] = r.handleCreateFolder
	r.tools["grafana_update_folder"] = r.handleUpdateFolder
	r.tools["grafana_delete_folder"] = r.handleDeleteFolder
	
	// Alerts
	r.tools["grafana_list_alert_rules"] = r.handleListAlertRules
	r.tools["grafana_get_alert_rule"] = r.handleGetAlertRule
	r.tools["grafana_create_alert_rule"] = r.handleCreateAlertRule
	r.tools["grafana_update_alert_rule"] = r.handleUpdateAlertRule
	r.tools["grafana_delete_alert_rule"] = r.handleDeleteAlertRule
	
	// Annotations
	r.tools["grafana_list_annotations"] = r.handleListAnnotations
	r.tools["grafana_create_annotation"] = r.handleCreateAnnotation
	r.tools["grafana_update_annotation"] = r.handleUpdateAnnotation
	r.tools["grafana_delete_annotation"] = r.handleDeleteAnnotation
	
	// Query
	r.tools["grafana_query"] = r.handleQuery
	
	// Organization
	r.tools["grafana_get_org"] = r.handleGetOrg
	r.tools["grafana_list_org_users"] = r.handleListOrgUsers
	
	// User
	r.tools["grafana_get_current_user"] = r.handleGetCurrentUser
	
	// Teams
	r.tools["grafana_list_teams"] = r.handleListTeams
	r.tools["grafana_get_team"] = r.handleGetTeam
	r.tools["grafana_create_team"] = r.handleCreateTeam
	r.tools["grafana_delete_team"] = r.handleDeleteTeam
}

// Helper functions
func jsonResult(v interface{}) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.ContentBlock{{Type: "text", Text: string(data)}},
	}, nil
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.ContentBlock{{Type: "text", Text: msg}},
	}
}

func getString(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt64(args map[string]interface{}, key string) int64 {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		case int:
			return int64(n)
		}
	}
	return 0
}

func getInt(args map[string]interface{}, key string) int {
	return int(getInt64(args, key))
}

func getBool(args map[string]interface{}, key string) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func getStringSlice(args map[string]interface{}, key string) []string {
	if v, ok := args[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]string, len(arr))
			for i, item := range arr {
				if s, ok := item.(string); ok {
					result[i] = s
				}
			}
			return result
		}
	}
	return nil
}

// ============== Tool Definitions ==============

func (r *Registry) grafanaHealthTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_health",
		Description: "Check Grafana server health and version information",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaSearchDashboardsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_search_dashboards",
		Description: "Search for dashboards by query, tags, or folder",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"query": {Type: "string", Description: "Search query string"},
				"tags":  {Type: "array", Description: "Filter by tags"},
				"type":  {Type: "string", Description: "Filter by type: dash-db or dash-folder", Enum: []string{"dash-db", "dash-folder"}},
				"limit": {Type: "integer", Description: "Maximum number of results (default 50)"},
			},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_dashboard",
		Description: "Get a dashboard by its UID, including all panels and configuration",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Dashboard UID"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaCreateDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_create_dashboard",
		Description: "Create a new dashboard with panels",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"title":      {Type: "string", Description: "Dashboard title"},
				"tags":       {Type: "array", Description: "Dashboard tags"},
				"folder_uid": {Type: "string", Description: "Folder UID to save dashboard in"},
				"panels":     {Type: "array", Description: "Array of panel configurations"},
				"refresh":    {Type: "string", Description: "Auto-refresh interval (e.g., 5s, 1m, 5m)"},
				"time_from":  {Type: "string", Description: "Time range from (e.g., now-6h)"},
				"time_to":    {Type: "string", Description: "Time range to (e.g., now)"},
			},
			Required: []string{"title"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: false,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaUpdateDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_update_dashboard",
		Description: "Update an existing dashboard",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid":        {Type: "string", Description: "Dashboard UID to update"},
				"title":      {Type: "string", Description: "New dashboard title"},
				"tags":       {Type: "array", Description: "Dashboard tags"},
				"panels":     {Type: "array", Description: "Array of panel configurations"},
				"folder_uid": {Type: "string", Description: "Folder UID to move dashboard to"},
				"message":    {Type: "string", Description: "Save message/commit description"},
				"overwrite":  {Type: "boolean", Description: "Overwrite existing dashboard"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaDeleteDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_delete_dashboard",
		Description: "Delete a dashboard by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Dashboard UID to delete"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaListDatasourcesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_list_datasources",
		Description: "List all configured datasources",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetDatasourceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_datasource",
		Description: "Get a datasource by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Datasource UID"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaCreateDatasourceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_create_datasource",
		Description: "Create a new datasource",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"name":       {Type: "string", Description: "Datasource name"},
				"type":       {Type: "string", Description: "Datasource type (e.g., prometheus, loki, elasticsearch)"},
				"url":        {Type: "string", Description: "Datasource URL"},
				"access":     {Type: "string", Description: "Access mode: proxy or direct", Enum: []string{"proxy", "direct"}},
				"is_default": {Type: "boolean", Description: "Set as default datasource"},
				"json_data":  {Type: "object", Description: "Additional JSON configuration"},
			},
			Required: []string{"name", "type", "url"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: false,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaUpdateDatasourceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_update_datasource",
		Description: "Update an existing datasource",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid":        {Type: "string", Description: "Datasource UID to update"},
				"name":       {Type: "string", Description: "New datasource name"},
				"url":        {Type: "string", Description: "New datasource URL"},
				"is_default": {Type: "boolean", Description: "Set as default datasource"},
				"json_data":  {Type: "object", Description: "Additional JSON configuration"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaDeleteDatasourceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_delete_datasource",
		Description: "Delete a datasource by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Datasource UID to delete"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaListFoldersTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_list_folders",
		Description: "List all folders",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetFolderTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_folder",
		Description: "Get a folder by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Folder UID"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaCreateFolderTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_create_folder",
		Description: "Create a new folder",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"title": {Type: "string", Description: "Folder title"},
				"uid":   {Type: "string", Description: "Optional folder UID"},
			},
			Required: []string{"title"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: false,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaUpdateFolderTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_update_folder",
		Description: "Update a folder",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid":     {Type: "string", Description: "Folder UID to update"},
				"title":   {Type: "string", Description: "New folder title"},
				"version": {Type: "integer", Description: "Current folder version for optimistic locking"},
			},
			Required: []string{"uid", "title", "version"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaDeleteFolderTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_delete_folder",
		Description: "Delete a folder and all its dashboards",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Folder UID to delete"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaListAlertRulesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_list_alert_rules",
		Description: "List all alert rules",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_alert_rule",
		Description: "Get an alert rule by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Alert rule UID"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaCreateAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_create_alert_rule",
		Description: "Create a new alert rule",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"title":          {Type: "string", Description: "Alert rule title"},
				"folder_uid":     {Type: "string", Description: "Folder UID to store the alert"},
				"rule_group":     {Type: "string", Description: "Rule group name"},
				"condition":      {Type: "string", Description: "Condition refId"},
				"queries":        {Type: "array", Description: "Array of query configurations"},
				"for_duration":   {Type: "string", Description: "Duration before alert fires (e.g., 5m)"},
				"no_data_state":  {Type: "string", Description: "State when no data: NoData, Alerting, OK", Enum: []string{"NoData", "Alerting", "OK"}},
				"exec_err_state": {Type: "string", Description: "State on execution error: Alerting, Error, OK", Enum: []string{"Alerting", "Error", "OK"}},
				"labels":         {Type: "object", Description: "Labels to attach to alert"},
				"annotations":    {Type: "object", Description: "Annotations (summary, description, etc.)"},
			},
			Required: []string{"title", "folder_uid", "rule_group", "condition", "queries"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: false,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaUpdateAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_update_alert_rule",
		Description: "Update an existing alert rule",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid":            {Type: "string", Description: "Alert rule UID to update"},
				"title":          {Type: "string", Description: "New alert rule title"},
				"queries":        {Type: "array", Description: "New query configurations"},
				"for_duration":   {Type: "string", Description: "Duration before alert fires"},
				"no_data_state":  {Type: "string", Description: "State when no data"},
				"exec_err_state": {Type: "string", Description: "State on execution error"},
				"labels":         {Type: "object", Description: "Labels to attach to alert"},
				"annotations":    {Type: "object", Description: "Annotations"},
				"is_paused":      {Type: "boolean", Description: "Pause the alert rule"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaDeleteAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_delete_alert_rule",
		Description: "Delete an alert rule by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {Type: "string", Description: "Alert rule UID to delete"},
			},
			Required: []string{"uid"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaListAnnotationsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_list_annotations",
		Description: "List annotations with optional filters",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"from":          {Type: "integer", Description: "Start time in epoch milliseconds"},
				"to":            {Type: "integer", Description: "End time in epoch milliseconds"},
				"dashboard_uid": {Type: "string", Description: "Filter by dashboard UID"},
				"panel_id":      {Type: "integer", Description: "Filter by panel ID"},
				"tags":          {Type: "array", Description: "Filter by tags"},
				"limit":         {Type: "integer", Description: "Maximum number of results"},
			},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaCreateAnnotationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_create_annotation",
		Description: "Create a new annotation",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"text":          {Type: "string", Description: "Annotation text"},
				"time":          {Type: "integer", Description: "Start time in epoch milliseconds (default: now)"},
				"time_end":      {Type: "integer", Description: "End time for region annotations"},
				"dashboard_uid": {Type: "string", Description: "Dashboard UID to attach annotation"},
				"panel_id":      {Type: "integer", Description: "Panel ID to attach annotation"},
				"tags":          {Type: "array", Description: "Annotation tags"},
			},
			Required: []string{"text"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: false,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaUpdateAnnotationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_update_annotation",
		Description: "Update an existing annotation",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id":       {Type: "integer", Description: "Annotation ID to update"},
				"text":     {Type: "string", Description: "New annotation text"},
				"time":     {Type: "integer", Description: "New start time"},
				"time_end": {Type: "integer", Description: "New end time"},
				"tags":     {Type: "array", Description: "New tags"},
			},
			Required: []string{"id"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaDeleteAnnotationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_delete_annotation",
		Description: "Delete an annotation by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {Type: "integer", Description: "Annotation ID to delete"},
			},
			Required: []string{"id"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaQueryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_query",
		Description: "Execute a query against a datasource",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"datasource_uid":  {Type: "string", Description: "Datasource UID to query"},
				"datasource_type": {Type: "string", Description: "Datasource type (e.g., prometheus, loki)"},
				"query":           {Type: "string", Description: "Query expression (PromQL for Prometheus, LogQL for Loki, etc.)"},
				"from":            {Type: "string", Description: "Start time (e.g., now-1h, 2024-01-01T00:00:00Z)"},
				"to":              {Type: "string", Description: "End time (e.g., now)"},
				"max_data_points": {Type: "integer", Description: "Maximum number of data points"},
				"interval_ms":     {Type: "integer", Description: "Query interval in milliseconds"},
			},
			Required: []string{"datasource_uid", "datasource_type", "query"},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetOrgTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_org",
		Description: "Get current organization information",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaListOrgUsersTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_list_org_users",
		Description: "List users in the current organization",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetCurrentUserTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_current_user",
		Description: "Get current authenticated user information",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaListTeamsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_list_teams",
		Description: "List teams in the organization",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"query":    {Type: "string", Description: "Search query"},
				"page":     {Type: "integer", Description: "Page number"},
				"per_page": {Type: "integer", Description: "Results per page"},
			},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaGetTeamTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_get_team",
		Description: "Get a team by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {Type: "integer", Description: "Team ID"},
			},
			Required: []string{"id"},
		},
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: true,
		},
	}
}

func (r *Registry) grafanaCreateTeamTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_create_team",
		Description: "Create a new team",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"name":  {Type: "string", Description: "Team name"},
				"email": {Type: "string", Description: "Team email address"},
			},
			Required: []string{"name"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: false,
			OpenWorldHint:   true,
		},
	}
}

func (r *Registry) grafanaDeleteTeamTool() mcp.Tool {
	return mcp.Tool{
		Name:        "grafana_delete_team",
		Description: "Delete a team by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"id": {Type: "integer", Description: "Team ID to delete"},
			},
			Required: []string{"id"},
		},
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: true,
			OpenWorldHint:   true,
		},
	}
}

// ============== Handler Implementations ==============

func (r *Registry) handleHealth(args map[string]interface{}) (*mcp.CallToolResult, error) {
	health, err := r.client.GetHealth()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get health: %v", err)), nil
	}
	return jsonResult(health)
}

func (r *Registry) handleSearchDashboards(args map[string]interface{}) (*mcp.CallToolResult, error) {
	query := getString(args, "query")
	tags := getStringSlice(args, "tags")
	dashType := getString(args, "type")
	limit := getInt(args, "limit")
	if limit == 0 {
		limit = 50
	}

	results, err := r.client.SearchDashboards(query, tags, nil, dashType, limit)
	if err != nil {
		return errorResult(fmt.Sprintf("Search failed: %v", err)), nil
	}
	return jsonResult(results)
}

func (r *Registry) handleGetDashboard(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	dashboard, err := r.client.GetDashboard(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get dashboard: %v", err)), nil
	}
	return jsonResult(dashboard)
}

func (r *Registry) handleCreateDashboard(args map[string]interface{}) (*mcp.CallToolResult, error) {
	title := getString(args, "title")
	if title == "" {
		return errorResult("title is required"), nil
	}

	dashboard := grafana.Dashboard{
		Title:         title,
		Tags:          getStringSlice(args, "tags"),
		SchemaVersion: 39,
		Refresh:       getString(args, "refresh"),
	}

	if timeFrom := getString(args, "time_from"); timeFrom != "" {
		dashboard.Time = &grafana.TimeRange{From: timeFrom, To: getString(args, "time_to")}
	}

	// Handle panels if provided
	if panelsRaw, ok := args["panels"]; ok {
		if panelsArr, ok := panelsRaw.([]interface{}); ok {
			panels := make([]grafana.Panel, 0, len(panelsArr))
			for _, p := range panelsArr {
				if pm, ok := p.(map[string]interface{}); ok {
					panel := grafana.Panel{
						Type:  getString(pm, "type"),
						Title: getString(pm, "title"),
					}
					panels = append(panels, panel)
				}
			}
			dashboard.Panels = panels
		}
	}

	req := grafana.SaveDashboardRequest{
		Dashboard: dashboard,
		FolderUID: getString(args, "folder_uid"),
		Message:   "Created via MCP",
	}

	result, err := r.client.SaveDashboard(req)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to create dashboard: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleUpdateDashboard(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	// Get existing dashboard
	existing, err := r.client.GetDashboard(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get dashboard: %v", err)), nil
	}

	// Update fields
	if title := getString(args, "title"); title != "" {
		existing.Title = title
	}
	if tags := getStringSlice(args, "tags"); len(tags) > 0 {
		existing.Tags = tags
	}

	req := grafana.SaveDashboardRequest{
		Dashboard: *existing,
		FolderUID: getString(args, "folder_uid"),
		Message:   getString(args, "message"),
		Overwrite: getBool(args, "overwrite"),
	}

	result, err := r.client.SaveDashboard(req)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to update dashboard: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleDeleteDashboard(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	if err := r.client.DeleteDashboard(uid); err != nil {
		return errorResult(fmt.Sprintf("Failed to delete dashboard: %v", err)), nil
	}
	return jsonResult(map[string]string{"status": "deleted", "uid": uid})
}

func (r *Registry) handleListDatasources(args map[string]interface{}) (*mcp.CallToolResult, error) {
	datasources, err := r.client.GetDatasources()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to list datasources: %v", err)), nil
	}
	return jsonResult(datasources)
}

func (r *Registry) handleGetDatasource(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	ds, err := r.client.GetDatasource(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get datasource: %v", err)), nil
	}
	return jsonResult(ds)
}

func (r *Registry) handleCreateDatasource(args map[string]interface{}) (*mcp.CallToolResult, error) {
	name := getString(args, "name")
	dsType := getString(args, "type")
	dsURL := getString(args, "url")

	if name == "" || dsType == "" || dsURL == "" {
		return errorResult("name, type, and url are required"), nil
	}

	ds := grafana.Datasource{
		Name:      name,
		Type:      dsType,
		URL:       dsURL,
		Access:    getString(args, "access"),
		IsDefault: getBool(args, "is_default"),
	}

	if ds.Access == "" {
		ds.Access = "proxy"
	}

	if jsonData, ok := args["json_data"].(map[string]interface{}); ok {
		ds.JSONData = jsonData
	}

	result, err := r.client.CreateDatasource(ds)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to create datasource: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleUpdateDatasource(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	existing, err := r.client.GetDatasource(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get datasource: %v", err)), nil
	}

	if name := getString(args, "name"); name != "" {
		existing.Name = name
	}
	if dsURL := getString(args, "url"); dsURL != "" {
		existing.URL = dsURL
	}
	if _, ok := args["is_default"]; ok {
		existing.IsDefault = getBool(args, "is_default")
	}

	result, err := r.client.UpdateDatasource(uid, *existing)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to update datasource: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleDeleteDatasource(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	if err := r.client.DeleteDatasource(uid); err != nil {
		return errorResult(fmt.Sprintf("Failed to delete datasource: %v", err)), nil
	}
	return jsonResult(map[string]string{"status": "deleted", "uid": uid})
}

func (r *Registry) handleListFolders(args map[string]interface{}) (*mcp.CallToolResult, error) {
	folders, err := r.client.GetFolders()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to list folders: %v", err)), nil
	}
	return jsonResult(folders)
}

func (r *Registry) handleGetFolder(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	folder, err := r.client.GetFolder(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get folder: %v", err)), nil
	}
	return jsonResult(folder)
}

func (r *Registry) handleCreateFolder(args map[string]interface{}) (*mcp.CallToolResult, error) {
	title := getString(args, "title")
	if title == "" {
		return errorResult("title is required"), nil
	}

	folder, err := r.client.CreateFolder(title, getString(args, "uid"))
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to create folder: %v", err)), nil
	}
	return jsonResult(folder)
}

func (r *Registry) handleUpdateFolder(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	title := getString(args, "title")
	version := getInt(args, "version")

	if uid == "" || title == "" || version == 0 {
		return errorResult("uid, title, and version are required"), nil
	}

	folder, err := r.client.UpdateFolder(uid, title, version)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to update folder: %v", err)), nil
	}
	return jsonResult(folder)
}

func (r *Registry) handleDeleteFolder(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	if err := r.client.DeleteFolder(uid); err != nil {
		return errorResult(fmt.Sprintf("Failed to delete folder: %v", err)), nil
	}
	return jsonResult(map[string]string{"status": "deleted", "uid": uid})
}

func (r *Registry) handleListAlertRules(args map[string]interface{}) (*mcp.CallToolResult, error) {
	rules, err := r.client.GetAlertRules()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to list alert rules: %v", err)), nil
	}
	return jsonResult(rules)
}

func (r *Registry) handleGetAlertRule(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	rule, err := r.client.GetAlertRule(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get alert rule: %v", err)), nil
	}
	return jsonResult(rule)
}

func (r *Registry) handleCreateAlertRule(args map[string]interface{}) (*mcp.CallToolResult, error) {
	title := getString(args, "title")
	folderUID := getString(args, "folder_uid")
	ruleGroup := getString(args, "rule_group")
	condition := getString(args, "condition")

	if title == "" || folderUID == "" || ruleGroup == "" || condition == "" {
		return errorResult("title, folder_uid, rule_group, and condition are required"), nil
	}

	rule := grafana.AlertRule{
		Title:        title,
		FolderUID:    folderUID,
		RuleGroup:    ruleGroup,
		Condition:    condition,
		NoDataState:  getString(args, "no_data_state"),
		ExecErrState: getString(args, "exec_err_state"),
		For:          getString(args, "for_duration"),
	}

	if labels, ok := args["labels"].(map[string]interface{}); ok {
		rule.Labels = make(map[string]string)
		for k, v := range labels {
			if s, ok := v.(string); ok {
				rule.Labels[k] = s
			}
		}
	}

	if annotations, ok := args["annotations"].(map[string]interface{}); ok {
		rule.Annotations = make(map[string]string)
		for k, v := range annotations {
			if s, ok := v.(string); ok {
				rule.Annotations[k] = s
			}
		}
	}

	result, err := r.client.CreateAlertRule(rule)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to create alert rule: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleUpdateAlertRule(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	existing, err := r.client.GetAlertRule(uid)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get alert rule: %v", err)), nil
	}

	if title := getString(args, "title"); title != "" {
		existing.Title = title
	}
	if forDuration := getString(args, "for_duration"); forDuration != "" {
		existing.For = forDuration
	}
	if noDataState := getString(args, "no_data_state"); noDataState != "" {
		existing.NoDataState = noDataState
	}
	if execErrState := getString(args, "exec_err_state"); execErrState != "" {
		existing.ExecErrState = execErrState
	}
	if _, ok := args["is_paused"]; ok {
		existing.IsPaused = getBool(args, "is_paused")
	}

	result, err := r.client.UpdateAlertRule(uid, *existing)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to update alert rule: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleDeleteAlertRule(args map[string]interface{}) (*mcp.CallToolResult, error) {
	uid := getString(args, "uid")
	if uid == "" {
		return errorResult("uid is required"), nil
	}

	if err := r.client.DeleteAlertRule(uid); err != nil {
		return errorResult(fmt.Sprintf("Failed to delete alert rule: %v", err)), nil
	}
	return jsonResult(map[string]string{"status": "deleted", "uid": uid})
}

func (r *Registry) handleListAnnotations(args map[string]interface{}) (*mcp.CallToolResult, error) {
	from := getInt64(args, "from")
	to := getInt64(args, "to")
	dashboardUID := getString(args, "dashboard_uid")
	panelID := getInt64(args, "panel_id")
	tags := getStringSlice(args, "tags")
	limit := getInt(args, "limit")

	annotations, err := r.client.GetAnnotations(from, to, dashboardUID, panelID, tags, limit)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to list annotations: %v", err)), nil
	}
	return jsonResult(annotations)
}

func (r *Registry) handleCreateAnnotation(args map[string]interface{}) (*mcp.CallToolResult, error) {
	text := getString(args, "text")
	if text == "" {
		return errorResult("text is required"), nil
	}

	ann := grafana.Annotation{
		Text:         text,
		Time:         getInt64(args, "time"),
		TimeEnd:      getInt64(args, "time_end"),
		DashboardUID: getString(args, "dashboard_uid"),
		PanelID:      getInt64(args, "panel_id"),
		Tags:         getStringSlice(args, "tags"),
	}

	if ann.Time == 0 {
		ann.Time = time.Now().UnixMilli()
	}

	result, err := r.client.CreateAnnotation(ann)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to create annotation: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleUpdateAnnotation(args map[string]interface{}) (*mcp.CallToolResult, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return errorResult("id is required"), nil
	}

	ann := grafana.Annotation{
		Text:    getString(args, "text"),
		Time:    getInt64(args, "time"),
		TimeEnd: getInt64(args, "time_end"),
		Tags:    getStringSlice(args, "tags"),
	}

	if err := r.client.UpdateAnnotation(id, ann); err != nil {
		return errorResult(fmt.Sprintf("Failed to update annotation: %v", err)), nil
	}
	return jsonResult(map[string]interface{}{"status": "updated", "id": id})
}

func (r *Registry) handleDeleteAnnotation(args map[string]interface{}) (*mcp.CallToolResult, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return errorResult("id is required"), nil
	}

	if err := r.client.DeleteAnnotation(id); err != nil {
		return errorResult(fmt.Sprintf("Failed to delete annotation: %v", err)), nil
	}
	return jsonResult(map[string]interface{}{"status": "deleted", "id": id})
}

func (r *Registry) handleQuery(args map[string]interface{}) (*mcp.CallToolResult, error) {
	dsUID := getString(args, "datasource_uid")
	dsType := getString(args, "datasource_type")
	query := getString(args, "query")

	if dsUID == "" || dsType == "" || query == "" {
		return errorResult("datasource_uid, datasource_type, and query are required"), nil
	}

	from := getString(args, "from")
	to := getString(args, "to")
	if from == "" {
		from = "now-1h"
	}
	if to == "" {
		to = "now"
	}

	req := grafana.QueryRequest{
		From: from,
		To:   to,
		Queries: []grafana.QueryTarget{
			{
				RefID:         "A",
				Datasource:    grafana.DatasourceRef{Type: dsType, UID: dsUID},
				Query:         query,
				MaxDataPoints: getInt(args, "max_data_points"),
				IntervalMs:    getInt(args, "interval_ms"),
			},
		},
	}

	result, err := r.client.Query(req)
	if err != nil {
		return errorResult(fmt.Sprintf("Query failed: %v", err)), nil
	}
	return jsonResult(result)
}

func (r *Registry) handleGetOrg(args map[string]interface{}) (*mcp.CallToolResult, error) {
	org, err := r.client.GetCurrentOrg()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get organization: %v", err)), nil
	}
	return jsonResult(org)
}

func (r *Registry) handleListOrgUsers(args map[string]interface{}) (*mcp.CallToolResult, error) {
	users, err := r.client.GetOrgUsers()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to list org users: %v", err)), nil
	}
	return jsonResult(users)
}

func (r *Registry) handleGetCurrentUser(args map[string]interface{}) (*mcp.CallToolResult, error) {
	user, err := r.client.GetCurrentUser()
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get current user: %v", err)), nil
	}
	return jsonResult(user)
}

func (r *Registry) handleListTeams(args map[string]interface{}) (*mcp.CallToolResult, error) {
	query := getString(args, "query")
	page := getInt(args, "page")
	perPage := getInt(args, "per_page")

	teams, err := r.client.GetTeams(query, page, perPage)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to list teams: %v", err)), nil
	}
	return jsonResult(teams)
}

func (r *Registry) handleGetTeam(args map[string]interface{}) (*mcp.CallToolResult, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return errorResult("id is required"), nil
	}

	team, err := r.client.GetTeam(id)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to get team: %v", err)), nil
	}
	return jsonResult(team)
}

func (r *Registry) handleCreateTeam(args map[string]interface{}) (*mcp.CallToolResult, error) {
	name := getString(args, "name")
	if name == "" {
		return errorResult("name is required"), nil
	}

	team, err := r.client.CreateTeam(name, getString(args, "email"))
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to create team: %v", err)), nil
	}
	return jsonResult(team)
}

func (r *Registry) handleDeleteTeam(args map[string]interface{}) (*mcp.CallToolResult, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return errorResult("id is required"), nil
	}

	if err := r.client.DeleteTeam(id); err != nil {
		return errorResult(fmt.Sprintf("Failed to delete team: %v", err)), nil
	}
	return jsonResult(map[string]interface{}{"status": "deleted", "id": id})
}
