import { connect, waitForPageLoad } from "@/client.js";

async function main() {
  const client = await connect();
  const page = await client.page("test-chat");

  try {
    // 1. Go to Settings
    console.log("Navigating to Settings...");
    await page.goto("http://localhost:5173/settings");
    await waitForPageLoad(page);

    // 2. Select Model Services
    console.log("Selecting Model Services...");
    // Find the Model Services menu item
    await page.click('text="Model Services"'); // Adjust selector if needed
    
    // 3. Configure API
    console.log("Configuring API...");
    // Find Base URL input
    // Using placeholder or label might be hard if Semi UI hides input
    // Let's try to fill inputs by order or specific selector if possible
    // Semi UI inputs often have a wrapper.
    // We can use page.fill with placeholder if available.
    
    // Select "OpenAI" provider if not selected (default is usually OpenAI)
    // Assuming OpenAI is selected.
    
    // Fill Base URL
    // Placeholder "https://api..."
    await page.fill('input[placeholder*="https://api"]', "");
    
    // Fill API Key
    // Placeholder "sk-..." or similar
    const apiKey = process.env.OPENAI_API_KEY ?? "";
    if (!apiKey) {
        throw new Error("Missing OPENAI_API_KEY");
    }
    await page.fill('input[type="password"]', apiKey);
    
    // 4. Fetch Models
    console.log("Fetching models...");
    await page.click('text="Fetch Available Models"'); // Button text from en.json
    
    // Wait for modal or success toast
    await page.waitForTimeout(3000);
    
    // 5. Select Model in Modal
    // Check if modal is visible
    // We need to select "gpt-4o-mini-ca"
    // Assuming it's in the list.
    // Click "Select All" or find specific checkbox
    // Let's try to click the checkbox for gpt-4o-mini-ca
    // The checkbox value is the model ID.
    // Semi UI Checkbox structure: label > input[type=checkbox]
    // Or text.
    await page.click('text="gpt-4o-mini-ca"');
    
    // Click "Add Selected"
    await page.click('text="Add Selected Models"');
    await page.waitForTimeout(1000);

    // 6. Go to Chat
    console.log("Going to Chat...");
    await page.goto("http://localhost:5173/");
    await waitForPageLoad(page);
    
    // 7. Create New Chat
    // "New Chat" button in sidebar
    await page.click('text="New Chat"');
    await page.waitForTimeout(1000);
    
    // Select Model?
    // Usually defaults to first or we need to select it.
    // Assuming default is set or we can select.
    
    // 8. Send Message
    console.log("Sending message...");
    await page.fill('textarea', "Hello, are you real?");
    await page.keyboard.press("Enter");
    
    // 9. Wait for Response
    await page.waitForTimeout(5000);
    
    // 10. Verify Response
    // Get last message content
    // We can take a screenshot or extract text
    const content = await page.content();
    console.log("Page content length:", content.length);
    
    // Check for "simulated response" text to ensure we are NOT mocked
    if (content.includes("simulated response")) {
        console.error("FAILURE: Received mocked response!");
    } else {
        console.log("SUCCESS: Did not detect mocked response text.");
    }
    
    await page.screenshot({ path: "tmp/chat_result.png" });
    
  } catch (e) {
    console.error("Error:", e);
    await page.screenshot({ path: "tmp/error.png" });
  } finally {
    await client.disconnect();
  }
}

main();
