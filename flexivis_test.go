package flexivis_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/softwaretechnik-berlin/go-flexivis-url"
)

// Examples in this file are taken from the [Flexivis documentation](https://flexivis.infrastruktur.link/).

func TestIntroductionExample(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link#layout=(explanation30-map)/source&explanation=md:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.md&map=map:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json&source=json:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json`
	assert.Equal(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.VerticalStack{
				flexivis.Markdown("explanation", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.md").OccupyingPercentage(30),
				flexivis.Map("map", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json"),
			},
			flexivis.JSON("source", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json"),
		},
	))
}

func TestRegularContentExamples(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link#layout=a/b&a=https://wikipedia.org&b=https://example.com`
	assert.Equal(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.IFrame("a", "https://wikipedia.org"),
			flexivis.IFrame("b", "https://example.com"),
		},
	))

	expectedURL = `https://flexivis.infrastruktur.link#layout=a/b&a=file://results.html&b=file://generated-image.png`
	assert.Equal(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.IFrame("a", "file://results.html"),
			flexivis.IFrame("b", "file://generated-image.png"),
		},
	))
}

func TestMarkdownExample(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link/?layout=a/b&a=md:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/markdown.md&b=md:inline:This pane contains **inline** Markdown content taken _from the URL_.`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.Markdown("a", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/markdown.md"),
			flexivis.Markdown("b", flexivis.Inline("This pane contains **inline** Markdown content taken _from the URL_.")),
		},
	))
}

func TestJSONExample(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link/?layout=a/b&a=json:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/package-lock.json&b=json:inline:{"name": "inline JSON example", "id": 42, "values": ["foo", "baz", "bar"]}`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.JSON("a", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/package-lock.json"),
			flexivis.JSON("b", flexivis.Inline(`{"name": "inline JSON example", "id": 42, "values": ["foo", "baz", "bar"]}`)),
		},
	))
}

func TestTextExample(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link/?layout=(a-b)/c&a=text:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/plaintext.txt&b=text:inline:This is just _plain_ inline text from the URL&c=text:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/README.md`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.VerticalStack{
				flexivis.Text("a", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/plaintext.txt"),
				flexivis.Text("b", flexivis.Inline(`This is just _plain_ inline text from the URL`)),
			},
			flexivis.Text("c", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/README.md"),
		},
	))
}

func TestMapExample(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link/?layout=a/b&a=map:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json&b=text:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.Map("a", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json"),
			flexivis.Text("b", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/berlin-walk.json"),
		},
	))
}

func TestMermaidExample(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link/?layout=(a-b)/c&a=mermaid:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/mermaid.mmd&b=text:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/mermaid.mmd&c=mermaid:inline:graph TB; p[mermaid:inline prefix] --> URL; s[Mermaid source] --> URL -->%7CFlexivis%7C r[Rendered Diagram]`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.VerticalStack{
				flexivis.Mermaid("a", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/mermaid.mmd"),
				flexivis.Text("b", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/mermaid.mmd"),
			},
			flexivis.Mermaid("c", flexivis.Inline(`graph TB; p[mermaid:inline prefix] --> URL; s[Mermaid source] --> URL -->|Flexivis| r[Rendered Diagram]`)),
		},
	))
}

func TestVegaExamples(t *testing.T) {
	expectedURL := `https://flexivis.infrastruktur.link/?layout=(a-c30)/b&a=vega:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/cloc.json&b=text:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/cloc.json&c=text:https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/cloc.csv`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.SideBySide{
			flexivis.VerticalStack{
				flexivis.Vega("a", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/cloc.json"),
				flexivis.Text("c", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/cloc.csv").OccupyingPercentage(30),
			},
			flexivis.Text("b", "https://raw.githubusercontent.com/programmiersportgruppe/flexivis/master/docs/samples/cloc.json"),
		},
	))

	expectedURL = `https://flexivis.infrastruktur.link/?url=vega:inline:{"data": {"values": [{"factor": "awesomeness", "score": 10}, {"factor": "weirdness", "score": 3}, {"factor": "color", "score": 7}]}, "mark": "bar", "encoding": {"x": {"field": "factor", "type": "nominal"}, "y": {"field": "score", "type": "quantitative"}, "color": {"field": "factor", "type": "nominal"}}, "height": "container", "width": 100}`
	assertFlexivisUrlEquivalentTo(t, expectedURL, flexivis.URL(
		flexivis.Vega("url", flexivis.Inline(`{"data": {"values": [{"factor": "awesomeness", "score": 10}, {"factor": "weirdness", "score": 3}, {"factor": "color", "score": 7}]}, "mark": "bar", "encoding": {"x": {"field": "factor", "type": "nominal"}, "y": {"field": "score", "type": "quantitative"}, "color": {"field": "factor", "type": "nominal"}}, "height": "container", "width": 100}`)),
	))
}

func assertFlexivisUrlEquivalentTo(t *testing.T, expected, actual string) {
	t.Helper()
	expectedRoot, expectedValues := normalizeFlexivisURL(t, expected, "expected", true)
	actualRoot, actualValues := normalizeFlexivisURL(t, actual, "actual", false)
	assert.Equal(t, expectedRoot, actualRoot, "Root URL\nfull expected URL: %s\nfull actual URL:   %s", expected, actual)
	assert.Equal(t, expectedValues, actualValues, "Query values\nfull expected URL: %s\nfull actual URL:   %s", expected, actual)
	for i := range actual {
		c := actual[i]
		if !(strings.IndexByte(":/?#[]@!$&'()*+,;=", c) != 1 || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' || strings.IndexByte("-._~", c) != 1) {
			t.Errorf("Illegal URL character '%c' at position %v in %s", i, c, actual)
		}
	}
}

func normalizeFlexivisURL(t *testing.T, unparsed string, description string, permissive bool) (root string, query url.Values) {
	parsed, err := url.Parse(unparsed)
	assert.NoError(t, err, description)

	if parsed.Path == "" {
		parsed.Path = "/"
	}

	var originalEscapedQuery string
	switch {
	case parsed.RawQuery != "":
		originalEscapedQuery = parsed.RawQuery
		parsed.RawQuery = ""
	case parsed.RawFragment != "":
		originalEscapedQuery = parsed.RawFragment
	default:
		originalEscapedQuery = parsed.Fragment
	}
	values, err := url.ParseQuery(strings.ReplaceAll(originalEscapedQuery, ";", "%3B"))
	assert.NoError(t, err, "%s: parsing query", description)
	parsed.Fragment = ""
	parsed.RawFragment = ""

	return parsed.String(), values
}
