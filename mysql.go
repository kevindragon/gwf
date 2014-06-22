// mysql操作。mysql数据模型

package gwf

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbHost = "localhost"
	dbPort = 3306
	dbUser = "root"
	dbPwd  = "mysql5"
	dbName = "myblog"
)

var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
	dbUser, dbPwd, dbHost, dbPort, dbName)

type MysqlDB struct {
	db *sql.DB
}

var mysqlDB *MysqlDB

func init() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("Conn't open mysql database connect")
	}
	mysqlDB = &MysqlDB{db}
}

// 获取新的MySQL操作
func NewMysqlDB() *MysqlDB {
	return mysqlDB
}

func fieldsAandScans(reciever interface{}) (string, []string, []interface{}) {
	elem := reflect.ValueOf(reciever).Elem()
	rtype := elem.Type()

	fmt.Println(rtype.Kind())

	if rtype.Kind() != reflect.Struct {
		panic("数据接收者错误")
	}

	tableName := rtype.Name()

	// 字段数量
	elemNum := elem.NumField()

	fields := make([]string, 0)
	for i := 0; i < elemNum; i++ {
		fields = append(fields, strings.ToLower(rtype.Field(i).Name))
	}

	scans := make([]interface{}, elemNum)
	for i := range scans {
		switch elem.Field(i).Kind() {
		case reflect.String:
			var inf string
			scans[i] = &inf
		case reflect.Int:
			var inf int
			scans[i] = &inf
		case reflect.Int16:
			var inf int16
			scans[i] = &inf
		case reflect.Int32:
			var inf int32
			scans[i] = &inf
		case reflect.Int64:
			var inf int64
			scans[i] = &inf
		default:
			var inf interface{}
			scans[i] = &inf
		}
	}

	return tableName, fields, scans
}

// 通过id获取一条数据
// 使用reciever名字作为数据库
func (d *MysqlDB) GetById(id int, reciever interface{}) {
	tableName, fields, scans := fieldsAandScans(reciever)

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?",
		strings.Join(fields, ", "), strings.ToLower(tableName))

	row := d.db.QueryRow(sqlStr, id)
	row.Scan(scans...)

	elem := reflect.ValueOf(reciever).Elem()
	for i, v := range scans {
		column := reflect.ValueOf(v).Elem().Interface()
		e := elem.Field(i)
		if e.CanSet() {
			e.Set(reflect.ValueOf(column))
		}
	}
}

func (d *MysqlDB) GetTop(num int, reciever interface{}) {
	v := reflect.Indirect(reflect.ValueOf(reciever))
	ev := v.Type().Elem()

	tableName := ev.Name()

	fields := make([]string, 0)
	for i := 0; i < ev.NumField(); i++ {
		fields = append(fields, strings.ToLower(ev.Field(i).Name))
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM %s LIMIT %d",
		strings.Join(fields, ", "), strings.ToLower(tableName), num)
	stmt, _ := d.db.Prepare(sqlStr)
	rows, _ := stmt.Query()

	for rows.Next() {
		// 新建一个结构体指针
		ns := reflect.New(ev).Elem()

		scans := make([]interface{}, 0)
		for i := 0; i < ns.NumField(); i++ {
			f := ns.Field(i)
			scans = append(scans, reflect.New(f.Type()).Interface())
		}

		rows.Scan(scans...)

		for i := 0; i < ns.NumField(); i++ {
			f := ns.Field(i)
			f.Set(reflect.ValueOf(scans[i]).Elem())
		}
		v.Set(reflect.Append(v, ns))
	}

}
