package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func gjr(db *sql.DB, rqParam string) {
	fmt.Println("date==yx=" + rqParam)
	dms := getDm(db, rqParam)
	for _, dm := range dms {
		sps := []float64{}

		dmSql := `SELECT id,date rq,code dm,close sp,high zg,pre_close qsp, m5,open kp  FROM dayline where code = '` +
			dm["code"].(string) + `'  ORDER BY dm,rq DESC limit 3`
		stmt, err := db.Prepare(dmSql)
		if err != nil {
			fmt.Printf("query prepare err:%s\n", err.Error())
			return
		}

		rows, err := stmt.Query()
		defer Closedb(stmt, rows)
		if err != nil {
			fmt.Printf("query err:%s\n", err.Error())
		}
		var lastsp float64
		yx := 0
		cnt := 0
		for rows.Next() {
			var id int
			var rq string
			var dm string
			var sp float64
			var zg float64
			var qsp float64
			var m5 float64
			var kp float64
			rows.Scan(&id, &rq, &dm, &sp, &zg, &qsp, &m5, &kp)
			//fmt.Println(rq, rqParam)
			if rq > rqParam {
				sps = append(sps, sp)
				continue
			}
			if cnt == 0 && sp > m5 {
				//yx = 1
				//fmt.Println(sp < zg*0.99, zg*0.99)
				lastsp = sp
			} else if cnt == 1 && lastsp > zg && zg > sp*1.03 && zg > kp*1.03 {
				yx = 1
				//fmt.Println(dm, yx, lastsp, sp, zg, sp*1.04)
				break
			} else {
				break
			}
			//yx = 1
			cnt++
		}
		if yx > 0 {
			fmt.Println(dm["code"].(string)[0:6], yx)
		}
		Closedb(stmt, rows)
	}
}
