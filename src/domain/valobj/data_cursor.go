package valobj

//
type DataCursor struct {
	// 异步处理数据获取器
	BlockSize int64
	BlockNo   int64
	StartTime string
	EndTime   string
}
