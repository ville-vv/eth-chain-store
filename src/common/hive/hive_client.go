package hive

//
import (
	"context"
	"fmt"

	"github.com/beltran/gohive"
)

//"github.com/dazheng/gohive"
//"github.com/uxff/gohive"

type HiveClient struct {
}

func New() {
	cfg := gohive.NewConnectConfiguration()
	cfg.Service = "hive"
	cfg.Database = "test"
	ghConn, err := gohive.Connect("localhost", 10000, "NONE", cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	cursor := ghConn.Cursor()
	cursor.Exec(context.Background(), "show databases")
	for cursor.HasMore(context.Background()) {
		if cursor.Err != nil {
			fmt.Println(cursor.Err)
			break
		}
		var data map[string]interface{}
		cursor.FetchOne(context.Background(), data)

		fmt.Println(data)
	}
}
