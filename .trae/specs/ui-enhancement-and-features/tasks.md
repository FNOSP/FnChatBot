# Tasks

* [x] Task 1: Backend - Implement MCP and Skill APIs

  * [x] SubTask 1.1: Implement CRUD endpoints for MCP configurations (`/api/mcp`).

  * [x] SubTask 1.2: Create `SkillParser` service to handle `.md` and `.zip` file parsing.

  * [x] SubTask 1.3: Implement Skill upload endpoint (`/api/skills/upload`) and CRUD endpoints.

  * [x] SubTask 1.4: Verify Model Config API supports necessary fields for the new UI.

* [x] Task 2: Frontend - Refactor Settings Layout

  * [x] SubTask 2.1: Create a sidebar layout for `SettingsView.vue`.

  * [x] SubTask 2.2: Define routes/tabs for Models, MCP, and Skills.

* [x] Task 3: Frontend - Implement Model Services UI

  * [x] SubTask 3.1: Create provider sidebar (OpenAI, Anthropic, Custom, etc.).

  * [x] SubTask 3.2: Implement configuration form (API Key, Base URL).

  * [x] SubTask 3.3: Implement Model List management (add/remove models).

* [x] Task 4: Frontend - Implement MCP Management UI

  * [x] SubTask 4.1: Create MCP Server list view.

  * [x] SubTask 4.2: Implement Add/Edit MCP Server dialog.

* [x] Task 5: Frontend - Implement Skill Management UI

  * [x] SubTask 5.1: Create Skill list view with enable/disable toggles.

  * [x] SubTask 5.2: Implement "Add Skill" dialog with file upload (Trae style).

  * [x] SubTask 5.3: Integrate with backend upload API and display parsed results.

# Task Dependencies

* Task 3, 4, 5 depend on Task 2 (Layout).

* Task 5 depends on Task 1 (Backend API).

