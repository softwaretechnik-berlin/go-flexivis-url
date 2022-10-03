package main

import (
	"encoding/json"

	"github.com/softwaretechnik-berlin/go-flexivis-url"
)

const someMarkdown = `
# Example

This markdown will be _included in the URL_ and **rendered into a view**!
`

func main() {
	println(flexivis.URL(
		flexivis.SideBySide{
			flexivis.VerticalStack{
				flexivis.SideBySide{
					flexivis.Markdown("description", flexivis.Inline(someMarkdown)),
					flexivis.Mermaid("diagram", flexivis.Inline(`graph TD; classDef empty stroke:none,fill:none; S( ):::empty -->|structured URL spec in code| G[go-flexivis-url] -->|URL| B[Browser] -->|URL| F[Flexivis] -->|a nicely rendered view| B`)),
				},
				flexivis.IFrame("flexivis", "https://flexivis.infrastruktur.link/"),
			}.OccupyingPercentage(40),
			flexivis.Map("a", flexivis.Inline(
				asJson(jobject{
					"type": "Feature",
					"geometry": jobject{
						"type": "LineString",
						"coordinates": []any{
							[]any{13.3907, 52.5074},
							[]any{13.3902, 52.5076},
							[]any{13.3891, 52.5076},
							[]any{13.3871, 52.5077},
							[]any{13.3855, 52.5073},
							[]any{13.3841, 52.5095},
							[]any{13.3838, 52.5109},
							[]any{13.3827, 52.5136},
							[]any{13.3813, 52.5156},
							[]any{13.3796, 52.5165},
							[]any{13.3785, 52.5163},
						},
					},
					"properties": jobject{
						"stroke":      "green",
						"id":          42,
						"title":       "Berlin Walk",
						"description": "Represents GPS data collected during a hypothetical walk through Berlin.",
						"source":      "handcrafted",
					},
				}),
			)),
		},
	))
}

type jobject = map[string]any

func asJson(value any) string {
	j, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(j)
}
