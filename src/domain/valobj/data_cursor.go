package valobj

import jsoniter "github.com/json-iterator/go"

//
type DataCursor struct {
	// 异步处理数据获取器
	BlockSize int64
	BlockNo   int64
	StartTime string
	EndTime   string
}

func (sel *DataCursor) ToString() string {
	str, _ := jsoniter.MarshalToString(sel)
	return str
}
