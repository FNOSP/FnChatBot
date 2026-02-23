import { connect, waitForPageLoad } from "@/client.js";

// Connect to the browser
const client = await connect();
const page = await client.page("fnchatbot-e2e");

console.log("Navigating to app...");
await page.goto("http://localhost:5173");
await waitForPageLoad(page);

// 1. Test Complex Workflow (Task Planning)
console.log("Testing Task Planning...");
// Wait for input
await page.waitForSelector("textarea");
await page.fill("textarea", "Create a plan to refactor the project structure. Please use TodoWrite to list tasks.");
await page.click("button:has(svg.lucide-send)");

// Wait for response and TaskPanel
// TaskPanel appears when tasks > 0. It has class "w-80 border-l ..."
// We can check for text "Plan & Progress"
console.log("Waiting for Task Panel...");
try {
  await page.waitForSelector("text=Plan & Progress", { timeout: 30000 });
  console.log("✅ Task Panel visible");
} catch (e) {
  console.error("❌ Task Panel not found");
}

// 2. Test Subagent
console.log("Testing Subagent...");
await page.fill("textarea", "Please delegate a task to explore the codebase using the Task tool.");
await page.click("button:has(svg.lucide-send)");

console.log("Waiting for Subagent execution...");
try {
  // Check for "Subagent:" text in message
  await page.waitForSelector("text=Subagent:", { timeout: 30000 });
  console.log("✅ Subagent execution visible");
} catch (e) {
  console.error("❌ Subagent execution not found");
}

// 3. Test Skill Loading
console.log("Testing Skill Loading...");
await page.fill("textarea", "Load the 'code-review' skill.");
await page.click("button:has(svg.lucide-send)");

console.log("Waiting for Skill Loading...");
try {
  // Check for "Skill Loaded" text
  await page.waitForSelector("text=Skill Loaded", { timeout: 30000 });
  console.log("✅ Skill Loading visible");
} catch (e) {
  console.error("❌ Skill Loading not found");
}

// 4. Test Sandbox Mode
console.log("Testing Sandbox Mode...");
// We ask for a bash command. Since we don't have a real LLM that might output bash reliably without prompt,
// we can try to force it or assume the backend mock returns a bash block if we ask for it?
// Actually, our tool service returns text. The LLM generates the bash block.
// We are using `gpt-4o-mini-ca` which is a real model (via proxy).
await page.fill("textarea", "List files in current directory using bash.");
await page.click("button:has(svg.lucide-send)");

console.log("Waiting for Sandbox Tag...");
try {
  // Check for "Sandbox Mode" text
  await page.waitForSelector("text=Sandbox Mode", { timeout: 30000 });
  console.log("✅ Sandbox Mode tag visible");
} catch (e) {
  console.error("❌ Sandbox Mode tag not found");
}

console.log("Test Complete.");
// Keep page open for manual inspection if needed, or close
// await client.disconnect();
