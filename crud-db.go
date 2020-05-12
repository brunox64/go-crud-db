package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

var server = "localhost"
var port = 1433
var user = "sankhya"
var password = "senhaaqui"
var database = "nomebaseaqui"

func connect() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func main() {
	connect()

	// Read employees
	// result, err := executeRead("SELECT TOP 10 NUFIN, VLRDESDOB, TIMORIGEM FROM TGFFIN WHERE TIMORIGEM IS NOT NULL", nil)
	result, err := executeRead("SELECT TOP 10 NUFIN, VLRDESDOB, TIMORIGEM FROM TGFFIN WHERE TIMORIGEM IS NULL", nil)
	if err != nil {
		log.Fatal("Error reading Employees: ", err.Error())
	}

	// o conceito está abaixo em x onde somente "x é nil 1:" é verdade
	var x interface{} // como só estou declarando ela assume zero value que é também nil
	if x == nil {
		println("x é nil 1:")
	}
	x = new(interface{})
	if x == nil {
		println("x é nil 2:")
	}
	x = "abc"
	if x == nil {
		println("x é nil 3:")
	}

	//var zeroValue = new(interface{})
	// o zeroValue de uma interface é nil
	var timOrigem interface{}
	for _, m := range result {
		if m["TIMORIGEM"] == nil {
			timOrigem = "[NULL]"
		} else {
			// timOrigem = (*(m["TIMORIGEM"])).(string)
			timOrigem = m["TIMORIGEM"]
		}
		fmt.Printf("val: %d %s %s\n", (m["NUFIN"]), (m["VLRDESDOB"]), timOrigem)
	}
}

// executeRead reads all employee records
func executeRead(tsql string, params map[string]interface{}) ([]map[string]interface{}, error) {
	var result = make([]map[string]interface{}, 0)

	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return result, err
	}

	var args = make([]interface{}, 0)
	if params != nil {
		for k, v := range params {
			args = append(args, sql.Named(k, v))
		}
	}

	tstm, err := db.Prepare(tsql)
	if err != nil {
		return result, err
	}

	// Execute query
	rows, err := tstm.QueryContext(ctx, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return result, err
	}

	// Iterate through the result set.
	for rows.Next() {
		var row = make([]interface{}, len(cols))

		// exemplo com reflection
		// var rowRV = make([]reflect.Value, len(row))
		// for i := range row {
		// 	rowRV[i] = reflect.ValueOf(&row[i])
		// }
		// var res = reflect.ValueOf(rows).MethodByName("Scan").Call(rowRV)[0]
		// if res.Interface() != nil {
		// 	return result, res.Interface().(error)
		// }
		// fim exemplo com reflection

		// exemplo sem reflection
		for i := range row {
			row[i] = &row[i]
		}
		err := rows.Scan(row...)
		if err != nil {
			return result, err
		}
		// fim exemplo sem reflection

		// daqui pra frente nada muda

		// exemplo retornando map[string]*interface{}
		// var mapRow = make(map[string]*interface{})

		// for i, v := range cols {
		// 	mapRow[v] = &row[i]
		// }

		// exemplo retornando map de values
		var mapRow = make(map[string]interface{})

		for i, v := range cols {
			mapRow[v] = row[i]
		}

		result = append(result, mapRow)
	}

	return result, err
}
