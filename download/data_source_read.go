package download

func (me DataSourceData) ProcessRead(dlConfig DownloadConfig, resDataMap ResourceData) error {
	for dsName, dsStruct := range DataSourceInfoMap {
		if !dlConfig.MatchDataSource(dsName, resDataMap) && dlConfig.ResourceNames != nil {
			continue
		}

		client := dsStruct.RESTClient(
			dlConfig.EnvironmentURL,
			dlConfig.APIToken,
		)
		if err := me.read(dlConfig, dsName, client); err != nil {
			return err
		}
	}
	return nil
}

func (me DataSourceData) read(dlConfig DownloadConfig, dsName string, client DataSourceClient) error {
	config, err := client.ListInterface()
	if err != nil {
		return err
	}

	if DataSourceInfoMap[dsName].MarshallHCL != nil {
		me[dsName] = DataSource{RESTMap: DataSourceInfoMap[dsName].MarshallHCL(config)}
	}

	return nil
}
