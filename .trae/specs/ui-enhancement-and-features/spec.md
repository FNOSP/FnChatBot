# UI Enhancement and Features Spec

## Why
The current frontend UI is basic and lacks the advanced configuration capabilities required for a robust AI Agent system. The user specifically requested a model management interface inspired by Cherry Studio and the implementation of MCP (Model Context Protocol) and Skill management features (referencing Trae's skill import dialog) as originally planned in the Architecture document.

## What Changes

### Frontend
- **Refactor Settings Layout**:
  - Adopt a sidebar-based layout for the Settings view.
  - Categories: Model Services, MCP Servers, Skill Management, etc.
- **Model Services UI (Cherry Studio Style)**:
  - **Sidebar**: List of providers (OpenAI, Anthropic, Custom, etc.) with toggle switches.
  - **Content Area**:
    - API Key input (masked with toggle).
    - API Address (Base URL) input.
    - Model List management (Add/Remove models for the provider).
    - "Check Connection" button.
- **MCP Management UI**:
  - List of configured MCP servers with status indicators.
  - "Add MCP Server" dialog (Name, Base URL, API Key).
  - Edit/Delete functionality.
- **Skill Management UI**:
  - List of enabled/disabled skills.
  - "Add Skill" dialog (Trae Style):
    - File upload area (supports `.md` and `.zip`).
    - Automatic parsing of skill name, description, and instructions.
    - Form to review/edit parsed details before saving.

### Backend
- **API Enhancements**:
  - **Model Config**: Ensure API supports retrieving and saving configurations grouped by provider.
  - **MCP Config**: Implement full CRUD endpoints for `mcp_configs` table.
  - **Skill Management**:
    - Implement file upload endpoint (`POST /api/skills/upload`).
    - Logic to parse `.md` files (extracting YAML frontmatter or sections).
    - Logic to parse `.zip` files (finding and parsing `SKILL.md`).
    - CRUD endpoints for `skills` table.

## Impact
- **Affected Specs**: `PRD.md`, `Architecture.md` (implementing pending features).
- **Affected Code**:
  - Frontend: `src/views/SettingsView.vue`, `src/components/settings/*`
  - Backend: `internal/api/handlers.go`, `internal/models/skill.go`, `internal/services/skill_parser.go` (new)

## ADDED Requirements
### Requirement: Enhanced Model Configuration
The system SHALL provide a model configuration interface with provider-specific settings, API key management, and model list customization, visually similar to Cherry Studio.

### Requirement: MCP Management
The system SHALL allow users to manually add, edit, and delete MCP server configurations via the UI.

### Requirement: Skill Import
The system SHALL allow users to import skills by uploading `.md` or `.zip` files. The system MUST automatically parse metadata (name, description) and instructions from the uploaded files.
