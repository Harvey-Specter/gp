package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Param struct {
	db *sql.DB
	rq string
}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3307)/stock?charset=utf8mb4")
	if err != nil {
		fmt.Printf("open sql err:%s\n", err.Error())
		return
	}
	defer db.Close()
	//cjlBs := 1.5
	start := time.Now().UnixNano()
	//genMx(db, 5) //计算均线
	//fmt.Println("123456")

	rq := "2021-01-01"
	tname := "dayline"
	if len(os.Args) != 3 {
		fmt.Println("命令行参数数量错误,应该是3, 日期,表名 ; 目前长度是:", len(os.Args))
		os.Exit(1)
	}
	for k, v := range os.Args {
		if k == 1 {
			rq = v
		} else if k == 2 {
			tname = v
		}
	}

	dltp3l(db, rq, 20, tname)
	dltp3l(db, rq, 30, tname)
	dltp3l(db, rq, 60, tname)
	// save5DayTP(db, rq)
	end := time.Now().UnixNano()
	fmt.Printf("dltp3l cost is :%d \n", (end-start)/1000)

	xc(db, rq, tname) //效果不好 暂停执行 调一调
	end = time.Now().UnixNano()
	fmt.Printf("xc cost is :%d \n", (end-start)/1000)

	qkht(db, rq, tname)
	stars(db, rq, tname)

	end = time.Now().UnixNano()
	fmt.Printf(" cost is :%d \n", (end-start)/1000)
}
