# Grafana MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that exposes the Grafana HTTP API as tools for AI assistants like Claude. Manage dashboards, datasources, folders, alert rules, annotations, teams, and more — all through natural language.

**34 tools across 9 Grafana API domains.**

> No external dependencies beyond the YAML config library — pure Go stdlib for all Grafana API communication.

---

## Requirements

- Go 1.22+
- A running Grafana instance (self-hosted or Grafana Cloud)
- A Grafana API key or service account token

---

## Build

```bash
go build -o bin/grafana-mcp ./cmd/server
```

Or use the Makefile:

```bash
make build          # current platform
make build-darwin   # macOS arm64
make build-linux    # Linux amd64
```

---

## Configuration

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `GRAFANA_URL` | `http://localhost:3000` | Grafana base URL |
| `GRAFANA_API_KEY` | — | API key or service account token |
| `GRAFANA_CONFIG_FILE` | `config.yaml` | Path to tool enable/disable config |

### Tool configuration (optional)

Tools are individually enabled or disabled via a YAML file. By default every tool is enabled. The file path is resolved in this order:

1. `GRAFANA_CONFIG_FILE` environment variable
2. `config.yaml` in the working directory
3. No file → all tools enabled

**Format:**

```yaml
tools:
  tool_name:
    enabled: false   # omit or set true to enable
```

See [Recommended Profiles](#recommended-configuration-profiles) for ready-to-use configurations.

---

## Running with Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "/path/to/bin/grafana-mcp",
      "env": {
        "GRAFANA_URL": "https://your-grafana.example.com",
        "GRAFANA_API_KEY": "glsa_YOUR_SERVICE_ACCOUNT_TOKEN",
        "GRAFANA_CONFIG_FILE": "/path/to/config.yaml"
      }
    }
  }
}
```

## Running with Claude Code (CLI)

Add to `~/.mcp.json`:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "/path/to/bin/grafana-mcp",
      "env": {
        "GRAFANA_URL": "https://your-grafana.example.com",
        "GRAFANA_API_KEY": "glsa_YOUR_SERVICE_ACCOUNT_TOKEN",
        "GRAFANA_CONFIG_FILE": "/path/to/config.yaml"
      }
    }
  }
}
```

---

## Tool Domains

### Health (1 tool)
| Tool | Description |
|---|---|
| `grafana_health` | Check Grafana server health and version |

### Dashboards (5 tools)
| Tool | Description |
|---|---|
| `grafana_search_dashboards` | Search dashboards by query, folder, or tags |
| `grafana_get_dashboard` | Get a dashboard by UID |
| `grafana_create_dashboard` | Create a new dashboard |
| `grafana_update_dashboard` | Update an existing dashboard |
| `grafana_delete_dashboard` | Delete a dashboard by UID |

### Datasources (5 tools)
| Tool | Description |
|---|---|
| `grafana_list_datasources` | List all configured datasources |
| `grafana_get_datasource` | Get a datasource by name or ID |
| `grafana_create_datasource` | Add a new datasource |
| `grafana_update_datasource` | Update a datasource configuration |
| `grafana_delete_datasource` | Remove a datasource |

### Folders (5 tools)
| Tool | Description |
|---|---|
| `grafana_list_folders` | List all dashboard folders |
| `grafana_get_folder` | Get a folder by UID |
| `grafana_create_folder` | Create a new folder |
| `grafana_update_folder` | Rename or move a folder |
| `grafana_delete_folder` | Delete a folder |

### Alert Rules (5 tools)
| Tool | Description |
|---|---|
| `grafana_list_alert_rules` | List all alert rules |
| `grafana_get_alert_rule` | Get an alert rule by UID |
| `grafana_create_alert_rule` | Create a new alert rule |
| `grafana_update_alert_rule` | Update an existing alert rule |
| `grafana_delete_alert_rule` | Delete an alert rule |

### Annotations (4 tools)
| Tool | Description |
|---|---|
| `grafana_list_annotations` | List annotations with optional time range and tag filters |
| `grafana_create_annotation` | Create a new annotation |
| `grafana_update_annotation` | Update an existing annotation |
| `grafana_delete_annotation` | Delete an annotation |

### Query (1 tool)
| Tool | Description |
|---|---|
| `grafana_query` | Execute a raw datasource query (PromQL, Loki, etc.) |

### Organization (2 tools)
| Tool | Description |
|---|---|
| `grafana_get_org` | Get current organization info |
| `grafana_list_org_users` | List users in the current organization |

### User (1 tool)
| Tool | Description |
|---|---|
| `grafana_get_current_user` | Get the currently authenticated user |

### Teams (4 tools)
| Tool | Description |
|---|---|
| `grafana_list_teams` | List all teams |
| `grafana_get_team` | Get a team by ID |
| `grafana_create_team` | Create a new team |
| `grafana_delete_team` | Delete a team |

---

## Recommended Configuration Profiles

Copy the relevant block into your `config.yaml` (or a file pointed to by `GRAFANA_CONFIG_FILE`).

### Read-Only

All list/get/search/query/health tools active. Every create, update, and delete tool disabled. Safe for shared environments and dashboards you don't want accidentally modified.

```yaml
# config-readonly.yaml
# Zero writes — safe for production read access.
tools:
  grafana_create_dashboard:
    enabled: false
  grafana_update_dashboard:
    enabled: false
  grafana_delete_dashboard:
    enabled: false
  grafana_create_datasource:
    enabled: false
  grafana_update_datasource:
    enabled: false
  grafana_delete_datasource:
    enabled: false
  grafana_create_folder:
    enabled: false
  grafana_update_folder:
    enabled: false
  grafana_delete_folder:
    enabled: false
  grafana_create_alert_rule:
    enabled: false
  grafana_update_alert_rule:
    enabled: false
  grafana_delete_alert_rule:
    enabled: false
  grafana_create_annotation:
    enabled: false
  grafana_update_annotation:
    enabled: false
  grafana_delete_annotation:
    enabled: false
  grafana_create_team:
    enabled: false
  grafana_delete_team:
    enabled: false
```

---

### SRE / On-Call

Observability focus: health checks, dashboard browsing, querying datasources, reading alert rules and annotations. No configuration changes or deletions. Annotations can be created (useful for marking incidents).

```yaml
# config-sre.yaml
# Read + incident annotations. No config mutations.
tools:
  grafana_create_dashboard:
    enabled: false
  grafana_update_dashboard:
    enabled: false
  grafana_delete_dashboard:
    enabled: false
  grafana_create_datasource:
    enabled: false
  grafana_update_datasource:
    enabled: false
  grafana_delete_datasource:
    enabled: false
  grafana_create_folder:
    enabled: false
  grafana_update_folder:
    enabled: false
  grafana_delete_folder:
    enabled: false
  grafana_create_alert_rule:
    enabled: false
  grafana_update_alert_rule:
    enabled: false
  grafana_delete_alert_rule:
    enabled: false
  grafana_update_annotation:
    enabled: false
  grafana_delete_annotation:
    enabled: false
  grafana_create_team:
    enabled: false
  grafana_delete_team:
    enabled: false
```

---

### Dashboard Author

Full dashboard and folder management. Cannot touch datasources, alert rules, or teams — reducing the blast radius to the dashboard layer only.

```yaml
# config-dashboard-author.yaml
# Dashboard + folder CRUD. Read-only for everything else.
tools:
  grafana_create_datasource:
    enabled: false
  grafana_update_datasource:
    enabled: false
  grafana_delete_datasource:
    enabled: false
  grafana_create_alert_rule:
    enabled: false
  grafana_update_alert_rule:
    enabled: false
  grafana_delete_alert_rule:
    enabled: false
  grafana_create_annotation:
    enabled: false
  grafana_update_annotation:
    enabled: false
  grafana_delete_annotation:
    enabled: false
  grafana_create_team:
    enabled: false
  grafana_delete_team:
    enabled: false
```

---

### Alerting Engineer

Full alert rule management plus read access to dashboards and datasources for context. No dashboard mutations, no team or folder management.

```yaml
# config-alerting.yaml
# Alert rule CRUD + read access. No dashboard/folder/team mutations.
tools:
  grafana_create_dashboard:
    enabled: false
  grafana_update_dashboard:
    enabled: false
  grafana_delete_dashboard:
    enabled: false
  grafana_create_datasource:
    enabled: false
  grafana_update_datasource:
    enabled: false
  grafana_delete_datasource:
    enabled: false
  grafana_create_folder:
    enabled: false
  grafana_update_folder:
    enabled: false
  grafana_delete_folder:
    enabled: false
  grafana_create_team:
    enabled: false
  grafana_delete_team:
    enabled: false
```

---

### Admin

All tools enabled. Default when no config file is present.

```yaml
# config-admin.yaml
# Full access — all 34 tools enabled.
tools: {}
```

---

## Grafana API Key Scopes

The minimum required permissions for a service account token depend on which tools you enable:

| Domain | Required role / permission |
|---|---|
| Health | No auth required (public endpoint) |
| Dashboards | `Viewer` to read; `Editor` to create/update/delete |
| Datasources | `Viewer` to list/get; `Admin` to create/update/delete |
| Folders | `Viewer` to list/get; `Editor` to create/update/delete |
| Alert Rules | `Viewer` to read; `Editor` to create/update/delete |
| Annotations | `Viewer` to read; `Editor` to create/update/delete |
| Query | `Viewer` (datasource query permissions apply) |
| Organization | `Viewer` |
| Teams | `Viewer` to read; `Admin` to create/delete |

For read-only profiles a **Viewer** service account is sufficient. For full admin profiles use an **Admin** service account or a token with `Admin` role.

---

## Project Structure

```
.
├── cmd/server/main.go          # Entry point — env config, MCP protocol loop
├── config.yaml                 # Tool enable/disable configuration
├── internal/
│   ├── config/config.go        # ToolsConfig, IsEnabled(), YAML loading
│   ├── grafana/client.go       # Grafana HTTP client (all API calls)
│   ├── mcp/types.go            # MCP JSON-RPC types
│   └── tools/registry.go       # Tool registry, definitions, and handlers
├── Makefile
└── go.mod
```
