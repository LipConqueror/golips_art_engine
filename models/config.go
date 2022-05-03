// config
package models

import (
	"encoding/json"
)

type Config struct {
	NamePrefix        string           `json:"namePrefix"`
	Description       string           `json:"description"`
	BaseUri           string           `json:"baseUri"`
	Format            OutputFormat     `json:"format"`
	Background        Background       `json:"background"`
	RarityDelimiter   string           `json:"rarityDelimiter"`
	ColorSetDelimiter string           `json:"colorSetDelimiter"`
	LimitDelimiter    string           `json:"limitDelimiter"`
	DnaDelimiter      string           `json:"dnaDelimiter"`
	DnaSettings       DnaSettings      `json:"dnaSettings"`
	MetadataSettings  MetadataSettings `json:"metadataSettings"`
	ProcessCount      json.Number      `json:"processCount"`
	LogSettings       LogSettings      `json:"logSettings"`

	MultiVersionSettings MultiVersionSettings `json:"multiVersionSettings"`

	SolanaMetadata SolanaMetadataSettings `json:"solanaMetadata"`

	LayerConfigurations []LayerConfiguration `json:"layerConfigurations"`
}

type MultiVersionSettings struct {
	LayerName string `json:"layerName"`
}

type SolanaMetadataSettings struct {
	Symbol               string          `json:"symbol"`                  // solana
	SellerFeeBasisPoints json.Number     `json:"seller_fee_basis_points"` // solana
	ExternalUrl          string          `json:"external_url"`            // solana
	Creators             []SolanaCreator `json:"creators"`
}

type SolanaCreator struct {
	Address string      `json:"address"`
	Share   json.Number `json:"share"`
}

type DnaSettings struct {
	SaveDnaHistory     bool   `json:"saveDnaHistory"`
	LoadDnaHistory     bool   `json:"loadDnaHistory"`
	LoadDnaHistoryName string `json:"loadDnaHistoryName"`
	StartId            int    `json:"startId"`
}

type MetadataSettings struct {
	SaveDnaInMetadata  bool   `json:"saveDnaInMetadata"`
	ShowNoneInMetadata bool   `json:"showNoneInMetadata"`
	NoneAttributeName  string `json:"noneAttributeName"`

	ExtraMetadata *ExtraMetadata `json:"extraMetadata"`

	OutputEthFormat bool `json:"outputEthFormat"`
	OutputSOLFormat bool `json:"outputSolanaFormat"`

	ShowEditionInMetadata bool `json:"showEditionInMetadata"`
}

type LogSettings struct {
	ShowGeneratingProgress bool `json:"showGeneratingProgress"`
	Debug                  bool `json:"debug"`
}

type LayerConfiguration struct {
	GrowEditionSizeTo int               `json:"growEditionSizeTo"`
	LayersOrder       []LayerOrder      `json:"layersOrder"`
	ConflictElements  map[string]string `json:"conflictElements"`
	ColorSets         map[string]string `json:"-"` // k-v: colorSet-color ie: hair-red
}

type LayerOrder struct {
	Name     string                    `json:"name"`
	Options  LayerOption               `json:"options"`
	Elements []LayerElement            `json:"-"`
	Limits   map[string][]LayerElement `json:"-"`
}

type LayerOption struct {
	BypassDNA      bool   `json:"bypassDNA"`
	DisplayName    string `json:"displayName"`
	IsColorBase    bool   `json:"isColorBase"`
	ColorSet       string `json:"colorSet"`
	HideInMetadata bool   `json:"hideInMetadata"`
}

type LayerElement struct {
	Id              int
	Name            string
	Color           string
	Path            string
	Weight          float64
	BelongLayerName string
	HideInMetadata  bool
}

type OutputFormat struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	Smoothing bool `json:"smoothing"`
}

type Background struct {
	Generate      bool    `json:"generate"`
	Brightness    string  `json:"brightness"`
	Static        bool    `json:"static"`
	Default       string  `json:"default"`
	BrightnessNum float64 `json:"-"`
}
