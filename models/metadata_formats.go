// metadata_formats
package models

import (
	"encoding/json"
)

type ExtraMetadata map[string]interface{}

type MetadataErc721 struct {
	Name          string              `json:"name"`
	Description   string              `json:"description,omitempty"`
	Image         string              `json:"image,omitempty"`
	Dna           string              `json:"dna,omitempty"`
	Edition       int                 `json:"edition,omitempty"`
	Date          int64               `json:"date,omitempty"`
	ExtraMetadata string              `json:"extra!@#,omitempty"`
	Attributes    []MetaDataAttribute `json:"attributes"`
	Compiler      string              `json:"compiler,omitempty"`
}

type MetaDataAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type MetadataSolana struct {
	Name                 string              `json:"name"`
	Symbol               string              `json:"symbol"`                  // solana
	SellerFeeBasisPoints json.Number         `json:"seller_fee_basis_points"` // solana
	Description          string              `json:"description,omitempty"`
	Image                string              `json:"image,omitempty"`
	ExternalUrl          string              `json:"external_url"` // solana
	Edition              int                 `json:"edition,omitempty"`
	Dna                  string              `json:"dna,omitempty"`
	ExtraMetadata        string              `json:"extra!@#,omitempty"`
	Attributes           []MetaDataAttribute `json:"attributes"`
	Properties           SolanaProperty      `json:"properties,omitempty"`
	Compiler             string              `json:"compiler,omitempty"`
}

type SolanaProperty struct {
	Category string               `json:"category"`
	Creators []SolanaCreator      `json:"creators"`
	Files    []SolanaPropertyFile `json:"files"`
}

type SolanaPropertyFile struct {
	Uri  string `json:"uri"`
	Type string `json:"type"`
}
