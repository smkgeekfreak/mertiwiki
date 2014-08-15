// ReflectivePostgFunc.go
package pgreflect

import (
	//"errors"
	//"fmt"
	"log"
	"mws/db"
	"mws/util"
	"strings"
)

func GetPgFuncInfo(name string) (pgfunc *PgFunc, reterr error) {
	funcName := util.GetCallerName()
	// Convert name to all lower case for proper lookup
	name = strings.ToLower(strings.Replace(name, " ", "", -1))
	log.Printf("%s Called with %s", funcName, name)

	conn, err := db.GetConnection()
	if err != nil {
		log.Printf("%s Connection err %#v", funcName, err)
		return nil, err
	}
	log.Printf("%s Connected to Postgres", funcName)
	defer conn.Close()
	stmt, err := conn.Prepare(
		`SELECT  proname, proargnames, pronargs,prosrc,format_type(prorettype, NULL),pg_get_function_identity_arguments('` + name + `'::regproc)
			FROM    pg_catalog.pg_namespace n
			JOIN    pg_catalog.pg_proc p
			ON      pronamespace = n.oid
			WHERE   nspname = 'public' AND p.proname = '` + name + `';`)

	if err != nil {
		log.Printf("Prepare err %#v", err)
		//stmt.Close()
		//conn.Close()
		return nil, err
	}

	result, err := stmt.Query()

	if err != nil {
		log.Printf("%s Execution problem: %#v", funcName, result)
		//result.Close()
		//stmt.Close()
		//conn.Close()
		return nil, err
	}
	for result.Next() {
		var name, args, src, rettype, argtypes string
		var numargs int

		if err := result.Scan(&name, &args, &numargs, &src, &rettype, &argtypes); err != nil {
			log.Printf("%s Results Error %v", funcName, err)
			//result.Close()
			//stmt.Close()
			//conn.Close()
			return nil, err
		}

		//log.Printf("%s  \r\n\tPg Function Name - %s,  \r\n\tNum Args - %d, \r\n\tArgs - %s, \r\n\tArgName,Types - %s, \r\n\tRetType - %s, \r\n\tSrc - \t%s,",
		//	funcName, name, numargs, args, argtypes, rettype, src)
		pgfunc := PgFunc{
			Name:          name,
			NumArgs:       numargs,
			ArgMap:        parsePgFuncArgTypes(numargs, argtypes),
			ReturnTypeStr: rettype,
			Describe:      src,
		}

		log.Printf("%s  \r\n\tPg Func - \r\n\t%v", funcName, pgfunc)
		return &pgfunc, nil
	}

	//result.Close()
	//stmt.Close()
	//conn.Close()
	return nil, nil

}

func parsePgFuncArgTypes(numArgs int, argNamesTypes string) map[string]string {
	funcName := util.GetCallerName()
	var argmap map[string]string
	argmap = make(map[string]string)

	var argName, argType string
	argNamesTypes = strings.TrimSpace(argNamesTypes)
	//log.Printf("arrNTypes %s", argNamesTypes)
	arr := strings.Split(argNamesTypes, ",")
	//log.Printf("arr %s", arr)
	for i, val := range arr {
		if i < numArgs {
			//log.Printf("arr[%d ]%s", i, val)
			val = strings.TrimSpace(val)
			arr2 := strings.Split(val, " ")
			for j, val2 := range arr2 {
				//log.Printf("arr2[ %d ]%s", j, val2)
				if j == 0 {
					//log.Printf("name[ %d ]%s", j, val2)
					argName = val2
				} else {
					val2 = strings.Join(append(arr2[j:]), " ")
					//log.Printf("type[ %d ]%s", j, val2)
					argType = val2
					argmap[argName] = argType
					//log.Printf("Argmap = %v , %s", argmap, argType)
					break
				}
			}
		}
	}
	log.Printf("%s Full Argmap = %v , %s", funcName, argmap, argType)
	// Check number of args parsed was number expected
	if len(argmap) != numArgs {
		log.Printf("%s Unacceptable number (%d) of arguments parsed, ", funcName, len(argmap))
		return nil //, errors.New(fmt.Sprintf("%s Unacceptable number (%d) of arguments parsed", funcName, len(argmap)))
	}
	return argmap //, nil
}

//Temp method
func GetPostgresFunctions() error {
	funcName := util.GetCallerName()
	//log.Printf("Calling %s (%d, %d)", funcName, tagId, revisionId)

	db, err := db.GetConnection()
	if err != nil {
		log.Printf("%s Connection err %#v", funcName, err)
		return err
	}
	log.Printf("%s Connected to Postgres", funcName)
	defer db.Close()
	stmt, err := db.Prepare(
		`SELECT  proname, proargnames, pronargs,prosrc
		FROM    pg_catalog.pg_namespace n
		JOIN    pg_catalog.pg_proc p
		ON      pronamespace = n.oid
		WHERE   nspname = 'public'`)

	if err != nil {
		log.Printf("Prepare err %#v", err)
		//stmt.Close()
		//db.Close()
		return err
	}

	result, err := stmt.Query()

	if err != nil {
		log.Printf("Execution problem: %#v", result)
		//result.Close()
		//stmt.Close()
		//db.Close()
		return err
	}
	for result.Next() {
		var name, args, src string
		var numargs int

		if err := result.Scan(&name, &args, &numargs, &src); err != nil {
			log.Printf("Results Error %v", err)
			//result.Close()
			//stmt.Close()
			//db.Close()
			return err
		}
		//		log.Printf("Processing return values for entity relationship ids %d, %d", tagId, revisionId)

		log.Printf("%s Function Name - %s, %s, %d, %s ", funcName, name, args, numargs, src)

		stmt2, err := db.Prepare(
			`SELECT pg_get_function_identity_arguments('` + name + `'::regproc)`)
		if err != nil {
			log.Printf("Prepare err %#v", err)
			//stmt.Close()
			//db.Close()
			return err
		}
		result2, err := stmt2.Query()
		if err != nil {
			log.Printf("Execution problem: %#v", result)
			//result.Close()
			//stmt.Close()
			//db.Close()
			return err
		}
		for result2.Next() {
			var argnameNtype string

			if err := result2.Scan(&argnameNtype); err != nil {
				log.Printf("Result2 Error %v", err)
				//result2.Close()
				//stmt2.Close()
				//db.Close()
				return err
			}
			log.Printf("%s  Function Args - %s %s", funcName, name, argnameNtype)
		}

	}

	//result.Close()
	//stmt.Close()
	//db.Close()
	return nil
}
