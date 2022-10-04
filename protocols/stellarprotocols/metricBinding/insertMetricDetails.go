package metricbinding

func InsertAndFindMetricID(metricID string, metricName string) (int64, error) {
	// var metricMapID int64
	// object := dao.Connection{}

	// metricMap, errInMetricMap := object.GetMetricMapID(metricID).Then(func(data interface{}) interface{} {
	// 	return data
	// }).Await()
	// if errInMetricMap != nil {
	// 	logrus.Error("Error when retrieving metric id from DB " + errInMetricMap.Error())
	// }
	// if metricMap == nil {
	// 	logrus.Error("Metric ID is not recorded in the DB")
	// 	data, errWhenGettingTheSequence := object.GetNextSequenceValue("METRICID")
	// 	if errWhenGettingTheSequence != nil {
	// 		logrus.Error("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
	// 		return -1, errors.New("Error when taking the sequence no Error : " + errWhenGettingTheSequence.Error())
	// 	}

	// 	//insert to metric map

	// }

	return 0, nil
}
