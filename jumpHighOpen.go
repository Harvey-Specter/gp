package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func jumpHighOpen(db *sql.DB, rqParam string) {
	fmt.Println("date==jumpHighOpen=" + rqParam)
	dms := getDm(db, rqParam)
	for _, dm := range dms {
		sps := []float64{}

		dmSql := `SELECT id,date rq,code dm,close sp,high zg,pre_close qsp, m5,open kp,low zd  FROM dayline where code = '` +
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
		//var lastsp float64
		var lastkp float64
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
			var zd float64
			rows.Scan(&id, &rq, &dm, &sp, &zg, &qsp, &m5, &kp, &zd)
			//fmt.Println(rq, rqParam)
			//fmt.Println(cnt, dm, yx, lastkp, sp, zg, sp)
			if rq >= rqParam {
				sps = append(sps, sp)
				continue
			}
			if cnt == 0 && sp > m5 && zg-sp < zg-zd*0.25 && zg != sp {
				//lastsp = sp
				lastkp = kp
			} else if cnt == 1 && lastkp > zg && zg-sp < zg-zd*0.25 && sp > m5 {
				yx = 1

				break
			} else {
				yx = 0
				break
			}
			//yx = 1
			cnt++
		}
		if yx > 0 {
			fmt.Println(dm["code"].(string)[0:6], reveSliceF(sps))
		}
		Closedb(stmt, rows)
	}
}
