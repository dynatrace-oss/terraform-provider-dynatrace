package openpipeline

type RoutingTable struct {
	CatchAllPipeline RoutingTableEntryTarget `json:"catchAllPipeline"`
	Editable         *bool                   `json:"editable,omitempty"`
	Entries          []RoutingTableEntry     `json:"entries"`
}

type RoutingTableEntryTarget struct {
	Editable   *bool  `json:"editable,omitempty"`
	PipelineId string `json:"pipelineId"`
}

type RoutingTableEntry struct {
	Builtin    *bool   `json:"builtin,omitempty"`
	Editable   *bool   `json:"editable,omitempty"`
	Enabled    bool    `json:"enabled"`
	Matcher    Matcher `json:"matcher"`
	Note       string  `json:"note"`
	PipelineId string  `json:"pipelineId"`
}
