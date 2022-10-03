// Package flexivis provides utilities for programatically creating Flexivis links.
// See https://flexivis.infrastruktur.link/ for information about Flexivis.
package flexivis

import (
	"fmt"
	"strings"
)

// Layout is a layout string as documented at https://flexivis.infrastruktur.link/#layout.
type Layout string

// URL creates a Flexivis URL using this layout and the given views
func (l Layout) URL(views []View) string {
	// We have our own URL encoding code so we can customise how spaces are encoded and how much percent encoding of reserved characters happens.
	var b strings.Builder
	b.WriteString("https://flexivis.infrastruktur.link#")
	hasWrittenParam := false
	if l != "" {
		b.WriteString("layout=")
		b.WriteString(string(l))
		hasWrittenParam = true
	}
	for _, view := range views {
		if hasWrittenParam {
			b.WriteByte('&')
		} else {
			hasWrittenParam = true
		}
		b.WriteString(string(view.Name))
		b.WriteByte('=')
		if view.Type != "" {
			b.WriteString(string(view.Type))
			b.WriteByte(':')
		}
		b.WriteString(escapeParameterValue(string(view.Resource)))
	}
	return b.String()
}

func parenthesize(layout Layout) Layout {
	return Layout("(" + string(layout) + ")")
}

type outerLayoutStructure byte

const (
	atomic outerLayoutStructure = iota
	scaled
	joined
)

// StructuredLayout represents a programatically created structured layout.
// In particular, you can mix and match SideBySide, VerticalStack, Sized and View objects to build a strcutured layout.
// See https://flexivis.infrastruktur.link/#layout for an indication of how these can be combined.
type StructuredLayout interface {
	layoutAndEmbeddedViews() (Layout, outerLayoutStructure, []View)
}

var _ StructuredLayout = (*SideBySide)(nil)
var _ StructuredLayout = (*VerticalStack)(nil)
var _ StructuredLayout = (*Sized)(nil)
var _ StructuredLayout = (*ViewName)(nil)
var _ StructuredLayout = (*View)(nil)

// URL builds a Flexivis URL from a StructuredLayout.
func URL(structuredLayout StructuredLayout) string {
	layout, _, views := structuredLayout.layoutAndEmbeddedViews()
	if layout == "url" {
		layout = ""
	}
	return layout.URL(views)
}

// SideBySide represents a collection of sublayouts that are joined into a composite layout by '/' sperators and rendered side-by-side.
// See https://flexivis.infrastruktur.link/#layout.
type SideBySide []StructuredLayout

func (s SideBySide) layoutAndEmbeddedViews() (Layout, outerLayoutStructure, []View) {
	return joinLayoutsAndEmbeddedViews('/', s)
}

// VerticalStack represents a collection of sublayouts that are joined into a composite layout by '-' sperators and rendered on on top of another.
// See https://flexivis.infrastruktur.link/#layout.
type VerticalStack []StructuredLayout

func (s VerticalStack) layoutAndEmbeddedViews() (Layout, outerLayoutStructure, []View) {
	return joinLayoutsAndEmbeddedViews('-', s)
}

// Sized specifies how much of the available vertical or horizontal space a sublayout should use.
// It's meant to be used as a child of SideBySide or VerticalStack layouts.
// See https://flexivis.infrastruktur.link/#layout.
type Sized struct {
	InnerLayout StructuredLayout
	Percentage  uint8
}

func (s Sized) layoutAndEmbeddedViews() (Layout, outerLayoutStructure, []View) {
	atomic, embedded := ensureOuterStructureAtMost(atomic, s.InnerLayout)
	return Layout(fmt.Sprintf("%s%v", atomic, s.Percentage)), scaled, embedded
}

// ViewName represents the name of a view.
type ViewName string

// AsLayout converts this view name to a layout. (A single view name is a valid layout).
func (n ViewName) AsLayout() Layout { return Layout(n) }
func (n ViewName) layoutAndEmbeddedViews() (Layout, outerLayoutStructure, []View) {
	return n.AsLayout(), atomic, nil
}

// ViewType is a type of view supported by Flexivis.
// See https://flexivis.infrastruktur.link/#view-types
type ViewType string

type Resource string

func URLResource(url string) Resource {
	return Resource(url)
}

func Inline(contents string) Resource {
	return Resource("inline:" + contents)
}

// View represents a view.
type View struct {
	Name     ViewName
	Type     ViewType // this is a simplified modeling of the more complex prefix structure documented in https://flexivis.infrastruktur.link/#view-specifications
	Resource Resource
}

func (v View) layoutAndEmbeddedViews() (Layout, outerLayoutStructure, []View) {
	return v.Name.AsLayout(), atomic, []View{v}
}

// IFrame creates an IFrame to display regular content.
// See https://flexivis.infrastruktur.link/#regular-content.
func IFrame(name ViewName, url string) View {
	if !(strings.HasPrefix(url, "https:") || strings.HasPrefix(url, "http:") || strings.HasPrefix(url, "file:")) {
		panic(fmt.Sprintf("IFrame content must have a URL with an https, http or file scheme, but got: %s", url))
	}
	return View{name, "", Resource(url)}
}

// JSON creates an interactive formatted JSON viewer with collapsible nodes.
// See https://flexivis.infrastruktur.link/#json.
func JSON(name ViewName, json Resource) View {
	return View{name, "json", json}
}

// Map creates an interactive map view from inline GeoJSON.
// See https://flexivis.infrastruktur.link/#map.
func Map(name ViewName, geoJSON Resource) View {
	return View{name, "map", geoJSON}
}

// Markdown creates view with rendered [markdown](https://en.wikipedia.org/wiki/Markdown).
// See https://flexivis.infrastruktur.link/#markdown.
func Markdown(name ViewName, markdown Resource) View {
	return View{name, "md", markdown}
}

// Mermaid creates view with rendered [Mermaid diagrams](https://mermaid-js.github.io/).
// See https://flexivis.infrastruktur.link/#mermaid.
func Mermaid(name ViewName, markdown Resource) View {
	return View{name, "mermaid", markdown}
}

// Text creates plain text view.
// See https://flexivis.infrastruktur.link/#text.
func Text(name ViewName, markdown Resource) View {
	return View{name, "text", markdown}
}

// Text creates view with a rendered [Vega](https://vega.github.io/vega/) or [Vega-Lite](https://vega.github.io/vega-lite/) diagram.
// See https://flexivis.infrastruktur.link/#vega.
func Vega(name ViewName, markdown Resource) View {
	return View{name, "vega", markdown}
}

func ensureOuterStructureAtMost(max outerLayoutStructure, structuredLayout StructuredLayout) (Layout, []View) {
	layout, outerStructure, embeddedViews := structuredLayout.layoutAndEmbeddedViews()
	if outerStructure <= max {
		return layout, embeddedViews
	}
	return parenthesize(layout), embeddedViews
}

func joinLayoutsAndEmbeddedViews(joiner byte, structuredLayouts []StructuredLayout) (Layout, outerLayoutStructure, []View) {
	switch len(structuredLayouts) {
	case 0:
		panic("Need at least 1 structured layout but got 0")
	case 1:
		return structuredLayouts[0].layoutAndEmbeddedViews()
	default:
		var joinedLayout strings.Builder
		var joinedViews []View
		for i, structuredLayout := range structuredLayouts {
			if i != 0 {
				joinedLayout.WriteByte(joiner)
			}
			layout, views := ensureOuterStructureAtMost(scaled, structuredLayout)
			joinedLayout.WriteString(string(layout))
			joinedViews = append(joinedViews, views...)
		}
		return Layout(joinedLayout.String()), joined, joinedViews
	}
}

func escapeParameterValue(value string) string {
	hasSpace := false
	charactersToPercentEncode := 0
	for i := 0; i < len(value); i++ {
		c := value[i]
		if shouldEscape(c) {
			if c == ' ' {
				hasSpace = true
			} else {
				charactersToPercentEncode++
			}
		}
	}
	if charactersToPercentEncode == 0 && !hasSpace {
		return value
	}

	encoded := make([]byte, len(value)+2*charactersToPercentEncode)
	encodedIndex := 0
	emit := func(c byte) {
		encoded[encodedIndex] = c
		encodedIndex++
	}
	for i := range value {
		c := value[i]
		if !shouldEscape(c) {
			emit(c)
			continue
		}
		if c == ' ' {
			emit('+')
			continue
		}
		emit('%')
		emit("0123456789ABCDEF"[c>>4])
		emit("0123456789ABCDEF"[c&0xf])
	}
	return string(encoded)
}

func shouldEscape(c byte) bool {
	// Don't escape the alphanumeric [Unreserved Characters](https://www.rfc-editor.org/rfc/rfc3986#section-2.3).
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~': // Don't escape the symbolic [Unreserved Characters](https://www.rfc-editor.org/rfc/rfc3986#section-2.3)
		return false
	case ':', '/', '?', '#', '[', ']', '@', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=': // Handle the [Reserved Characters](https://www.rfc-editor.org/rfc/rfc3986#section-2.2).
		switch c {
		case '&': // Escape the parameter separator.
			return true
		case '+': // Escape the plus character because Flexivis treats a literal plus as an encoded space.
			return true
		case '#': // Escape the number sign because iTerm2 doesn't recognise URLs with this character in their fragments.
			return true
			// case ' ', '[', ']', '{', '}': // Escape characters that cause URL recognition logic in e.g. terminal emulators to fail to associate all of the characters as being part of a URL
			// 	return true
		}
		// Leave all other reserved characters unescaped
		return false
	}
	// Escape all other characters (include the precent sign itself).
	return true
}
