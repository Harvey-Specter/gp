package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func Closedb(stmt *sql.Stmt, rows *sql.Rows) {
	rows.Close()
	stmt.Close()
}
func del(db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM user where uid=?")
	if err != nil {
		fmt.Printf("delete prepare err:%s\n", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(4)
	if err != nil {
		fmt.Printf("delete err:%s\n", err.Error())
		return
	}
	lastId, _ := res.LastInsertId()
	affect, _ := res.RowsAffected()
	fmt.Printf("lastId:%d affectRow:%d\n", lastId, affect)
}

func batchInsertTp(db *sql.DB, data []map[string]string) {
	sqlStr := "INSERT INTO tp(date, code,price) VALUES "
	vals := []interface{}{}

	for _, row := range data {
		sqlStr += "(?, ?,?),"
		vals = append(vals, row["date"], row["code"], row["price"])
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)

	//format all vals at once
	res, err := stmt.Exec(vals...)
	if err != nil {
		fmt.Printf("batchInsert err:%s\n", err.Error())
		return
	}
	lastId, _ := res.LastInsertId()
	affect, _ := res.RowsAffected()
	fmt.Printf("lastId:%d affectRow:%d\n", lastId, affect)
}

func save5DayTP(db *sql.DB, rqParam string) {

	dms := getTpDmByDate(db, rqParam)

	for _, dm := range dms {
		dmSql := "SELECT  date rq, open kp,close sp FROM dayline where code = '" +
			dm["code"] + "' and date>'" + rqParam + "' ORDER BY date limit 5"

		// fmt.Println("dmSql=", dmSql)
		stmt, err := db.Prepare(dmSql)
		if err != nil {
			fmt.Printf("query prepare err:%s\n", err.Error())
		}
		rows, err := stmt.Query()
		defer Closedb(stmt, rows)
		if err != nil {
			fmt.Printf("query err:%s\n", err.Error())
		}

		var pns []float64
		var kp0 float64
		var endDate string
		cnt := 0
		r1 := ""
		for rows.Next() {
			var rq string
			var kp float64
			var sp float64

			rows.Scan(&rq, &kp, &sp)

			if cnt == 0 {
				kp0 = kp
			}
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (sp-kp0)/kp0*100), 64)

			pns = append(pns, value)
			if cnt == 4 {
				endDate = rq
				r1 += fmt.Sprintf("%.2f", sp) + "," + endDate
			} else {
				r1 += fmt.Sprintf("%.2f", sp) + ","
			}
			// fmt.Println(cnt, sp, r1)
			cnt++
		}
		UpdateTp(db, pns, rqParam, dm["code"], r1)
	}
}
func UpdateTp(db *sql.DB, data []float64, date string, code string, r1 string) {

	sqlStr := "update tp set p1=? , p2=?, p3=?, p4=?, p5=? , r1=? where date=? and code = ? "

	stmt, _ := db.Prepare(sqlStr)
	defer stmt.Close()
	// res, err := stmt.Exec(data[0], data[1], data[2], data[3], data[4], r1, date, code)
	_, err := stmt.Exec(data[0], data[1], data[2], data[3], data[4], r1, date, code)

	if err != nil {
		fmt.Printf("UpdateTp err:%s\n", err.Error())
		return
	}
	// lastId, _ := res.LastInsertId()
	// affect, _ := res.RowsAffected()
	// fmt.Printf("lastId:%d affectRow:%d\n", lastId, affect)
}

func insert(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO user (uid,name,age) VALUES (?,?,?)")
	if err != nil {
		fmt.Printf("insert prepare err:%s\n", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(4, "xlxing", 25)
	if err != nil {
		fmt.Printf("insert err:%s\n", err.Error())
		return
	}
	lastId, _ := res.LastInsertId()
	affect, _ := res.RowsAffected()
	fmt.Printf("lastId:%d affectRow:%d\n", lastId, affect)
}
