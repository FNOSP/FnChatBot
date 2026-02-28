package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ToolTypeFunction = "function"
)

type Tool struct {
	Type     string     `json:"type"`
	Function ToolSchema `json:"function"`
}

type ToolSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ToolService struct {
	UserID uint
}

// NewToolService creates a ToolService scoped to a specific user.
func NewToolService(userID uint) *ToolService {
	return &ToolService{UserID: userID}
}

func (s *ToolService) GetAvailableTools() ([]Tool, error) {
	var tools []Tool

	tools = append(tools, Tool{
		Type: ToolTypeFunction,
		Function: ToolSchema{
			Name:        "TodoWrite",
			Description: "Update the task list. Use to plan and track progress. ALWAYS call this when starting a multi-step task.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"items": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"content":    map[string]interface{}{"type": "string", "description": "Task description"},
								"status":     map[string]interface{}{"type": "string", "enum": []string{"pending", "in_progress", "completed"}},
								"activeForm": map[string]interface{}{"type": "string", "description": "Present tense action, e.g. 'Reading files'"},
							},
							"required": []string{"content", "status", "activeForm"},
						},
					},
				},
				"required": []string{"items"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: ToolTypeFunction,
		Function: ToolSchema{
			Name:        "Task",
			Description: "Delegate a sub-task to a specialized agent. Use this for complex steps like 'explore codebase' or 'write detailed plan'.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"description":   map[string]interface{}{"type": "string", "description": "Short description of the task"},
					"prompt":        map[string]interface{}{"type": "string", "description": "Detailed instructions for the subagent"},
					"subagent_type": map[string]interface{}{"type": "string", "enum": []string{"explore", "code", "plan"}},
				},
				"required": []string{"description", "prompt", "subagent_type"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: ToolTypeFunction,
		Function: ToolSchema{
			Name:        "Skill",
			Description: "Load a specialized skill/knowledge. Use this when you need domain expertise (e.g. 'how to review code', 'how to build mcp').",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string", "description": "Name of the skill to load"},
				},
				"required": []string{"name"},
			},
		},
	})

	var skills []models.Skill
	if err := db.DB.Where("enabled = ? AND user_id = ?", true, s.UserID).Find(&skills).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch skills: %v", err)
	}

	for _, skill := range skills {
		tool, err := s.convertSkillToTool(skill)
		if err != nil {
			log.Printf("Skipping invalid skill %s: %v", skill.Name, err)
			continue
		}
		tools = append(tools, tool)
	}

	// Collect tools from connected MCP clients (config in mcp.json, no user filter)
	if DefaultMCPService != nil {
		clients := DefaultMCPService.GetConnectedClients()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for name, c := range clients {
			mcpTools, err := s.listMCPTools(ctx, name, c)
			if err != nil {
				log.Printf("Failed to list tools from MCP %s: %v", name, err)
				continue
			}
			tools = append(tools, mcpTools...)
		}
	}

	return tools, nil
}

func (s *ToolService) convertSkillToTool(skill models.Skill) (Tool, error) {
	var configMap map[string]interface{}
	if err := json.Unmarshal(skill.Config, &configMap); err != nil {
		return Tool{}, fmt.Errorf("invalid config json: %v", err)
	}

	params, ok := configMap["parameters"]
	if !ok {
		params = map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		paramsMap = map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}
	}

	return Tool{
		Type: ToolTypeFunction,
		Function: ToolSchema{
			Name:        skill.Name,
			Description: skill.Description,
			Parameters:  paramsMap,
		},
	}, nil
}

// listMCPTools calls MCP ListTools and converts result to our Tool slice.
func (s *ToolService) listMCPTools(ctx context.Context, serverName string, c *client.Client) ([]Tool, error) {
	res, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	out := make([]Tool, 0, len(res.Tools))
	for _, t := range res.Tools {
		params := make(map[string]interface{})
		data, _ := json.Marshal(t.InputSchema)
		_ = json.Unmarshal(data, &params)
		if len(params) == 0 {
			params = map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}
		}
		out = append(out, Tool{
			Type: ToolTypeFunction,
			Function: ToolSchema{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}
	return out, nil
}

func (s *ToolService) ExecuteSkill(name string, args string) (string, error) {
	if name == "TodoWrite" {
		return fmt.Sprintf("Tasks updated. Current state: %s", args), nil
	}

	if name == "Task" {
		var taskArgs struct {
			Description string `json:"description"`
			Prompt      string `json:"prompt"`
			Subagent    string `json:"subagent_type"`
		}
		if err := json.Unmarshal([]byte(args), &taskArgs); err != nil {
			return "", fmt.Errorf("invalid task args: %v", err)
		}

		if taskArgs.Subagent == "explore" {
			return fmt.Sprintf("Subagent [explore] completed task: %s.\nFindings:\n- Found project root at /app\n- Found go.mod (go 1.21)\n- Found main.go\nAnalysis complete.", taskArgs.Description), nil
		}
		return fmt.Sprintf("Subagent [%s] executed task: %s. Result: Success.", taskArgs.Subagent, taskArgs.Description), nil
	}

	if name == "Skill" {
		var skillArgs struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal([]byte(args), &skillArgs); err != nil {
			return "", fmt.Errorf("invalid skill args: %v", err)
		}

		return fmt.Sprintf(`<skill-loaded name="%s">
# Skill: %s

## Best Practices
1. Keep functions small and focused.
2. Use descriptive variable names.
3. Handle errors explicitly.

## Common Patterns
- Repository Pattern
- Service Layer
- Dependency Injection
</skill-loaded>

Skill loaded successfully. You can now use this knowledge to assist the user.`, skillArgs.Name, skillArgs.Name), nil
	}

	if name == "get_current_time" {
		return "2023-10-27 10:00:00", nil
	}

	// Execute via MCP client CallTool
	if DefaultMCPService != nil {
		clients := DefaultMCPService.GetConnectedClients()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		var argsMap map[string]interface{}
		if args != "" {
			_ = json.Unmarshal([]byte(args), &argsMap)
		}
		if argsMap == nil {
			argsMap = make(map[string]interface{})
		}
		req := mcp.CallToolRequest{Params: mcp.CallToolParams{Name: name, Arguments: argsMap}}
		for _, c := range clients {
			res, err := c.CallTool(ctx, req)
			if err != nil {
				continue
			}
			return formatCallToolResult(res), nil
		}
	}

	return fmt.Sprintf("Tool %s not found or execution failed", name), fmt.Errorf("tool not found")
}

func formatCallToolResult(res *mcp.CallToolResult) string {
	if res == nil {
		return ""
	}
	var buf string
	for _, c := range res.Content {
		buf += mcp.GetTextFromContent(c)
	}
	return buf
}
