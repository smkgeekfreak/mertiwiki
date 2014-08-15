// PgFunc.go
package pgreflect

import (
	"log"
	"mws/db"
	"mws/util"
	//"os"
	"reflect"
	"strconv"
)

type PgFunc struct {
	Name          string
	NumArgs       int
	ArgMap        map[string]string
	ReturnTypeStr string
	Describe      string
}

//
// Generate PostgreSQL function syntax for selecting all return values from
// using the specified function name and a variable number of arguments
// passed into the function.
// Example Output from a function name "samplefunction" that accepts
// 2 arguments `SELECT * from samplefunction ($1, $2)`
func (pgf *PgFunc) generatePgFuncSQL() string {
	funcName := util.GetCallerName()
	if pgf == nil {
		return ""
	}
	sqlstring := `SELECT * from ` + pgf.Name
	if pgf.NumArgs > 0 {
		sqlstring += `(`
		for i := 0; i < pgf.NumArgs; i++ {
			sqlstring += `$` + strconv.Itoa(i+1)
			if i < pgf.NumArgs-1 {
				sqlstring += ", "
			} else {
				sqlstring += ") "
			}
		}

	}
	log.Printf("%s SQL String for %s = %s", funcName, pgf.Name, sqlstring)

	return sqlstring
}

func (pgf *PgFunc) VariadicScan(variadicArgs ...interface{}) (results []map[string]interface{}, retCode int, err error) {
	funcName := util.GetCallerName()

	if pgf == nil {
		log.Printf("%s PgFunc invalid", funcName)
		return nil, -1, nil
	}
	log.Printf("%s  args %#v ", funcName, variadicArgs)

	db, err := db.GetConnection()
	if err != nil {
		log.Printf("%s Connection err %#v", err)
		return nil, -1, err
	}

	query := pgf.generatePgFuncSQL()
	defer db.Close()
	log.Printf("%s Connected to Postgres, %v preparing %s with %v", funcName, pgf, query, variadicArgs)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("%s Prepare err %#v", err)
		//stmt.Close()
		//db.Close()
		return nil, -1, err
	}
	result, err := stmt.Query(variadicArgs...)
	if err != nil {
		log.Printf("%s Execution problem: %#v", funcName, err)
		//stmt.Close()
		//db.Close()
		return nil, -1, err
	}
	if result == nil {
		log.Printf("%s No results returned: %#v", funcName, result)
		//stmt.Close()
		//db.Close()
		return nil, 1, nil
	}
	columns, err := result.Columns()
	if err != nil {
		log.Printf("%s Colums err %#v", funcName, err)
		//stmt.Close()
		//db.Close()
		return nil, -1, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	rows := make([]map[string]interface{}, 0)
	for result.Next() {
		retVal := make(map[string]interface{})
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := result.Scan(valuePtrs...); err != nil {
			log.Printf("%s Scan Error %#v", funcName, err)
			//result.Close()
			//stmt.Close()
			//db.Close()
			return nil, -1, err
		}
		switch pgf.ReturnTypeStr {
		case "void":
			log.Printf("%s Returns %s type", funcName, "void")
			retVal = nil
		case "integer":
			fallthrough
		case "record":
			log.Printf("%s Returns %s type", funcName, pgf.ReturnTypeStr)
			for i, col := range columns {

				var v interface{}

				val := values[i]

				b, ok := val.([]byte)

				if ok {
					v = string(b)
				} else {
					v = val
				}
				retVal[col] = v
				valType := reflect.TypeOf(v)
				log.Printf("%s returning %v [%v]  of type (%s)", funcName, col, v, valType)
			}
		default:
			log.Printf("%s Returns [%s] type", funcName, pgf.ReturnTypeStr)
			for i, col := range columns {

				var v interface{}

				val := values[i]

				b, ok := val.([]byte)

				if ok {
					v = string(b)
				} else {
					v = val
				}
				retVal[col] = v
				valType := reflect.TypeOf(v)
				log.Printf("%s returning %v [%v]  of type (%s)", funcName, col, v, valType)
			}

		}
		rows = append(rows, retVal)
		log.Printf("%s appened %#v", funcName, rows)
	}

	return rows, 0, nil
}
