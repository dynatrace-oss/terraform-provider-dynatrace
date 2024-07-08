package openpipeline

type DataExtractionProcessor struct {
	Description     string                   `json:"description"`
	Editable        *bool                    `json:"editable,omitempty"`
	Enabled         bool                     `json:"enabled"`
	EventProvider   *EventProvider           `json:"eventProvider"`   //bizevent extraction processor
	EventType       *EventProvider           `json:"eventType"`       //bizevent extraction processor
	FieldExtraction *BizeventFieldExtraction `json:"fieldExtraction"` //bizevent extraction processor
	Id              string                   `json:"id"`
	Matcher         Matcher                  `json:"matcher"`
	Properties      *DavisEventProperty      `json:"properties"` // davis extraction processor
	SampleData      *string                  `json:"sampleData,omitempty"`
	Type            string                   `json:"type"`
}

type Matcher = string

type DavisEventProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventProvider struct {
	Constant *string `json:"constant"` //constant value assignment
	Field    *string `json:"field"`    // field value assignment
	Type     *string `json:"type"`     //constant value assignment
}

type BizeventFieldExtraction struct {
	Fields   []string                        `json:"fields"`
	Semantic BizeventFieldExtractionSemantic `json:"semantic"`
}

type BizeventFieldExtractionSemantic string

type MetricExtractionProcessor struct {
	Description string    `json:"description"`
	Dimensions  *[]string `json:"dimensions,omitempty"`
	Editable    *bool     `json:"editable,omitempty"`
	Enabled     bool      `json:"enabled"`
	Field       *string   `json:"field"` //value metric extraction processor
	Id          string    `json:"id"`
	Matcher     Matcher   `json:"matcher"`
	MetricKey   string    `json:"metricKey"`
	SampleData  *string   `json:"sampleData,omitempty"`
	Type        string    `json:"type"`
}

type SxqlProcessor struct {
	Description string  `json:"description"`
	Editable    *bool   `json:"editable,omitempty"`
	Enabled     bool    `json:"enabled"`
	Id          string  `json:"id"`
	Matcher     Matcher `json:"matcher"`
	SampleData  *string `json:"sampleData,omitempty"`
	SxqlScript  string  `json:"sxqlScript"`
	Type        string  `json:"type"`
}

type SecurityContextProcessor struct {
	Description string        `json:"description"`
	Editable    *bool         `json:"editable,omitempty"`
	Enabled     bool          `json:"enabled"`
	Id          string        `json:"id"`
	Matcher     Matcher       `json:"matcher"`
	SampleData  *string       `json:"sampleData,omitempty"`
	Type        string        `json:"type"`
	Value       EventProvider `json:"value"`
}

type StorageStageProcessor struct {
	BucketName  string  `json:"bucketName"`
	Description string  `json:"description"`
	Editable    *bool   `json:"editable,omitempty"`
	Enabled     bool    `json:"enabled"`
	Id          string  `json:"id"`
	Matcher     Matcher `json:"matcher"`
	SampleData  *string `json:"sampleData,omitempty"`
	Type        string  `json:"type"`
}
