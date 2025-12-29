---
intent: 001-chat-log-viewer
type: reference
created: 2025-12-29T10:45:00Z
updated: 2025-12-29T11:00:00Z
---

# UI Standards: Claude.ai Visual Style

## Overview

The generated HTML must closely match the Claude.ai web chat interface styling. This document contains the exact CSS variables and patterns extracted from the Claude.ai interface.

## Color Palette (HSL Values)

### Light Mode (`data-theme="claude" data-mode="light"`)

```css
:root {
  /* Always colors */
  --always-white: 0 0% 100%;
  --always-black: 0 0% 0%;

  /* Accent - Main (Claude Orange) */
  --accent-main-000: 15 54.2% 51.2%;
  --accent-main-100: 15 54.2% 51.2%;
  --accent-main-200: 15 63.1% 59.6%;
  --accent-main-900: 0 0% 0%;
  --accent-brand: 15 63.1% 59.6%;

  /* Accent - Secondary (Blue) */
  --accent-secondary-000: 210 73.7% 40.2%;
  --accent-secondary-100: 210 70.9% 51.6%;
  --accent-secondary-200: 210 70.9% 51.6%;
  --accent-secondary-900: 211 72% 90%;

  /* Accent - Pro (Purple) */
  --accent-pro-000: 251 34.2% 33.3%;
  --accent-pro-100: 251 40% 45.1%;
  --accent-pro-200: 251 61% 72.2%;
  --accent-pro-900: 253 33.3% 91.8%;

  /* Backgrounds (Warm Cream Tones) */
  --bg-000: 0 0% 100%;           /* Pure white */
  --bg-100: 48 33.3% 97.1%;      /* Cream - main background */
  --bg-200: 53 28.6% 94.5%;      /* Slightly darker cream */
  --bg-300: 48 25% 92.2%;        /* Hover/active states */
  --bg-400: 50 20.7% 88.6%;      /* Pressed states */
  --bg-500: 50 20.7% 88.6%;

  /* Text */
  --text-000: 60 2.6% 7.6%;      /* Near black */
  --text-100: 60 2.6% 7.6%;      /* Primary text */
  --text-200: 60 2.5% 23.3%;     /* Secondary text */
  --text-300: 60 2.5% 23.3%;
  --text-400: 51 3.1% 43.7%;     /* Muted text */
  --text-500: 51 3.1% 43.7%;     /* Subtle text */

  /* Borders */
  --border-100: 30 3.3% 11.8%;
  --border-200: 30 3.3% 11.8%;
  --border-300: 30 3.3% 11.8%;
  --border-400: 30 3.3% 11.8%;

  /* Status Colors */
  --danger-000: 0 58.6% 34.1%;
  --danger-100: 0 56.2% 45.4%;
  --danger-200: 0 56.2% 45.4%;
  --danger-900: 0 50% 95%;

  --success-000: 125 100% 18%;
  --success-100: 103 72.3% 26.9%;
  --success-200: 103 72.3% 26.9%;
  --success-900: 86 45.1% 90%;

  --warning-000: 45 91.8% 19%;
  --warning-100: 39 88.8% 28%;
  --warning-200: 39 88.8% 28%;
  --warning-900: 38 65.9% 92%;

  /* On-color (text on colored backgrounds) */
  --oncolor-100: 0 0% 100%;
  --oncolor-200: 60 6.7% 97.1%;
  --oncolor-300: 60 6.7% 97.1%;
}
```

### Dark Mode (`data-theme="claude" data-mode="dark"`)

```css
[data-theme="claude"][data-mode="dark"] {
  /* Backgrounds (Dark warm tones) */
  --bg-000: 60 2.1% 18.4%;
  --bg-100: 60 2.7% 14.5%;
  --bg-200: 30 3.3% 11.8%;
  --bg-300: 60 2.6% 7.6%;
  --bg-400: 0 0% 0%;
  --bg-500: 0 0% 0%;

  /* Text (Light cream tones) */
  --text-000: 48 33.3% 97.1%;
  --text-100: 48 33.3% 97.1%;
  --text-200: 50 9% 73.7%;
  --text-300: 50 9% 73.7%;
  --text-400: 48 4.8% 59.2%;
  --text-500: 48 4.8% 59.2%;

  /* Borders (Light for visibility) */
  --border-100: 51 16.5% 84.5%;
  --border-200: 51 16.5% 84.5%;
  --border-300: 51 16.5% 84.5%;
  --border-400: 51 16.5% 84.5%;
}
```

## Typography

### Font Families

```css
:root {
  --font-mono: var(--font-jetbrains);
  --font-ui: var(--font-anthropic-sans);
  --font-ui-serif: var(--font-anthropic-serif);
  --font-claude-response: var(--font-anthropic-serif);
  --font-user-message: var(--font-ui);
  --font-sans-serif: var(--font-ui);
  --font-serif: var(--font-ui-serif);
  --font-system: system-ui, sans-serif;
}
```

**Fallback fonts for our implementation:**
```css
:root {
  --font-ui: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  --font-serif: Georgia, "Times New Roman", serif;
  --font-mono: "JetBrains Mono", "Fira Code", Consolas, monospace;
}
```

### Font Weights

```css
.font-normal { font-weight: 430; }
.font-medium { font-weight: 550; }
.font-semibold { font-weight: 580; }
.font-bold { font-weight: 600; }
```

### Font Classes

```css
.font-base { font-size: 14px; }
.font-small { font-size: 12px; }
.font-large { font-size: 16px; }
.font-xl { font-size: 18px; }
```

## Message Styling

### User Messages

```css
.user-message {
  font-family: var(--font-user-message);  /* Sans-serif */
  background: hsl(var(--bg-300));
  border-radius: 0.75rem;  /* rounded-xl */
  padding: 0.625rem 1rem;  /* py-2.5 px-4 */
  color: hsl(var(--text-100));
  max-width: 75ch;
}
```

### Claude/Assistant Messages

```css
.claude-response {
  font-family: var(--font-claude-response);  /* Serif */
  line-height: 1.7;
  font-weight: 360;  /* Slightly lighter in dark mode */
}

.font-claude-response-body {
  word-wrap: break-word;
  white-space: normal;
  line-height: 1.7;
}
```

## Layout Patterns

### Sidebar Navigation

```css
nav.sidebar {
  width: 18rem;  /* 288px */
  position: fixed;
  left: 0;
  height: 100vh;
  border-right: 0.5px solid hsl(var(--border-300));
  background: linear-gradient(to top, hsl(var(--bg-200) / 5%), hsl(var(--bg-200) / 30%));
}
```

### Sidebar Items

```css
.sidebar-item {
  display: inline-flex;
  align-items: center;
  height: 2rem;  /* h-8 */
  padding: 0.375rem 1rem;  /* py-1.5 px-4 */
  border-radius: 0.5rem;  /* rounded-lg */
  font-size: 0.875rem;  /* text-sm */
  width: 100%;
  transition: 75ms;
}

.sidebar-item:hover {
  background: hsl(var(--bg-300));
}

.sidebar-item:active {
  background: hsl(var(--bg-300));
  transform: scale(1.0);
}

.sidebar-item.active {
  background: hsl(var(--bg-300));
}
```

### Content Area

```css
.content-area {
  max-width: 48rem;  /* max-w-3xl */
  margin: 0 auto;
  padding: 0 0.5rem;
}
```

## Component Patterns

### Buttons

```css
/* Ghost Button (most common) */
.button-ghost {
  background: transparent;
  color: hsl(var(--text-300));
  border: transparent;
  transition: 300ms cubic-bezier(0.165, 0.85, 0.45, 1);
}

.button-ghost:hover {
  background: hsl(var(--bg-300));
  color: hsl(var(--text-100));
}

.button-ghost:active {
  transform: scale(0.95);
}

/* Primary Button */
.button-primary {
  background: hsl(var(--text-000));
  color: hsl(var(--bg-000));
}

/* Claude Orange Button */
.button-claude {
  background: hsl(var(--accent-main-000));
  color: hsl(var(--oncolor-100));
}
```

### Inline Code

```css
code {
  background: hsl(var(--text-200) / 5%);
  border: 0.5px solid hsl(var(--border-300));
  color: hsl(var(--danger-000));  /* Red-ish for code */
  white-space: pre-wrap;
  border-radius: 0.4rem;
  padding: 1px 4px;
  font-size: 0.9rem;
}
```

### Code Blocks (Syntax Highlighting)

```css
pre code.hljs {
  display: block;
  overflow-x: auto;
  padding: 1em;
  color: #abb2bf;
  background: #282c34;
}

.hljs-keyword { color: #c678dd; }
.hljs-string { color: #98c379; }
.hljs-number { color: #d19a66; }
.hljs-comment { color: #5c6370; font-style: italic; }
.hljs-title { color: #61aeee; }
.hljs-built_in { color: #e6c07b; }
```

## Spacing System

```css
/* Tailwind-compatible spacing */
--space-0: 0;
--space-px: 1px;
--space-0.5: 0.125rem;  /* 2px */
--space-1: 0.25rem;     /* 4px */
--space-1.5: 0.375rem;  /* 6px */
--space-2: 0.5rem;      /* 8px */
--space-2.5: 0.625rem;  /* 10px */
--space-3: 0.75rem;     /* 12px */
--space-4: 1rem;        /* 16px */
--space-5: 1.25rem;     /* 20px */
--space-6: 1.5rem;      /* 24px */
--space-8: 2rem;        /* 32px */
```

## Animation & Transitions

```css
/* Standard transition */
.transition-standard {
  transition: 300ms cubic-bezier(0.165, 0.85, 0.45, 1);
}

/* Quick transition */
.transition-quick {
  transition: 75ms ease;
}

/* Scale on active */
.active-scale:active {
  transform: scale(0.95);
}
```

## Z-Index Scale

```css
--z-header: 20;
--z-sidebar: 30;
--z-modal: 40;
--z-dropdown: 50;
--z-overlay: 50;
--z-tooltip: 50;
--z-toast: 60;
```

## Footer Branding

```html
<footer class="footer-branding">
  claude-code-logs by <a href="https://fabriqa.ai" target="_blank">fabriqa.ai</a>
</footer>
```

```css
.footer-branding {
  text-align: center;
  padding: 1rem;
  font-size: 0.75rem;
  color: hsl(var(--text-500));
  border-top: 0.5px solid hsl(var(--border-300));
}

.footer-branding a {
  color: hsl(var(--accent-main-100));
  text-decoration: none;
}

.footer-branding a:hover {
  text-decoration: underline;
}
```

## Key Visual Characteristics

1. **Warm Cream Backgrounds** - Not pure white, uses HSL with warm undertones
2. **Serif for Claude** - Assistant responses use serif font (Georgia fallback)
3. **Sans-serif for User** - User messages use system sans-serif
4. **Subtle Borders** - 0.5px borders, not harsh lines
5. **Orange Accent** - Claude brand color `hsl(15, 54.2%, 51.2%)`
6. **Smooth Transitions** - Cubic bezier easing for polish
7. **Light Touch Shadows** - Very subtle, almost imperceptible
