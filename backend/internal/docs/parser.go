package docs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// Doc represents a parsed markdown document
type Doc struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	SidebarLabel string   `yaml:"sidebar_label"`
	SidebarOrder int      `yaml:"sidebar_order"`
	Section     string    `yaml:"section"`
	Keywords    []string  `yaml:"keywords"`
	RelatedPages []string `yaml:"related_pages"`
	LastUpdated time.Time `yaml:"last_updated"`
	ContentHTML string    `yaml:"-"`
	TOC         []TOCItem `yaml:"-"`
	Path        string    `yaml:"-"`
}

// TOCItem represents a table of contents entry
type TOCItem struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Level int    `json:"level"`
}

// SidebarItem represents a navigation item in the docs sidebar
type SidebarItem struct {
	Title    string        `json:"title"`
	Label    string        `json:"label"`
	Href     string        `json:"href"`
	Order    int           `json:"order"`
	Section  string        `json:"section"`
	Children []SidebarItem `json:"children"`
	Active   bool          `json:"active"`
}

// DocsCache manages parsed documents and navigation
type DocsCache struct {
	docs     map[string]*Doc
	sidebar  []SidebarItem
	lastMod  map[string]time.Time
	mu       sync.RWMutex
	markdown goldmark.Markdown
}

// NewDocsCache creates a new docs cache instance
func NewDocsCache() *DocsCache {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Linkify,
			extension.Strikethrough,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &DocsCache{
		docs:     make(map[string]*Doc),
		lastMod:  make(map[string]time.Time),
		markdown: md,
	}
}

// LoadContent loads all markdown files from the content directory
func (c *DocsCache) LoadContent(contentDir string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Convert absolute path to relative path from content dir
		relPath, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}

		// Convert file path to URL path (remove .md extension)
		urlPath := strings.TrimSuffix(relPath, ".md")
		if urlPath == "index" {
			urlPath = ""
		}
		urlPath = strings.ReplaceAll(urlPath, "\\", "/") // Windows compatibility

		doc, err := c.parseFile(path)
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", path, err)
		}

		doc.Path = urlPath
		c.docs[urlPath] = doc
		c.lastMod[path] = info.ModTime()

		return nil
	})

	if err != nil {
		return err
	}

	c.buildSidebar()
	return nil
}

// GetDoc retrieves a document by path
func (c *DocsCache) GetDoc(path string) (*Doc, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	doc, exists := c.docs[path]
	return doc, exists
}

// GetSidebar returns the navigation sidebar
func (c *DocsCache) GetSidebar() []SidebarItem {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.sidebar
}

// parseFile parses a markdown file with frontmatter
func (c *DocsCache) parseFile(filePath string) (*Doc, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Split frontmatter and content
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter in %s", filePath)
	}

	// Parse frontmatter
	var doc Doc
	if err := yaml.Unmarshal(parts[1], &doc); err != nil {
		return nil, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	// Parse markdown content
	var buf bytes.Buffer
	if err := c.markdown.Convert(parts[2], &buf); err != nil {
		return nil, fmt.Errorf("error converting markdown: %w", err)
	}

	doc.ContentHTML = buf.String()
	doc.TOC = c.extractTOC(doc.ContentHTML)

	return &doc, nil
}

// extractTOC extracts table of contents from HTML
func (c *DocsCache) extractTOC(html string) []TOCItem {
	var toc []TOCItem
	
	// Regex to match headings with IDs
	headingRegex := regexp.MustCompile(`<h([2-6])[^>]*id="([^"]*)"[^>]*>([^<]*)</h[2-6]>`)
	matches := headingRegex.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) == 4 {
			level := int(match[1][0] - '0') // Convert '2'-'6' to 2-6
			id := match[2]
			text := strings.TrimSpace(match[3])
			
			toc = append(toc, TOCItem{
				ID:    id,
				Text:  text,
				Level: level,
			})
		}
	}

	return toc
}

// buildSidebar constructs the navigation sidebar from loaded documents
func (c *DocsCache) buildSidebar() {
	sections := make(map[string][]SidebarItem)
	
	for _, doc := range c.docs {
		if doc.Section == "" {
			doc.Section = "root"
		}

		item := SidebarItem{
			Title:   doc.Title,
			Label:   doc.SidebarLabel,
			Href:    "/docs/" + doc.Path,
			Order:   doc.SidebarOrder,
			Section: doc.Section,
		}

		if item.Label == "" {
			item.Label = doc.Title
		}

		sections[doc.Section] = append(sections[doc.Section], item)
	}

	// Sort items within each section
	for section := range sections {
		sort.Slice(sections[section], func(i, j int) bool {
			return sections[section][i].Order < sections[section][j].Order
		})
	}

	// Build final sidebar structure
	var sidebar []SidebarItem

	// Add root items first (like index page)
	if rootItems, exists := sections["root"]; exists {
		sidebar = append(sidebar, rootItems...)
	}

	// Define section order
	sectionOrder := []string{"getting-started", "guides", "privacy", "reference", "community"}
	
	for _, sectionName := range sectionOrder {
		if items, exists := sections[sectionName]; exists {
			// Create section header if there are items
			if len(items) > 0 {
				sidebar = append(sidebar, items...)
			}
		}
	}

	c.sidebar = sidebar
}

// RefreshDoc refreshes a single document if it has been modified
func (c *DocsCache) RefreshDoc(contentDir, path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullPath := filepath.Join(contentDir, path+".md")
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}

	lastMod, exists := c.lastMod[fullPath]
	if exists && info.ModTime().Equal(lastMod) {
		return nil // No changes
	}

	doc, err := c.parseFile(fullPath)
	if err != nil {
		return err
	}

	doc.Path = path
	c.docs[path] = doc
	c.lastMod[fullPath] = info.ModTime()
	c.buildSidebar()

	return nil
}