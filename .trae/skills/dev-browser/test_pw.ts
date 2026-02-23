import { chromium } from 'playwright';
(async () => {
  try {
    const browser = await chromium.launch();
    console.log('Browser launched successfully');
    await browser.close();
  } catch (e) {
    console.error('Failed to launch browser:', e);
  }
})();