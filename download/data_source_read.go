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
		dataSource := &DataSource{RESTMap: DataSourceInfoMap[dsName].MarshallHCL(config, dlConfig)}
		for _, v := range dataSource.RESTMap {
			v.UniqueName = escape(UniqueDSName(dsName, v.Values))
		}
		me[dsName] = dataSource
	}

	return nil
}
