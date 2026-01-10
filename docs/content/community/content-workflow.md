---
title: "Documentation Content Workflow"
description: "Guide for adding, updating, and removing content in the CannaNote documentation system."
sidebar_label: "Content Workflow"
sidebar_order: 4
section: "community"
keywords: ["documentation", "content", "workflow", "markdown", "frontmatter"]
related_pages: ["community/contributing"]
---

# Documentation Content Workflow

This guide explains how to manage content in CannaNote's markdown-driven documentation system.

## Content Structure Overview

All documentation content is stored in markdown files with YAML frontmatter in the `/docs/content/` directory:

```
docs/content/
├── index.md                     # Homepage (/docs)
├── getting-started/
│   ├── index.md                 # (/docs/getting-started)
│   ├── first-entry.md          # (/docs/getting-started/first-entry)
│   └── patterns.md             # (/docs/getting-started/patterns)
├── guides/
│   ├── cannabinoids.md         # (/docs/guides/cannabinoids)
│   ├── terpenes.md             # (/docs/guides/terpenes)
│   └── harm-reduction.md       # (/docs/guides/harm-reduction)
├── privacy/
│   ├── data-privacy.md         # (/docs/privacy/data-privacy)
│   └── export.md               # (/docs/privacy/export)
├── reference/
│   └── api.md                  # (/docs/reference/api)
└── community/
    ├── contributing.md         # (/docs/community/contributing)
    └── content-workflow.md     # (/docs/community/content-workflow)
```

## URL Structure

The documentation system supports two URL patterns:

### Direct Files
- **File**: `docs/content/guides/cannabinoids.md`
- **URL**: `/docs/guides/cannabinoids`

### Directory Index Files
- **File**: `docs/content/getting-started/index.md`
- **URL**: `/docs/getting-started` (automatically resolves to index)

## Frontmatter Schema

Every markdown file must include YAML frontmatter at the top:

```yaml
---
title: "Page Title"                    # Required: Used in <title> and <h1>
description: "Page description"        # Required: Used for SEO meta description
sidebar_label: "Sidebar Label"        # Optional: Label in navigation (defaults to title)
sidebar_order: 1                      # Required: Order within section (lower = higher)
section: "guides"                     # Required: Section for grouping in sidebar
keywords: ["keyword1", "keyword2"]   # Optional: For search and SEO
related_pages: ["page1", "page2"]    # Optional: Show related content
last_updated: "2024-01-02"           # Optional: Content freshness indicator
---
```

### Frontmatter Fields Explained

#### Required Fields

- **`title`**: Page title displayed in browser tab and as main heading
- **`description`**: SEO meta description and page subtitle
- **`sidebar_order`**: Numeric order within section (1 = first, 2 = second, etc.)
- **`section`**: Groups pages in sidebar navigation

#### Valid Section Values

- `"root"` - Homepage and root-level pages
- `"getting-started"` - Getting started guides
- `"guides"` - Cannabis education guides  
- `"privacy"` - Privacy and data handling
- `"reference"` - Technical documentation
- `"community"` - Community and contribution guides

#### Optional Fields

- **`sidebar_label`**: Custom label for sidebar (use when title is too long)
- **`keywords`**: Array of keywords for search functionality
- **`related_pages`**: Array of page paths to show as related content
- **`last_updated`**: ISO date string for content freshness

## Adding New Content

### 1. Create the Markdown File

Create a new `.md` file in the appropriate section directory:

```bash
# Example: Adding a new terpenes guide
touch docs/content/guides/terpenes.md
```

### 2. Add Required Frontmatter

Start the file with complete frontmatter:

```markdown
---
title: "Understanding Terpenes"
description: "Learn about aromatic compounds in cannabis and their effects on your experience."
sidebar_label: "Terpenes"
sidebar_order: 2
section: "guides"
keywords: ["terpenes", "aromatic", "effects", "cannabis"]
related_pages: ["guides/cannabinoids", "guides/methods"]
---

# Understanding Terpenes

Your content here...
```

### 3. Write Content

Use standard markdown syntax with these conventions:

#### Headings
- Use `# Title` for the main heading (matches frontmatter title)
- Use `##` for major sections
- Use `###` for subsections
- Headings automatically generate table of contents entries

#### Links
- **Internal links**: `[Getting Started](../getting-started)` or `[Privacy](../privacy/data-privacy)`
- **External links**: `[GitHub](https://github.com/kcanamar/cannanote)`

#### Code Blocks
```markdown
\`\`\`javascript
// Example code
const example = "hello world";
\`\`\`
```

#### Callouts (Future Enhancement)
```markdown
> **Note**: Important information for users
```

### 4. Test Locally

The docs system automatically loads new files:

1. **Save the file** - No restart needed
2. **Visit the URL** - Navigate to `/docs/your-new-page`
3. **Check sidebar** - Ensure it appears in the correct section and order

### 5. Verify Navigation

Ensure the new page appears in the sidebar:
- Check the section grouping
- Verify the order within the section
- Test the sidebar link works with HTMX

## Updating Existing Content

### 1. Edit the Markdown File

Simply edit the `.md` file - changes are reflected immediately.

### 2. Update Frontmatter if Needed

- **Change position**: Update `sidebar_order`
- **Move sections**: Update `section` value
- **Update SEO**: Modify `title`, `description`, or `keywords`

### 3. Test Changes

Visit the page to verify:
- Content renders correctly
- Sidebar navigation updates
- Related pages still work

## Removing Content

### 1. Delete the File

```bash
rm docs/content/section/page-name.md
```

### 2. Update Related Pages

Check other pages that might reference the deleted page:
- Remove from `related_pages` arrays
- Update any internal links
- Remove from any manual navigation lists

### 3. Consider Redirects

For public pages, consider:
- Adding redirect logic in the handler
- Updating external links
- Notifying users of removed content

## Directory Organization

### Creating New Sections

To add a new documentation section:

1. **Create directory**: `mkdir docs/content/new-section`
2. **Add to parser**: Update `sectionOrder` in `internal/docs/parser.go`
3. **Add content**: Create markdown files with `section: "new-section"`

### Section Guidelines

- **Keep sections focused**: Each section should have a clear purpose
- **Logical grouping**: Group related content together
- **Consistent naming**: Use kebab-case for directories and files
- **Reasonable size**: Aim for 3-7 pages per section

## Content Best Practices

### Writing Guidelines

- **Clear titles**: Make page purposes obvious
- **Descriptive descriptions**: Help users understand content value
- **Logical flow**: Order pages from basic to advanced
- **Consistent tone**: Match CannaNote's educational, non-judgmental voice

### Technical Guidelines

- **Valid markdown**: Ensure proper syntax
- **Complete frontmatter**: Include all required fields
- **Relative links**: Use relative paths for internal links
- **Image optimization**: Optimize any images for web

### SEO Considerations

- **Unique titles**: Each page needs a unique, descriptive title
- **Meta descriptions**: Write compelling 150-160 character descriptions
- **Keywords**: Include relevant search terms naturally
- **Internal linking**: Link between related pages

## Development Workflow

### Local Development

1. **Edit content**: Make changes to markdown files
2. **Test immediately**: No build step required
3. **Check logs**: Monitor server output for any errors
4. **Verify navigation**: Test both direct URLs and sidebar navigation

### Content Review Process

1. **Technical review**: Ensure markdown and frontmatter are valid
2. **Content review**: Verify accuracy and tone
3. **Navigation testing**: Test all links and navigation
4. **SEO check**: Review titles, descriptions, and keywords

## Troubleshooting

### Common Issues

**Page not appearing in sidebar**:
- Check `section` value matches existing sections
- Verify `sidebar_order` is set
- Ensure frontmatter is valid YAML

**404 errors**:
- Verify file path matches URL structure
- Check for typos in frontmatter
- Ensure file has `.md` extension

**Broken links**:
- Use relative paths for internal links
- Test all links after content changes
- Check `related_pages` array for valid paths

**HTMX not working**:
- Verify page loads correctly with direct URL
- Check browser developer tools for JavaScript errors
- Ensure HTMX script is loading properly

### Getting Help

- **Issues**: Report bugs on [GitHub Issues](https://github.com/kcanamar/cannanote/issues)
- **Questions**: Ask in the [Community Support](support) section
- **Contributions**: See [Contributing Guidelines](contributing)

## Technical Details

### Markdown Processing

The system uses [goldmark](https://github.com/yuin/goldmark) with these features:
- **GitHub Flavored Markdown**: Tables, task lists, strikethrough
- **Syntax highlighting**: Code blocks with language detection
- **Auto-heading IDs**: Automatic anchor links for headings
- **Table of contents**: Generated from heading structure

### Caching

- **File watching**: Changes detected automatically in development
- **Memory cache**: Parsed content cached for performance
- **Sidebar generation**: Navigation rebuilt when content changes

### Performance

- **Server-side rendering**: Full HTML for SEO and speed
- **HTMX partial loading**: Instant navigation between pages
- **Asset optimization**: CSS and JavaScript optimized for production

This documentation workflow ensures maintainable, SEO-friendly content with excellent user experience.