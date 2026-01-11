# Grafana MCP Server

A comprehensive Model Context Protocol (MCP) server for Grafana, enabling AI assistants to interact with Grafana's full API surface.

## Features

**38 Tools** covering:

### Dashboards
- `grafana_search_dashboards` - Search dashboards by query, tags, folder
- `grafana_get_dashboard` - Get dashboard by UID with full configuration
- `grafana_create_dashboard` - Create new dashboards with panels
- `grafana_update_dashboard` - Update existing dashboards
- `grafana_delete_dashboard` - Delete dashboards

### Datasources
- `grafana_list_datasources` - List all configured datasources
- `grafana_get_datasource` - Get datasource details by UID
- `grafana_create_datasource` - Create new datasources (Prometheus, Loki, etc.)
- `grafana_update_datasource` - Update datasource configuration
- `grafana_delete_datasource` - Delete datasources

### Folders
- `grafana_list_folders` - List all folders
- `grafana_get_folder` - Get folder by UID
- `grafana_create_folder` - Create new folders
- `grafana_update_folder` - Update folder title
- `grafana_delete_folder` - Delete folders (and contents)

### Alert Rules
- `grafana_list_alert_rules` - List all alert rules
- `grafana_get_alert_rule` - Get alert rule by UID
- `grafana_create_alert_rule` - Create new alert rules
- `grafana_update_alert_rule` - Update alert rules
- `grafana_delete_alert_rule` - Delete alert rules

### Annotations
- `grafana_list_annotations` - List annotations with filters
- `grafana_create_annotation` - Create annotations (point or region)
- `grafana_update_annotation` - Update existing annotations
- `grafana_delete_annotation` - Delete annotations

### Query
- `grafana_query` - Execute queries against any datasource (PromQL, LogQL, etc.)

### Organization
- `grafana_get_org` - Get current organization info
- `grafana_list_org_users` - List users in organization

### Users
- `grafana_get_current_user` - Get authenticated user info

### Teams
- `grafana_list_teams` - List teams with search
- `grafana_get_team` - Get team by ID
- `grafana_create_team` - Create new teams
- `grafana_delete_team` - Delete teams

### Health
- `grafana_health` - Check Grafana server health and version

## Installation

### Build from source

```bash
go build -o grafana-mcp-server ./cmd/server
```

### Run

```bash
export GRAFANA_URL="http://localhost:3000"
export GRAFANA_API_KEY="your-api-key-here"
./grafana-mcp-server
```

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `GRAFANA_URL` | Grafana server URL | `http://localhost:3000` |
| `GRAFANA_API_KEY` | Service account token or API key | (none) |

### Creating an API Key

1. In Grafana, go to **Administration** → **Service accounts**
2. Create a new service account
3. Add a token with appropriate permissions
4. Use the token as `GRAFANA_API_KEY`

## Claude Desktop Configuration

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "/path/to/grafana-mcp-server",
      "env": {
        "GRAFANA_URL": "http://localhost:3000",
        "GRAFANA_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Example Usage

### Search dashboards
```json
{
  "name": "grafana_search_dashboards",
  "arguments": {
    "query": "kubernetes",
    "tags": ["production"],
    "limit": 10
  }
}
```

### Execute PromQL query
```json
{
  "name": "grafana_query",
  "arguments": {
    "datasource_uid": "prometheus",
    "datasource_type": "prometheus",
    "query": "rate(http_requests_total[5m])",
    "from": "now-1h",
    "to": "now"
  }
}
```

### Create annotation
```json
{
  "name": "grafana_create_annotation",
  "arguments": {
    "text": "Deployment v2.0.0",
    "tags": ["deployment", "production"],
    "dashboard_uid": "abc123"
  }
}
```

## Tool Annotations

All tools include MCP annotations for safe operation:

- `readOnlyHint: true` - Read operations (list, get, search)
- `destructiveHint: true` - Delete and update operations
- `openWorldHint: true` - External API interactions

## License

MIT
