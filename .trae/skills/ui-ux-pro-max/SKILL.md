---
name: "ui-ux-pro-max"
description: "UI/UX design intelligence. Invoke when user asks for UI design, improvement, review, or frontend implementation. Provides styles, palettes, fonts, and UX guidelines."
---

# UI/UX Pro Max - Design Intelligence

Comprehensive design guide for web and mobile applications. Contains design intelligence for styles, color palettes, typography, UX guidelines, and chart types.

## When to Apply
Reference these guidelines when:
- Designing new UI components or pages
- Choosing color palettes and typography
- Reviewing code for UX issues
- Building landing pages or dashboards
- Implementing accessibility requirements

## Rule Categories by Priority

1. **Accessibility (CRITICAL)**
   - `color-contrast`: Minimum 4.5:1 ratio for normal text
   - `focus-states`: Visible focus rings on interactive elements
   - `alt-text`: Descriptive alt text for meaningful images
   - `aria-labels`: `aria-label` for icon-only buttons
   - `keyboard-nav`: Tab order matches visual order
   - `form-labels`: Use `<label>` with `for` attribute

2. **Touch & Interaction (CRITICAL)**
   - `touch-target-size`: Minimum 44x44px touch targets
   - `hover-vs-tap`: Use click/tap for primary interactions (no hover-only critical actions)
   - `loading-buttons`: Disable button during async operations
   - `error-feedback`: Clear error messages near problem
   - `cursor-pointer`: Add `cursor-pointer` to clickable elements

3. **Performance (HIGH)**
   - `image-optimization`: Use WebP, srcset, lazy loading
   - `reduced-motion`: Check `prefers-reduced-motion`
   - `content-jumping`: Reserve space for async content (CLS)

4. **Layout & Responsive (HIGH)**
   - `viewport-meta`: `width=device-width initial-scale=1`
   - `readable-font-size`: Minimum 16px body text on mobile
   - `horizontal-scroll`: Ensure content fits viewport width (no accidental scroll)
   - `z-index-management`: Define z-index scale (e.g., 10, 20, 30, 50)

5. **Typography & Color (MEDIUM)**
   - `line-height`: Use 1.5-1.75 for body text
   - `line-length`: Limit to 65-75 characters per line
   - `font-pairing`: Match heading/body font personalities
   - `consistency`: Use consistent H1-H6 hierarchy

6. **Animation (MEDIUM)**
   - `duration-timing`: Use 150-300ms for micro-interactions
   - `transform-performance`: Use transform/opacity, not width/height/top/left
   - `loading-states`: Skeleton screens or spinners

7. **Style Selection (MEDIUM)**
   - `style-match`: Match style to product type (e.g., Trust for Fintech, Playful for Gaming)
   - `consistency`: Use same style across all pages
   - `no-emoji-icons`: Use SVG icons (Lucide, Heroicons), not emojis

8. **Charts & Data (LOW)**
   - `chart-type`: Match chart type to data type
   - `color-guidance`: Use accessible color palettes for data
   - `data-table`: Provide table alternative for accessibility

## Workflow for AI Assistant

When user requests UI/UX work (design, build, create, implement, review, fix, improve), follow this workflow:

### Step 1: Analyze User Requirements
Extract key information:
- **Product type**: SaaS, e-commerce, portfolio, dashboard, landing page, etc.
- **Style keywords**: minimal, playful, professional, elegant, dark mode, etc.
- **Industry**: healthcare, fintech, gaming, education, etc.
- **Stack**: React, Vue, Next.js, or default to html-tailwind

### Step 2: Generate Design System (Mental Draft)
Before coding, formulate a design system:
- **Pattern**: e.g., Hero-Centric, Dashboard, Feed
- **Colors**: Primary, Secondary, CTA, Background, Text (ensure contrast)
- **Typography**: Heading font, Body font (Google Fonts)
- **Effects**: Shadows, border-radius, transitions

### Step 3: Implementation Guidelines
- **HTML/Tailwind**: Use utility classes for layout, spacing, colors.
- **React/Vue/Svelte**: Componentize reusable elements (Buttons, Cards).
- **Icons**: Use SVG icons (e.g., `lucide-react` or similar).
- **Responsive**: Mobile-first approach (`class`, `md:class`, `lg:class`).

### Example Design System Output
When asked to design, produce something like:
> **Recommended Design System**
> - **Style**: Soft UI / Minimalism
> - **Colors**: Navy Blue (#1e293b) + Emerald Green (#10b981) for accents
> - **Font**: Inter (Body) + Playfair Display (Headings)
> - **Shadows**: `shadow-md` for cards, `shadow-lg` for hover
> - **Radius**: `rounded-lg` (0.5rem)
