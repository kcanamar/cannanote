# CannaNote Git Commit Message Convention

## Overview

This document establishes consistent commit message formatting for the CannaNote project. Following these conventions helps maintain clear project history and makes it easier to understand changes across development, design, and brand work.

## Format Structure

```
<type>: <description>

[optional body]

[optional footer]
```

**Example:**
```
feat: implement responsive header with brand colors

Add new header component with Ganjier-inspired color palette.
Includes mobile breakpoints and accessibility improvements.

Closes #123
```

## Commit Types

### **feat:** - New features or additions
New functionality, brand elements, pages, or major additions to the project.

**Examples:**
- `feat: implement responsive header with brand colors`
- `feat: add cannabis strain tracking interface`
- `feat: create brand voice documentation system`
- `feat: integrate Supabase authentication flow`

### **fix:** - Bug fixes or corrections
Resolving broken functionality, correcting errors, or fixing issues.

**Examples:**
- `fix: resolve broken link in brand guidelines PDF`
- `fix: correct database connection timeout issues`
- `fix: repair mobile navigation dropdown behavior`
- `fix: update outdated terpene information in content`

### **docs:** - Documentation updates
Changes to documentation including brand style guides, READMEs, API docs, or educational content.

**Examples:**
- `docs: update typography section in brand book`
- `docs: add deployment instructions to engineering guide`
- `docs: revise cannabis science content for accuracy`
- `docs: create onboarding documentation for new developers`

### **style:** - Formatting or aesthetic changes
Visual improvements, CSS updates, design tweaks that don't affect functionality.

**Examples:**
- `style: adjust spacing in logo usage examples`
- `style: improve button hover states across interface`
- `style: standardize color values in CSS variables`
- `style: refine typography hierarchy in templates`

### **refactor:** - Code or structure improvements
Improving code structure, organization, or architecture without changing behavior.

**Examples:**
- `refactor: reorganize asset folders for better modularity`
- `refactor: restructure HTTP handlers using hexagonal architecture`
- `refactor: simplify template component structure`
- `refactor: extract reusable brand color functions`

### **perf:** - Performance optimizations
Improvements that make the application faster or more efficient.

**Examples:**
- `perf: compress images in brand kit for quicker web delivery`
- `perf: optimize database queries for strain search`
- `perf: implement lazy loading for cannabis education content`
- `perf: reduce bundle size by removing unused dependencies`

### **test:** - Adding or updating tests
Test creation, updates, or validation scripts for code or brand consistency.

**Examples:**
- `test: add unit tests for brand color palette generator`
- `test: create integration tests for user authentication`
- `test: implement visual regression tests for brand components`
- `test: add validation for cannabis data entry forms`

### **chore:** - Maintenance tasks
Dependency updates, tooling changes, or general maintenance work.

**Examples:**
- `chore: update npm packages for design system`
- `chore: configure development environment for new team members`
- `chore: update Go dependencies to latest versions`
- `chore: clean up unused brand asset files`

### **wip:** - Work-in-progress commits
Temporary commits for feature branches (not for main branch merges).

**Examples:**
- `wip: sketch initial wireframes for rebranding`
- `wip: prototype cannabis tracking interface`
- `wip: draft brand messaging framework`
- `wip: experiment with HTMX form interactions`

### **asset:** - Brand asset updates
Changes to images, fonts, icons, or other brand-specific files.

**Examples:**
- `asset: add high-res versions of primary logo`
- `asset: update cannabis strain icons with new design`
- `asset: compress and optimize brand photography`
- `asset: include new Ganjier-inspired color swatches`

### **brand:** - High-level brand strategy changes
Major brand guideline revisions, identity shifts, or strategic updates.

**Examples:**
- `brand: revise tone-of-voice section in guidelines`
- `brand: update positioning to emphasize Ganjier approach`
- `brand: establish new visual identity standards`
- `brand: refine target audience definitions`

### **ci:** - Continuous integration/config changes
Build scripts, deployment configurations, or development tooling updates.

**Examples:**
- `ci: add linting for SCSS files in brand theme`
- `ci: configure automated deployment pipeline for Fly.io`
- `ci: set up brand asset validation in build process`
- `ci: implement automated testing for template compilation`

### **revert:** - Reverting previous commits
Undoing previous changes (can be auto-generated or manual).

**Examples:**
- `revert: undo faulty color scheme update`
- `revert: remove problematic database migration`
- `revert: restore previous brand voice guidelines`
- `revert: rollback breaking authentication changes`

## Writing Guidelines

### **Description Format**
- Use imperative mood ("add" not "added" or "adds")
- Keep first line under 72 characters
- Don't end with punctuation
- Be specific and descriptive

### **Body Guidelines** (when needed)
- Explain the "what" and "why", not the "how"
- Reference issue numbers if applicable
- Include breaking changes if relevant
- Keep lines under 72 characters

### **Good Examples:**
```
feat: implement cannabis strain autocomplete search

Add autocomplete functionality to strain name input fields.
Uses Supabase full-text search with user history prioritization.

Closes #45
```

```
brand: establish Ganjier-inspired voice guidelines

Update brand voice from "knowledgeable friend" to "educated 
stoner who reads research" to better align with Ganjier 
certification approach and combat cannabis misinformation.
```

```
fix: resolve mobile navigation menu z-index issue

Navigation dropdown was appearing behind modal overlays on 
mobile devices. Increased z-index value to ensure proper 
layer stacking.
```

### **Avoid:**
```
❌ fix: stuff
❌ update: changed some things
❌ feat: new feature added
❌ docs: updated documentation
```

## Usage Notes

- **This is not automated** - These conventions are documentation only
- All team members should follow these formats manually
- Use conventional types even for small changes
- When in doubt, err on the side of being more descriptive
- Branch commits can be more casual; clean up for main branch merges

## Integration with Development Workflow

### **Feature Development**
1. Use `wip:` commits during active development
2. Use appropriate type for final feature commit
3. Squash or clean up commit history before merging to main

### **Brand Work**
1. Use `brand:`, `asset:`, or `docs:` for most brand-related changes
2. Use `style:` for pure visual improvements
3. Use `feat:` when adding substantial new brand elements

### **Bug Fixes and Maintenance**
1. Use `fix:` for any broken functionality
2. Use `chore:` for dependency updates and cleanup
3. Use `refactor:` when improving code without changing behavior

This convention ensures clear communication about all changes to the CannaNote project, whether they're technical improvements, brand development, or content updates.