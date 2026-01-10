# CannaNote Color Palette

## Color System Overview

CannaNote's color palette reflects our brand values of trustworthiness, calm professionalism, and subtle natural connection while maintaining strict app store compliance. The palette avoids cannabis clichés and instead focuses on sophisticated, desaturated tones that convey evidence-based wellness.

---

## Primary Colors
*Core brand colors for headers, primary actions, and brand presence*

### Deep Navy
- **Hex:** `#1a1f36`
- **Usage:** Primary brand color, headers, high-trust elements
- **Personality:** Professional, trustworthy, calming authority

### Midnight
- **Hex:** `#0f1419` 
- **Usage:** Dark backgrounds, depth, premium feel
- **Personality:** Sophisticated, private, focused

### Charcoal
- **Hex:** `#2d3748`
- **Usage:** Secondary text, borders, subtle structure
- **Personality:** Reliable, readable, unobtrusive

---

## Natural Colors
*Subtle earth tones that reference wellness without cannabis clichés*

### Soft Sage
- **Hex:** `#7a9b8e`
- **Usage:** Subtle natural accent, wellness indicators
- **Personality:** Calming, organic, mindful

### Muted Forest
- **Hex:** `#4a6358`
- **Usage:** Secondary natural accent, environmental themes
- **Personality:** Grounded, sustainable, protective

### Warm Earth
- **Hex:** `#8b7d6b`
- **Usage:** Tertiary organic tone, comfort elements
- **Personality:** Stable, nurturing, authentic

---

## Neutral Colors
*Clean foundation colors for readability and accessibility*

### Clean White
- **Hex:** `#ffffff`
- **Usage:** Light backgrounds, cards, clean space
- **Personality:** Pure, clear, transparent

### Off White
- **Hex:** `#f8f9fa`
- **Usage:** Subtle backgrounds, reduced glare
- **Personality:** Gentle, comfortable, approachable

### Light Gray
- **Hex:** `#e8eaed`
- **Usage:** Borders, dividers, subtle structure
- **Personality:** Organized, clean, unobtrusive

---

## Accent Colors
*Functional colors for UI states and data visualization*

### Insight Blue
- **Hex:** `#5b8db8`
- **Usage:** Data insights, links, informational elements
- **Personality:** Intelligent, trustworthy, clear

### Success Green
- **Hex:** `#6b9080`
- **Usage:** Confirmations, progress indicators, positive states
- **Personality:** Accomplishing, healthy, encouraging

### Warning Amber
- **Hex:** `#c9a86a`
- **Usage:** Alerts, cautions, attention needed
- **Personality:** Thoughtful, measured, protective

### Error Crimson
- **Hex:** `#a85857`
- **Usage:** Errors, danger states, critical alerts
- **Personality:** Serious, protective, clear boundaries

---

## Usage Guidelines

### Accessibility
- All color combinations meet WCAG 2.1 AA contrast requirements
- Text on colored backgrounds uses sufficient contrast ratios
- Color is never the only way to convey information

### Brand Consistency
- **Primary Colors:** Use Deep Navy for most brand applications
- **Natural Colors:** Use sparingly as accents, never as primary brand colors
- **Neutral Colors:** Foundation for most interface elements
- **Accent Colors:** Functional use only, never decorative

### App Store Compliance
- **No bright "cannabis greens"** that could trigger platform restrictions
- **Professional palette** suitable for health and wellness apps
- **Sophisticated tones** that convey medical legitimacy

### Dark Mode Support
- Primary colors work well against dark backgrounds
- Natural and accent colors maintain visibility in dark themes
- All colors tested for dark mode accessibility

---

## Implementation Notes

### CSS Variables
```css
:root {
  /* Primary */
  --color-deep-navy: #1a1f36;
  --color-midnight: #0f1419;
  --color-charcoal: #2d3748;
  
  /* Natural */
  --color-soft-sage: #7a9b8e;
  --color-muted-forest: #4a6358;
  --color-warm-earth: #8b7d6b;
  
  /* Neutral */
  --color-clean-white: #ffffff;
  --color-off-white: #f8f9fa;
  --color-light-gray: #e8eaed;
  
  /* Accent */
  --color-insight-blue: #5b8db8;
  --color-success-green: #6b9080;
  --color-warning-amber: #c9a86a;
  --color-error-crimson: #a85857;
}
```

### Design System Integration
- Colors align with minimalist, tech-forward aesthetic
- Palette supports both light and dark interface themes
- Flexible enough for data visualization while maintaining brand consistency

---

*This color palette serves as the foundation for all CannaNote visual communications, ensuring brand consistency while supporting user trust and app store compliance.*