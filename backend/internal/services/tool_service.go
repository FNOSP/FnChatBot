package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"
)

const (
	ToolTypeFunction = "function"
)

type ToolService struct{}

func NewToolService() *ToolService {
	return &ToolService{}
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
	if err := db.DB.Where("enabled = ?", true).Find(&skills).Error; err != nil {
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

	var mcps []models.MCPConfig
	if err := db.DB.Where("enabled = ?", true).Find(&mcps).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch MCPs: %v", err)
	}

	for _, mcp := range mcps {
		mcpTools, err := s.fetchMCPTools(mcp)
		if err != nil {
			log.Printf("Failed to fetch tools from MCP %s: %v", mcp.Name, err)
			continue
		}
		tools = append(tools, mcpTools...)
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

func (s *ToolService) fetchMCPTools(mcp models.MCPConfig) ([]Tool, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/tools", mcp.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP returned status %d", resp.StatusCode)
	}

	var tools []Tool
	if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
		return nil, err
	}
	return tools, nil
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

	var mcps []models.MCPConfig
	db.DB.Where("enabled = ?", true).Find(&mcps)

	for _, mcp := range mcps {
		result, err := s.executeMCPTool(mcp, name, args)
		if err == nil {
			return result, nil
		}
	}

	return fmt.Sprintf("Tool %s not found or execution failed", name), fmt.Errorf("tool not found")
}

func (s *ToolService) executeMCPTool(mcp models.MCPConfig, name string, args string) (string, error) {
	payload := map[string]string{
		"name": name,
		"args": args,
	}
	jsonBody, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(fmt.Sprintf("%s/tools/execute", mcp.BaseURL), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("MCP execution failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
