// trait_output_formats
package models

type TraitLayer struct {
	Name     string  `json:"name"`
	Total    int     `json:"trait_count"`
	Elements []Trait `json:"elements"`
}

type Trait struct {
	Name  string `json:"name"`
	Total int    `json:"show_count"`
	Rate  string `json:"rate"`
}
