package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a Grafana API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Grafana client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request to the Grafana API
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Provide user-friendly error for connection failures
		if strings.Contains(err.Error(), "connection refused") {
			return nil, fmt.Errorf("cannot connect to Grafana at %s: connection refused. Ensure Grafana is running and accessible", c.baseURL)
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, fmt.Errorf("cannot connect to Grafana at %s: host not found. Check GRAFANA_URL configuration", c.baseURL)
		}
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("connection to Grafana at %s timed out. Check network connectivity", c.baseURL)
		}
		return nil, fmt.Errorf("request to Grafana failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ============== Dashboard Operations ==============

// Dashboard represents a Grafana dashboard
type Dashboard struct {
	ID          int64             `json:"id,omitempty"`
	UID         string            `json:"uid,omitempty"`
	Title       string            `json:"title"`
	Tags        []string          `json:"tags,omitempty"`
	Timezone    string            `json:"timezone,omitempty"`
	SchemaVersion int             `json:"schemaVersion,omitempty"`
	Version     int               `json:"version,omitempty"`
	Panels      []Panel           `json:"panels,omitempty"`
	Templating  *Templating       `json:"templating,omitempty"`
	Annotations *AnnotationConfig `json:"annotations,omitempty"`
	Refresh     string            `json:"refresh,omitempty"`
	Time        *TimeRange        `json:"time,omitempty"`
}

type Panel struct {
	ID          int64              `json:"id,omitempty"`
	Type        string             `json:"type"`
	Title       string             `json:"title"`
	Description string             `json:"description,omitempty"`
	GridPos     GridPos            `json:"gridPos"`
	Targets     []Target           `json:"targets,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	FieldConfig *FieldConfig       `json:"fieldConfig,omitempty"`
}

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Target struct {
	RefID      string `json:"refId"`
	Expr       string `json:"expr,omitempty"`
	Query      string `json:"query,omitempty"`
	Datasource *DatasourceRef `json:"datasource,omitempty"`
}

type DatasourceRef struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type FieldConfig struct {
	Defaults  map[string]interface{} `json:"defaults,omitempty"`
	Overrides []interface{}          `json:"overrides,omitempty"`
}

type Templating struct {
	List []TemplateVar `json:"list,omitempty"`
}

type TemplateVar struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Query   string `json:"query,omitempty"`
	Current map[string]interface{} `json:"current,omitempty"`
}

type AnnotationConfig struct {
	List []AnnotationQuery `json:"list,omitempty"`
}

type AnnotationQuery struct {
	Name       string `json:"name"`
	Enable     bool   `json:"enable"`
	Datasource *DatasourceRef `json:"datasource,omitempty"`
}

type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// SaveDashboardRequest represents a request to save a dashboard
type SaveDashboardRequest struct {
	Dashboard Dashboard `json:"dashboard"`
	FolderUID string    `json:"folderUid,omitempty"`
	Message   string    `json:"message,omitempty"`
	Overwrite bool      `json:"overwrite,omitempty"`
}

// SaveDashboardResponse represents the response from saving a dashboard
type SaveDashboardResponse struct {
	ID      int64  `json:"id"`
	UID     string `json:"uid"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Version int    `json:"version"`
}

// SearchDashboardsResponse represents search results
type SearchDashboardsResponse struct {
	ID          int64    `json:"id"`
	UID         string   `json:"uid"`
	Title       string   `json:"title"`
	URI         string   `json:"uri"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	IsStarred   bool     `json:"isStarred"`
	FolderID    int64    `json:"folderId"`
	FolderUID   string   `json:"folderUid"`
	FolderTitle string   `json:"folderTitle"`
}

// SearchDashboards searches for dashboards
func (c *Client) SearchDashboards(query string, tags []string, folderIDs []int64, dashboardType string, limit int) ([]SearchDashboardsResponse, error) {
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	for _, tag := range tags {
		params.Add("tag", tag)
	}
	for _, fid := range folderIDs {
		params.Add("folderIds", fmt.Sprintf("%d", fid))
	}
	if dashboardType != "" {
		params.Set("type", dashboardType)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	path := "/api/search"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var results []SearchDashboardsResponse
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// GetDashboard retrieves a dashboard by UID
func (c *Client) GetDashboard(uid string) (*Dashboard, error) {
	resp, err := c.doRequest("GET", "/api/dashboards/uid/"+uid, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Dashboard Dashboard `json:"dashboard"`
		Meta      struct {
			FolderUID string `json:"folderUid"`
			URL       string `json:"url"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.Dashboard, nil
}

// SaveDashboard creates or updates a dashboard
func (c *Client) SaveDashboard(req SaveDashboardRequest) (*SaveDashboardResponse, error) {
	resp, err := c.doRequest("POST", "/api/dashboards/db", req)
	if err != nil {
		return nil, err
	}

	var result SaveDashboardResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteDashboard deletes a dashboard by UID
func (c *Client) DeleteDashboard(uid string) error {
	_, err := c.doRequest("DELETE", "/api/dashboards/uid/"+uid, nil)
	return err
}

// ============== Datasource Operations ==============

// Datasource represents a Grafana datasource
type Datasource struct {
	ID        int64                  `json:"id,omitempty"`
	UID       string                 `json:"uid,omitempty"`
	OrgID     int64                  `json:"orgId,omitempty"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	TypeName  string                 `json:"typeName,omitempty"`
	Access    string                 `json:"access"`
	URL       string                 `json:"url"`
	User      string                 `json:"user,omitempty"`
	Database  string                 `json:"database,omitempty"`
	BasicAuth bool                   `json:"basicAuth,omitempty"`
	IsDefault bool                   `json:"isDefault,omitempty"`
	JSONData  map[string]interface{} `json:"jsonData,omitempty"`
	SecureJSONData map[string]string `json:"secureJsonData,omitempty"`
	ReadOnly  bool                   `json:"readOnly,omitempty"`
}

// GetDatasources retrieves all datasources
func (c *Client) GetDatasources() ([]Datasource, error) {
	resp, err := c.doRequest("GET", "/api/datasources", nil)
	if err != nil {
		return nil, err
	}

	var results []Datasource
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// GetDatasource retrieves a datasource by UID
func (c *Client) GetDatasource(uid string) (*Datasource, error) {
	resp, err := c.doRequest("GET", "/api/datasources/uid/"+uid, nil)
	if err != nil {
		return nil, err
	}

	var result Datasource
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateDatasource creates a new datasource
func (c *Client) CreateDatasource(ds Datasource) (*Datasource, error) {
	resp, err := c.doRequest("POST", "/api/datasources", ds)
	if err != nil {
		return nil, err
	}

	var result struct {
		Datasource Datasource `json:"datasource"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.Datasource, nil
}

// UpdateDatasource updates an existing datasource
func (c *Client) UpdateDatasource(uid string, ds Datasource) (*Datasource, error) {
	resp, err := c.doRequest("PUT", "/api/datasources/uid/"+uid, ds)
	if err != nil {
		return nil, err
	}

	var result struct {
		Datasource Datasource `json:"datasource"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.Datasource, nil
}

// DeleteDatasource deletes a datasource by UID
func (c *Client) DeleteDatasource(uid string) error {
	_, err := c.doRequest("DELETE", "/api/datasources/uid/"+uid, nil)
	return err
}

// ============== Folder Operations ==============

// Folder represents a Grafana folder
type Folder struct {
	ID        int64  `json:"id,omitempty"`
	UID       string `json:"uid,omitempty"`
	Title     string `json:"title"`
	URL       string `json:"url,omitempty"`
	HasACL    bool   `json:"hasAcl,omitempty"`
	CanSave   bool   `json:"canSave,omitempty"`
	CanEdit   bool   `json:"canEdit,omitempty"`
	CanAdmin  bool   `json:"canAdmin,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
	Created   string `json:"created,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
	Updated   string `json:"updated,omitempty"`
	Version   int    `json:"version,omitempty"`
}

// GetFolders retrieves all folders
func (c *Client) GetFolders() ([]Folder, error) {
	resp, err := c.doRequest("GET", "/api/folders", nil)
	if err != nil {
		return nil, err
	}

	var results []Folder
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// GetFolder retrieves a folder by UID
func (c *Client) GetFolder(uid string) (*Folder, error) {
	resp, err := c.doRequest("GET", "/api/folders/"+uid, nil)
	if err != nil {
		return nil, err
	}

	var result Folder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateFolder creates a new folder
func (c *Client) CreateFolder(title, uid string) (*Folder, error) {
	body := map[string]string{"title": title}
	if uid != "" {
		body["uid"] = uid
	}

	resp, err := c.doRequest("POST", "/api/folders", body)
	if err != nil {
		return nil, err
	}

	var result Folder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateFolder updates a folder
func (c *Client) UpdateFolder(uid, title string, version int) (*Folder, error) {
	body := map[string]interface{}{
		"title":   title,
		"version": version,
	}

	resp, err := c.doRequest("PUT", "/api/folders/"+uid, body)
	if err != nil {
		return nil, err
	}

	var result Folder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteFolder deletes a folder by UID
func (c *Client) DeleteFolder(uid string) error {
	_, err := c.doRequest("DELETE", "/api/folders/"+uid, nil)
	return err
}

// ============== Alert Operations ==============

// AlertRule represents a Grafana alert rule
type AlertRule struct {
	ID           int64                  `json:"id,omitempty"`
	UID          string                 `json:"uid,omitempty"`
	OrgID        int64                  `json:"orgId,omitempty"`
	FolderUID    string                 `json:"folderUID"`
	RuleGroup    string                 `json:"ruleGroup"`
	Title        string                 `json:"title"`
	Condition    string                 `json:"condition"`
	Data         []AlertQuery           `json:"data"`
	NoDataState  string                 `json:"noDataState,omitempty"`
	ExecErrState string                 `json:"execErrState,omitempty"`
	For          string                 `json:"for,omitempty"`
	Annotations  map[string]string      `json:"annotations,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
	IsPaused     bool                   `json:"isPaused,omitempty"`
}

type AlertQuery struct {
	RefID             string                 `json:"refId"`
	QueryType         string                 `json:"queryType,omitempty"`
	RelativeTimeRange RelativeTimeRange      `json:"relativeTimeRange,omitempty"`
	DatasourceUID     string                 `json:"datasourceUid"`
	Model             map[string]interface{} `json:"model"`
}

type RelativeTimeRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

// GetAlertRules retrieves all alert rules
func (c *Client) GetAlertRules() ([]AlertRule, error) {
	resp, err := c.doRequest("GET", "/api/v1/provisioning/alert-rules", nil)
	if err != nil {
		return nil, err
	}

	var results []AlertRule
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// GetAlertRule retrieves an alert rule by UID
func (c *Client) GetAlertRule(uid string) (*AlertRule, error) {
	resp, err := c.doRequest("GET", "/api/v1/provisioning/alert-rules/"+uid, nil)
	if err != nil {
		return nil, err
	}

	var result AlertRule
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateAlertRule creates a new alert rule
func (c *Client) CreateAlertRule(rule AlertRule) (*AlertRule, error) {
	resp, err := c.doRequest("POST", "/api/v1/provisioning/alert-rules", rule)
	if err != nil {
		return nil, err
	}

	var result AlertRule
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAlertRule updates an alert rule
func (c *Client) UpdateAlertRule(uid string, rule AlertRule) (*AlertRule, error) {
	resp, err := c.doRequest("PUT", "/api/v1/provisioning/alert-rules/"+uid, rule)
	if err != nil {
		return nil, err
	}

	var result AlertRule
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteAlertRule deletes an alert rule by UID
func (c *Client) DeleteAlertRule(uid string) error {
	_, err := c.doRequest("DELETE", "/api/v1/provisioning/alert-rules/"+uid, nil)
	return err
}

// ============== Annotation Operations ==============

// Annotation represents a Grafana annotation
type Annotation struct {
	ID          int64    `json:"id,omitempty"`
	AlertID     int64    `json:"alertId,omitempty"`
	DashboardID int64    `json:"dashboardId,omitempty"`
	DashboardUID string  `json:"dashboardUID,omitempty"`
	PanelID     int64    `json:"panelId,omitempty"`
	Time        int64    `json:"time"`
	TimeEnd     int64    `json:"timeEnd,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Text        string   `json:"text"`
	Type        string   `json:"type,omitempty"`
	Created     int64    `json:"created,omitempty"`
	Updated     int64    `json:"updated,omitempty"`
	UserID      int64    `json:"userId,omitempty"`
	UserName    string   `json:"userName,omitempty"`
	UserEmail   string   `json:"email,omitempty"`
	UserAvatarURL string `json:"avatarUrl,omitempty"`
}

// GetAnnotations retrieves annotations with optional filters
func (c *Client) GetAnnotations(from, to int64, dashboardUID string, panelID int64, tags []string, limit int) ([]Annotation, error) {
	params := url.Values{}
	if from > 0 {
		params.Set("from", fmt.Sprintf("%d", from))
	}
	if to > 0 {
		params.Set("to", fmt.Sprintf("%d", to))
	}
	if dashboardUID != "" {
		params.Set("dashboardUID", dashboardUID)
	}
	if panelID > 0 {
		params.Set("panelId", fmt.Sprintf("%d", panelID))
	}
	for _, tag := range tags {
		params.Add("tags", tag)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	path := "/api/annotations"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var results []Annotation
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// CreateAnnotation creates a new annotation
func (c *Client) CreateAnnotation(ann Annotation) (*Annotation, error) {
	resp, err := c.doRequest("POST", "/api/annotations", ann)
	if err != nil {
		return nil, err
	}

	var result struct {
		ID      int64  `json:"id"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	ann.ID = result.ID
	return &ann, nil
}

// UpdateAnnotation updates an annotation
func (c *Client) UpdateAnnotation(id int64, ann Annotation) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/annotations/%d", id), ann)
	return err
}

// DeleteAnnotation deletes an annotation by ID
func (c *Client) DeleteAnnotation(id int64) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/annotations/%d", id), nil)
	return err
}

// ============== Organization Operations ==============

// Organization represents a Grafana organization
type Organization struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"name"`
	Address Address `json:"address,omitempty"`
}

type Address struct {
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	City     string `json:"city,omitempty"`
	ZipCode  string `json:"zipCode,omitempty"`
	State    string `json:"state,omitempty"`
	Country  string `json:"country,omitempty"`
}

// GetCurrentOrg retrieves the current organization
func (c *Client) GetCurrentOrg() (*Organization, error) {
	resp, err := c.doRequest("GET", "/api/org", nil)
	if err != nil {
		return nil, err
	}

	var result Organization
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== User Operations ==============

// User represents a Grafana user
type User struct {
	ID             int64  `json:"id,omitempty"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Login          string `json:"login"`
	Theme          string `json:"theme,omitempty"`
	OrgID          int64  `json:"orgId,omitempty"`
	IsGrafanaAdmin bool   `json:"isGrafanaAdmin,omitempty"`
	IsDisabled     bool   `json:"isDisabled,omitempty"`
	IsExternal     bool   `json:"isExternal,omitempty"`
	AuthLabels     []string `json:"authLabels,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	AvatarURL      string `json:"avatarUrl,omitempty"`
}

// GetCurrentUser retrieves the current user
func (c *Client) GetCurrentUser() (*User, error) {
	resp, err := c.doRequest("GET", "/api/user", nil)
	if err != nil {
		return nil, err
	}

	var result User
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetOrgUsers retrieves users in the current organization
func (c *Client) GetOrgUsers() ([]User, error) {
	resp, err := c.doRequest("GET", "/api/org/users", nil)
	if err != nil {
		return nil, err
	}

	var results []User
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// ============== Query Operations ==============

// QueryRequest represents a query request
type QueryRequest struct {
	From    string        `json:"from"`
	To      string        `json:"to"`
	Queries []QueryTarget `json:"queries"`
}

type QueryTarget struct {
	RefID         string                 `json:"refId"`
	Datasource    DatasourceRef          `json:"datasource"`
	MaxDataPoints int                    `json:"maxDataPoints,omitempty"`
	IntervalMs    int                    `json:"intervalMs,omitempty"`
	Query         string                 `json:"expr,omitempty"`
	RawQuery      string                 `json:"rawQuery,omitempty"`
	QueryType     string                 `json:"queryType,omitempty"`
	Extra         map[string]interface{} `json:"-"`
}

// QueryResponse represents query results
type QueryResponse struct {
	Results map[string]QueryResult `json:"results"`
}

type QueryResult struct {
	Frames []DataFrame `json:"frames,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type DataFrame struct {
	Schema FrameSchema `json:"schema,omitempty"`
	Data   FrameData   `json:"data,omitempty"`
}

type FrameSchema struct {
	Name   string        `json:"name,omitempty"`
	Fields []FieldSchema `json:"fields,omitempty"`
}

type FieldSchema struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels,omitempty"`
}

type FrameData struct {
	Values [][]interface{} `json:"values,omitempty"`
}

// Query executes a query against datasources
func (c *Client) Query(req QueryRequest) (*QueryResponse, error) {
	resp, err := c.doRequest("POST", "/api/ds/query", req)
	if err != nil {
		return nil, err
	}

	var result QueryResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== Health Operations ==============

// Health represents Grafana health status
type Health struct {
	Commit   string `json:"commit"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

// GetHealth retrieves the health status
func (c *Client) GetHealth() (*Health, error) {
	resp, err := c.doRequest("GET", "/api/health", nil)
	if err != nil {
		return nil, err
	}

	var result Health
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== Team Operations ==============

// Team represents a Grafana team
type Team struct {
	ID          int64  `json:"id,omitempty"`
	OrgID       int64  `json:"orgId,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	MemberCount int    `json:"memberCount,omitempty"`
	Permission  int    `json:"permission,omitempty"`
}

// GetTeams retrieves all teams
func (c *Client) GetTeams(query string, page, perPage int) ([]Team, error) {
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		params.Set("perpage", fmt.Sprintf("%d", perPage))
	}

	path := "/api/teams/search"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Teams []Team `json:"teams"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.Teams, nil
}

// GetTeam retrieves a team by ID
func (c *Client) GetTeam(id int64) (*Team, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/teams/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var result Team
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateTeam creates a new team
func (c *Client) CreateTeam(name, email string) (*Team, error) {
	body := map[string]string{"name": name}
	if email != "" {
		body["email"] = email
	}

	resp, err := c.doRequest("POST", "/api/teams", body)
	if err != nil {
		return nil, err
	}

	var result struct {
		TeamID  int64  `json:"teamId"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &Team{ID: result.TeamID, Name: name, Email: email}, nil
}

// DeleteTeam deletes a team by ID
func (c *Client) DeleteTeam(id int64) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/teams/%d", id), nil)
	return err
}
