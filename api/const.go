package api

const (
	//URL_BJ_HOME               = "http://www.bjbus.com/home/fun_rtbus.php"
	URL_BJ_HOME                     = "http://www.bjbus.com/home/index.php"
	URL_BJ_FMT_LINE_DIRECTION       = "http://www.bjbus.com/home/ajax_search_bus_stop_token.php?act=getLineDirOption&selBLine=%s"
	URL_BJ_FMT_LINE_STATION         = "http://www.bjbus.com/home/ajax_search_bus_stop_token.php?act=getDirStationOption&selBLine=%s&selBDir=%s"
	URL_BJ_FMT_FRESH_STATION_STATUS = "http://www.bjbus.com/home/ajax_search_bus_stop_token.php?act=busTime&selBLine=%s&selBDir=%s&selBStop=%d"
)
