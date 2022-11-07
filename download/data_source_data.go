package download

type DataSourceClient interface {
	ListInterface() (interface{}, error)
}

type DataSource struct {
	RESTMap map[string]map[string]interface{}
}

type DataSourceData map[string]DataSource
