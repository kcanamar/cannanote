# CannaNote Typography System

## Typography Philosophy

CannaNote's typography reflects our brand personality: wise, calm, trustworthy, and approachable. The system combines modern geometric sans-serif with subtle technical monospace touches, creating a sophisticated yet grounded aesthetic that conveys both data precision and human warmth.

---

## Primary Font: Space Grotesk

### Font Family
**Space Grotesk** - A modern, geometric sans-serif font that balances professionalism with approachability.

### Characteristics
- **Style:** Geometric sans-serif with subtle humanist touches
- **Personality:** Modern, trustworthy, clean, approachable
- **Usage:** Primary brand applications, UI text, marketing materials

### Weight Variations
- **400 (Regular):** Body text, secondary content
- **500 (Medium):** Emphasized text, navigation
- **600 (Semi-Bold):** Subheadings, important information
- **700 (Bold):** Headlines, primary emphasis

### Google Fonts Import
```css
@import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&display=swap');
```

---

## Secondary Font: JetBrains Mono

### Font Family
**JetBrains Mono** - A technical monospace font designed for code and data display.

### Characteristics
- **Style:** Monospace with excellent readability
- **Personality:** Technical, precise, trustworthy, sophisticated
- **Usage:** Data display, code snippets, technical contexts, subtle tech accent

### Weight Variations
- **400 (Regular):** Data displays, technical information
- **500 (Medium):** Emphasized technical content, labels

### Google Fonts Import
```css
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap');
```

---

## CannaNote Wordmark Treatment

### Standard Wordmark
- **Font:** Space Grotesk
- **Structure:** "Canna" (Regular 400) + "Note" (Medium 500)
- **Styling:** Single line, tight tracking, camel case emphasis

### Implementation
```css
.cannanote-wordmark {
  font-family: 'Space Grotesk', sans-serif;
  letter-spacing: -0.02em;
}

.cannanote-wordmark .canna {
  font-weight: 400; /* Regular */
}

.cannanote-wordmark .note {
  font-weight: 500; /* Medium */
}
```

### Size Variations
- **Small (sm):** `text-xl` - Navigation, compact spaces
- **Medium (md):** `text-3xl` - Standard applications
- **Large (lg):** `text-5xl` - Headers, marketing
- **Extra Large (xl):** `text-7xl` - Hero sections, print

---

## Typography Hierarchy

### Headlines
- **Font:** Space Grotesk
- **Weight:** 700 (Bold) for H1, 600 (Semi-Bold) for H2-H3
- **Usage:** Page titles, section headers, marketing headlines

### Body Text
- **Font:** Space Grotesk
- **Weight:** 400 (Regular)
- **Usage:** Main content, descriptions, general reading

### Navigation & UI
- **Font:** Space Grotesk
- **Weight:** 500 (Medium)
- **Usage:** Menu items, buttons, interactive elements

### Data & Technical
- **Font:** JetBrains Mono
- **Weight:** 400 (Regular) for data, 500 (Medium) for labels
- **Usage:** Consumption logs, analytics, timestamps, codes

---

## Implementation Guidelines

### CSS Font Stack
```css
:root {
  --font-primary: 'Space Grotesk', system-ui, -apple-system, sans-serif;
  --font-mono: 'JetBrains Mono', 'SF Mono', Monaco, Inconsolata, 'Roboto Mono', monospace;
  --font-system: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
}

/* Primary Typography */
.text-primary {
  font-family: var(--font-primary);
}

/* Technical/Data Typography */
.text-mono {
  font-family: var(--font-mono);
}

/* Fallback Typography */
.text-system {
  font-family: var(--font-system);
}
```

### Letter Spacing (Tracking)
- **Tight:** `-0.02em` for wordmarks and headlines
- **Normal:** `0` for body text
- **Wide:** `0.05em` for technical labels and small caps

### Line Height
- **Headlines:** 1.1 - 1.2 for impact
- **Body text:** 1.5 - 1.6 for readability
- **Technical:** 1.4 for data clarity

---

## Brand Voice Alignment

### Space Grotesk Personality Match
- **Wise:** Mature, sophisticated letterforms
- **Calm:** Geometric stability, even rhythm
- **Trustworthy:** Professional appearance, clear readability
- **Approachable:** Subtle humanist touches, friendly curves
- **Grounded:** Solid construction, reliable spacing

### JetBrains Mono Technical Integration
- **Evidence-based:** Precise, data-focused presentation
- **Transparent:** Clear, unambiguous information display
- **Professional:** Technical credibility for health data

---

## Accessibility Considerations

### Readability
- All fonts tested for dyslexia-friendly characteristics
- Sufficient character spacing for easy scanning
- Clear distinction between similar characters (0/O, 1/I)

### Size Minimums
- **Body text:** Minimum 16px for mobile, 14px for desktop
- **Technical text:** Minimum 14px for data displays
- **Navigation:** Minimum 16px for touch targets

### Contrast Requirements
- All text meets WCAG 2.1 AA contrast ratios
- Technical data uses enhanced contrast for precision reading
- Sufficient spacing between text elements

---

## Platform-Specific Considerations

### Mobile Applications
- Font sizes optimized for touch interfaces
- Adequate spacing for thumb navigation
- Readable at high pixel densities

### Web Applications
- Font loading optimization with display: swap
- Progressive enhancement with system font fallbacks
- Responsive sizing across device breakpoints

### Marketing Materials
- High-impact headline treatments
- Consistent brand voice across all materials
- Print-optimized font rendering

---

*This typography system ensures consistent, accessible, and brand-aligned text presentation across all CannaNote touchpoints while supporting both human readability and technical data precision.*