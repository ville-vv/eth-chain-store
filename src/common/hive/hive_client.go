package hive

//
//import (
//	"fmt"
//
//	"github.com/beltran/gohive"
//	//"github.com/dazheng/gohive"
//)
//
//type HiveClient struct {
//}
//
//func New() {
//
//	gh, err := gohive.Connect("172.16.16.155:10000", gohive.DefaultOptions)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	rs, err := gh.Query("select count(*) from etherum.transaction_records ;")
//	if err != nil {
//		panic(err)
//	}
//
//	var i int
//	for rs.Next() {
//		rs.Scan(&i)
//		fmt.Println(i)
//	}
//
//	gh.Close()
//}
