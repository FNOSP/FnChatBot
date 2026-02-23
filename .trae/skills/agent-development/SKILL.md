---
name: "agent-development"
description: "Use when creating agents, writing subagents, or needing guidance on agent structure, frontmatter, system prompts, triggers, and best practices for Claude Code plugins."
---

# Agent Development for Claude Code Plugins

## Overview

Agents are autonomous subprocesses that handle complex, multi-step tasks independently. Understanding agent structure, triggering conditions, and system prompt design enables creating powerful autonomous capabilities.

**Key concepts:**
- Agents are FOR autonomous work, commands are FOR user-initiated actions
- Markdown file format with YAML frontmatter
- Triggering via description field with examples
- System prompt defines agent behavior
- Model and color customization

## Agent File Structure

### Complete Format

```markdown
---
name: agent-identifier
description: Use this agent when [triggering conditions]. Examples:

<example>
Context: [Situation description]
user: "[User request]"
assistant: "[How assistant should respond and use this agent]"
<commentary>
[Why this agent should be triggered]
</commentary>
</example>

<example>
[Additional example...]
</example>

model: inherit
color: blue
tools: ["Read", "Write", "Grep"]
---

You are [agent role description]...

**Your Core Responsibilities:**
1. [Responsibility 1]
2. [Responsibility 2]

**Analysis Process:**
[Step-by-step workflow]

**Output Format:**
[What to return]
```

## Frontmatter Fields

### name (required)

Agent identifier used for namespacing and invocation.

**Format:** lowercase, numbers, hyphens only
**Length:** 3-50 characters
**Pattern:** Must start and end with alphanumeric

**Good examples:**
- `code-reviewer`
- `test-generator`
- `api-docs-writer`
- `security-analyzer`

**Bad examples:**
- `helper` (too generic)
- `-agent-` (starts/ends with hyphen)
- `my_agent` (underscores not allowed)
- `ag` (too short, < 3 chars)

### description (required)

Defines when Claude should trigger this agent. **This is the most critical field.**

**Must include:**
1. Triggering conditions ("Use this agent when...")
2. Multiple `<example>` blocks showing usage
3. Context, user request, and assistant response in each example
4. `<commentary>` explaining why agent triggers

**Format:**
```
Use this agent when [conditions]. Examples:

<example>
Context: [Scenario description]
user: "[What user says]"
assistant: "[How Claude should respond]"
<commentary>
[Why this agent is appropriate]
</commentary>
</example>

[More examples...]
```

**Best practices:**
- Include 2-4 concrete examples
- Show proactive and reactive triggering
- Cover different phrasings of same intent
- Explain reasoning in commentary
- Be specific about when NOT to use the agent

### model (required)

Which model the agent should use.

**Options:**
- `inherit` - Use same model as parent (recommended)
- `sonnet` - Claude Sonnet (balanced)
- `opus` - Claude Opus (most capable, expensive)
- `haiku` - Claude Haiku (fast, cheap)

**Recommendation:** Use `inherit` unless agent needs specific model capabilities.

### color (required)

Visual identifier for agent in UI.

**Options:** `blue`, `cyan`, `green`, `yellow`, `magenta`, `red`

**Guidelines:**
- Choose distinct colors for different agents in same plugin
- Use consistent colors for similar agent types
- Blue/cyan: Analysis, review
- Green: Success-oriented tasks
- Yellow: Caution, validation
- Red: Critical, security
- Magenta: Creative, generation

### tools (optional)

Restrict agent to specific tools.

**Format:** Array of tool names

```yaml
tools: ["Read", "Write", "Grep", "Bash"]
```

**Default:** If omitted, agent has access to all tools

**Best practice:** Limit tools to minimum needed (principle of least privilege)

**Common tool sets:**
- Read-only analysis: `["Read", "Grep", "Glob"]`
- Code generation: `["Read", "Write", "Grep"]`
- Testing: `["Read", "Bash", "Grep"]`
- Full access: Omit field or use `["*"]`

## System Prompt Design

The markdown body becomes the agent's system prompt. Write in second person, addressing the agent directly.

### Structure

**Standard template:**
```markdown
You are [role] specializing in [domain].

**Your Core Responsibilities:**
1. [Primary responsibility]
2. [Secondary responsibility]
3. [Additional responsibilities...]

**Analysis Process:**
1. [Step one]
2. [Step two]
3. [Step three]
[...]

**Quality Standards:**
- [Standard 1]
- [Standard 2]

**Output Format:**
Provide results in this format:
- [What to include]
- [How to structure]

**Edge Cases:**
Handle these situations:
- [Edge case 1]: [How to handle]
- [Edge case 2]: [How to handle]
```

### Best Practices

✅ **DO:**
- Write in second person ("You are...", "You will...")
- Be specific about responsibilities
- Provide step-by-step process
- Define output format
- Include quality standards
- Address edge cases
- Keep under 10,000 characters

❌ **DON'T:**
- Write in first person ("I am...")
- Be vague ("Do good work")
- Include internal thoughts (unless requested)
- Leave output format undefined
