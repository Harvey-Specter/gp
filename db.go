package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

// `day` date NOT NULL,
// `code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
// `pattern` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
// `price` double NOT NULL DEFAULT 0,
// `market` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',

func batchInsertPattern(db *sql.DB, data []map[string]string) {
	sqlStr := "replace INTO pattern_day(day, code,pattern,price,market) VALUES "
	vals := []interface{}{}

	for _, row := range data {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, row["day"], row["code"], row["pattern"], row["price"], row["market"])
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)

	//format all vals at once
	res, err := stmt.Exec(vals...)
	if err != nil {
		fmt.Printf("batchInsertPattern err:%s\n", err.Error())
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

// | id          | int unsigned | NO   | PRI | NULL    | auto_increment |
// | price_id    | int unsigned | NO   |     | NULL    |                |
// | day         | date         | NO   | MUL | NULL    |                |
// | code        | varchar(255) | NO   | MUL | NULL    |                |
// | user_id     | int unsigned | NO   |     | NULL    |                |
// | category_id | int unsigned | NO   |     | NULL    |                |
// | pattern     | varchar(255) | NO   |     | NULL    |                |
// | market      | varchar(255) | NO   |     | NULL    |                |
// | remark      | varchar(255) | YES  |     | NULL    |                |
// | created_at  | timestamp    | YES  |     | NULL    |                |
// | updated_at  | timestamp    | YES  |     | NULL    |                |
func batchSaveStockPG(db *sqlx.DB, data []map[string]string, cateId int64) {
	fmt.Println("len(data)===", len(data))
	sqlStr := "insert INTO stocks (price_id, day,code,user_id,category_id,pattern,market,remark,created_at) VALUES "
	//vals := []interface{}{}
	codes := ""

	// tx := db.MustBegin()
	// tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	// tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	// tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	// tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	// // Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	// tx.Commit()

	for _, row := range data {
		if !strings.Contains(codes, row["code"]) {
			// sqlStr += "(?, ?, ?, ?, ? , ?, ?, ?, ?),"

			sqlStr += fmt.Sprintf("(%s,'%s','%s',%s,%d,%s,%s,'%s','%s'),", row["price_id"], row["day"], row["code"], row["user_id"], cateId, row["pattern"], row["market"], row["remark"], row["created_at"])

			//vals = append(vals, row["price_id"], row["day"], row["code"], row["user_id"], cateId, row["pattern"], row["market"], row["remark"], row["created_at"])

			codes += row["code"] + ","
		}
	}
	//trim the last ,
	sqlStr = sqlStr[:len(sqlStr)-1] // remove trailing ","

	// fmt.Println("sqlStr===" + sqlStr)
	db.MustExec(sqlStr)
}
func SaveCategoyStockPG(db *sqlx.DB, name string, code string, remark string, stocklist []map[string]string) int64 {

	seSql := "SELECT  id from categories where name=$1 and code=$2"
	rows, seErr := db.Query(seSql, name, code)
	if seErr != nil {
		fmt.Printf("SaveCategoyStockPG err:%s\n", seErr.Error())
	}
	var id int64
	id = -1
	for rows.Next() {
		seErr = rows.Scan(&id)
	}
	if id == -1 {
		now := time.Now().Format("2006-01-02 15:04:05")

		err := db.QueryRow("INSERT INTO categories (name,code,remark,user_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id", name, code, remark, 1, now, now).Scan(&id)

		// err := db.QueryRow(`INSERT INTO users(name, favorite_fruit, age) VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)

		if err != nil {
			fmt.Printf("insert categories err:%s\n", err.Error())
			return -1
		}
		//id, _ = res.LastInsertId()
		// affect, _ := res.RowsAffected()

		// fmt.Printf("SaveCategoy lastId:%d affectRow:%d\n", id)
		batchSaveStockPG(db, stocklist, id)
	} else {
		batchSaveStockPG(db, stocklist, id)
	}
	return id
}

//-------------------mysql------------------
func batchSaveStock(db *sql.DB, data []map[string]string, cateId int64) {
	fmt.Println("len(data)===", len(data))
	sqlStr := "replace INTO stocks (price_id, day,code,user_id,category_id,pattern,market,remark,created_at) VALUES "
	vals := []interface{}{}

	for _, row := range data {
		sqlStr += "(?, ?, ?, ?, ? , ?, ?, ?, ?),"
		vals = append(vals, row["price_id"], row["day"], row["code"], row["user_id"], cateId, row["pattern"], row["market"], row["remark"], row["created_at"])
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("batchSaveStock Prepare err:%s\n", err.Error())
		return
	}
	//format all vals at once
	res, err := stmt.Exec(vals...)
	if err != nil {
		fmt.Printf("batchInsertPattern err:%s\n", err.Error())
		return
	}
	lastId, _ := res.LastInsertId()
	affect, _ := res.RowsAffected()
	fmt.Printf("lastId:%d affectRow:%d\n", lastId, affect)
}
func SaveCategoyStock(db *sql.DB, name string, code string, remark string, stocklist []map[string]string) int64 {

	seSql := "SELECT  id from categories where name=? and code=?"

	seStmt, seErr := db.Prepare(seSql)
	if seErr != nil {
		fmt.Printf("query prepare err:%s\n", seErr.Error())
	}
	rows, seErr := seStmt.Query(name, code)
	var id int64
	id = -1
	for rows.Next() {
		seErr = rows.Scan(&id)
	}
	if id == -1 {
		stmt, err := db.Prepare("INSERT INTO categories (name,code,remark,user_id,created_at,updated_at) VALUES (?,?,?,?,?,?)")
		if err != nil {
			fmt.Printf("insert prepare err:%s\n", err.Error())
			return -1
		}
		defer stmt.Close()
		now := time.Now().Format("2006-01-02 15:04:05")
		res, err := stmt.Exec(name, code, remark, 1, now, now)
		if err != nil {
			fmt.Printf("insert err:%s\n", err.Error())
			return -1
		}
		id, _ = res.LastInsertId()
		affect, _ := res.RowsAffected()
		fmt.Printf("SaveCategoy lastId:%d affectRow:%d\n", id, affect)
		batchSaveStock(db, stocklist, id)
	} else {
		batchSaveStock(db, stocklist, id)
	}
	return id
}

// | name        | varchar(255) | NO   |     | NULL    |                |
// | code        | varchar(255) | NO   |     |         |                |
// | remark      | varchar(255) | YES  |     | NULL    |                |
// | user_id     | int unsigned | NO   |     | 0       |                |
// | stock_count | int unsigned | NO   |     | 0       |                |
// | created_at  | timestamp    | YES  |     | NULL    |                |
// | updated_at  | timestamp    | YES  |     | NULL    |                |

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
