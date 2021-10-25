package postgres

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"

	"database/sql"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"
)

type (
	//Connection datastructure for storing PostgreSQL connection details
	Connection struct {
		DatabaseURL string `json:"databaseURL"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		DbName      string `json:"databaseName"`
		SSLMode     string `json:"sslmode"`
		db          *sql.DB
	}

	//Connector is a representation of connector.json metadata for the postgres connection
	Connector struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Title       string `json:"title"`
		Version     string `json:"version"`
		Type        string `json:"type"`
		Ref         string `json:"ref"`
		Settings    []struct {
			Name  string      `json:"name"`
			Type  string      `json:"type"`
			Value interface{} `json:"value"`
		} `json:"settings"`
	}

	//Query structure for SQL Queries
	Query struct {
		TableName string            `json:"tableName"`
		Cols      []string          `json:"columns"`
		Filters   map[string]string `json:"filters"`
	}

	//record is the one row of the ResultSet retrieved from database after execution of SQL Query
	Recd map[string]interface{}

	//ResultSet is an aggregation of SQL Query data records fetched from the database
	ResultSet struct {
		Record []*Recd `json:"records"`
	}

	//Input is a representation of acitivity's input parametres
	Input struct {
		Parameters map[string]interface{}   `json:"parameters,omitempty"`
		Values     []map[string]interface{} `json:"values,omitempty"`
	}

	//QueryActivity provides Activity metadata for Flogo
	QueryActivity struct {
		metadata *activity.Metadata
	}
)

var connectorCache map[string]*Connection
var connectorCacheMutex sync.Mutex

var preparedQueryCache map[string]*sql.Stmt
var preparedQueryCacheMutex sync.Mutex

func init() {
	connectorCache = make(map[string]*Connection, 100)
	preparedQueryCache = make(map[string]*sql.Stmt, 100)
	go keepalive()
}

type PgSharedConfigManager struct {
	db *sql.DB
}

func decodeBlob(blob string) []byte {
	decodedBlob, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return []byte(blob)
	}
	return decodedBlob
}

func (p *PgSharedConfigManager) getSchema(fields interface{}) map[string]string {
	schema := map[string]string{}
	for _, fieldObject := range fields.([]interface{}) {
		if fieldName, ok := fieldObject.(map[string]interface{})["FieldName"]; ok {
			schema[fieldName.(string)] = fieldObject.(map[string]interface{})["Type"].(string)
		}
	}
	return schema
}

func keepalive() {
	tick := time.Tick(20 * time.Minute)
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-tick:
			for key, conn := range connectorCache {
				err := conn.db.Ping()
				if err != nil {
					log.RootLogger().Warnf("sharedconnection::keepalive Ping connection [%s] is dead: [%s]", key, err)
					conn.db = nil
				} else {
					log.RootLogger().Debugf("sharedconnection::keepalive Ping connection: [%s]", key)
				}
			}
		}
	}
}

//GetStatement
func (p *PgSharedConfigManager) getStatement(prepared string) (stmt *sql.Stmt, err error) {
	preparedQueryCacheMutex.Lock()
	defer preparedQueryCacheMutex.Unlock()
	stmt, ok := preparedQueryCache[prepared]
	if !ok {
		stmt, err = p.db.Prepare(prepared)
		if err != nil {
			return nil, err
		}
		preparedQueryCache[prepared] = stmt
	}
	return stmt, nil
}

//PreparedDelete allows to delete records
func (p *PgSharedConfigManager) PreparedDelete(queryString string, inputData *Input, logCache log.Logger) (results map[string]interface{}, err error) {
	// log.Debugf("Executing prepared query %s", queryString)

	prepared, inputArgs, paramsArray, err := EvaluateQuery(queryString, *inputData)
	if err != nil {
		return nil, err
	}
	logCache.Debugf("Prepared delete statement: [%s]  and  Parameters: [%v], Parameter Values : [%v], Parameters : [%v] ", prepared, inputArgs, paramsArray)

	stmt, err := p.getStatement(prepared)
	if err != nil {
		logCache.Errorf("Failed to prepare statement: %s", err)
		return nil, err
	}
	logCache.Debug("----------- DB Stats in Delete activity -----------")
	p.printDBStats(p.db, logCache)
	// defer stmt.Close()
	result, err := stmt.Exec(inputArgs...)
	if err != nil {
		logCache.Errorf("Executing prepared query got error: %s", err)
		// stmt.Close()
		return nil, err
	}
	// sharedconnection.returnStatement(prepared, stmt)
	output := make(map[string]interface{})
	output["rowsAffected"], _ = result.RowsAffected()
	return output, nil
}

//PreparedInsert allows querying database with named parameters
func (p *PgSharedConfigManager) PreparedInsert(queryString string, inputData *Input, fields interface{}, logCache log.Logger) (results *ResultSet, err error) {

	schema := p.getSchema(fields)
	prepared, inputArgs, paramsArray, err := EvaluateQuery(queryString, *inputData)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(paramsArray); i++ {
		parameterType, ok := schema[paramsArray[i]]
		if ok && parameterType == "BYTEA" {
			inputArgs[i] = decodeBlob(inputArgs[i].(string))
		}
	}

	// log.Debug("prepared insert [%s]", prepared)
	stmt, err := p.getStatement(prepared)
	// stmt, err := p.db.Prepare(prepared)
	if err != nil {
		logCache.Errorf("Failed to prepare statement: %s", err)
		return nil, err
	}

	logCache.Debug("----------- DB Stats in Insert activity -----------")
	p.printDBStats(p.db, logCache)
	rows, err := stmt.Query(inputArgs...)
	if err != nil {
		logCache.Errorf("Executing prepared query got error: %s", err)
		// stmt.Close()
		return nil, err
	}
	if rows == nil {
		// logCache.Debugf("no rows returned for query %s", prepared)
		return nil, nil
	}
	// p.returnStatement(prepared, stmt)
	defer rows.Close()
	// logCache.Debug("return from prepared insert")
	return UnmarshalRows(rows)
}

// UnmarshalRows function
func UnmarshalRows(rows *sql.Rows) (results *ResultSet, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column information, %s", err.Error())
	}

	count := len(columns)
	cols := make([]interface{}, count)
	args := make([]interface{}, count)
	coltypes, err := rows.ColumnTypes()

	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, fmt.Errorf("error determining column types, %s", err.Error())
	}

	for i := range cols {
		args[i] = &cols[i]
	}

	var resultSet ResultSet
	columnCount := make(map[string]int)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			fmt.Printf("%s", err.Error())
			return nil, fmt.Errorf("error scanning rows, %s", err.Error())
		}
		m := make(Recd)
		for i, b := range cols {
			dbType := coltypes[i].DatabaseTypeName()
			//added check for duplicate columns in case of JOIN, WIPGRS-452
			if _, found := m[columns[i]]; found {
				columnCount[columns[i]]++
				columns[i] = fmt.Sprintf("%s_%d", columns[i], columnCount[columns[i]])
			}
			if b == nil {
				m[columns[i]] = nil
				continue
			}
			switch dbType {
			case "NUMERIC":
				x := b.([]uint8)
				if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
					m[columns[i]] = nx
				}
			case "BOX", "LINE", "POINT", "PATH", "POLYGON", "UUID", "BIT", "VARBIT", "INTERVAL", "MONEY":
				x := b.([]byte)
				m[columns[i]] = string(x)
			case "SMALLINT", "MEDIUMINT", "INT", "BIGINT":
				m[columns[i]] = b.(int64)
			case "DOUBLE", "REAL":
				m[columns[i]] = b.(float64)
			case "CHAR", "NCHAR", "BPCHAR", "":
				x := b.([]byte)
				if len, ok := coltypes[i].Length(); ok == true {
					if len > 0 {
						m[columns[i]] = string(x[0:len])
					}
				}
			case "BYTEA":
				m[columns[i]] = base64.StdEncoding.EncodeToString(b.([]byte))
			case "ANYENUM":
				str := fmt.Sprintf("%v", b)
				if g, ok := base64.StdEncoding.DecodeString(str); ok == nil {
					m[columns[i]] = string(g[:])
				}
			default:
				m[columns[i]] = b
			}
		}
		if len(m) > 0 {
			resultSet.Record = append(resultSet.Record, &m)
		}
	}
	return &resultSet, nil
}

//PreparedUpdate allows updating database with named parameters
func (p *PgSharedConfigManager) PreparedUpdate(queryString string, inputData *Input, fields interface{}, logCache log.Logger) (results *ResultSet, err error) {
	// logCache.Debugf("Executing prepared query %s", queryString)
	// logCache.Debugf("Query parameters: %v", inputData)

	schema := p.getSchema(fields)
	prepared, inputArgs, paramsArray, err := EvaluateQuery(queryString, *inputData)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(paramsArray); i++ {
		parameterType, ok := schema[paramsArray[i]]
		if ok && parameterType == "BYTEA" {
			inputArgs[i] = decodeBlob(inputArgs[i].(string))
		}
		logCache.Debug("--InputArgs: ", inputArgs)
	}
	logCache.Debugf("Prepared Update Statement is ", prepared, " Parameters ", inputArgs)
	stmt, err := p.getStatement(prepared)
	//stmt, err := p.db.Prepare(prepared)
	if err != nil {
		logCache.Errorf("Failed to prepare statement: %s", err)
		return nil, err
	}

	logCache.Debug("----------- DB Stats in Update activity -----------")
	p.printDBStats(p.db, logCache)

	// defer stmt.Close()
	rows, err := stmt.Query(inputArgs...)
	if err != nil {
		logCache.Errorf("Executing prepared query got error: %s", err)
		// stmt.Close()
		return nil, err
	}

	if rows == nil {
		logCache.Debugf("no rows returned for query %s", prepared)
		return nil, nil
	}
	// p.returnStatement(prepared, stmt)
	defer rows.Close()
	// logCache.Debug("Return from PreparedQuery")
	return UnmarshalRows(rows)
}

func checkCount(rows *sql.Rows) (count int, err error) {
	logCache.Debugf("Inside check count rows for update")
	// var counter int
	// defer rows.Close()
	for rows.Next() {
		logCache.Debugf("Row found")
		if err := rows.Scan(&count); err != nil {
			//log.Fatal(err)
		}
		logCache.Debugf("Counter: %d", count)
	}
	return count, err
}

//PreparedQuery allows querying database with named parameters
func (p *PgSharedConfigManager) PreparedQuery(queryString string, inputData *Input, logCache log.Logger) (results *ResultSet, err error) {
	// logCache.Debugf("Executing prepared query %s", queryString)

	prepared, inputArgs, paramsArray, err := EvaluateQuery(queryString, *inputData)
	if err != nil {
		return nil, err
	}
	logCache.Debug("Parameters : ", paramsArray, " Parameter Values : ", inputArgs)

	logCache.Debugf("Prepared Query Statement is : %s ", prepared)
	stmt, err := p.getStatement(prepared)
	// stmt, err := p.db.Prepare(prepared)
	if err != nil {
		logCache.Errorf("Failed to prepare statement: %s", err)
		return nil, err
	}

	logCache.Debug("----------- DB Stats in Query activity -----------")
	p.printDBStats(p.db, logCache)

	// defer stmt.Close()
	rows, err := stmt.Query(inputArgs...)
	if err != nil {
		logCache.Errorf("Executing prepared query got error: %s", err)
		// stmt.Close()
		return nil, err
	}

	if rows == nil {
		logCache.Debugf("no rows returned for query %s", prepared)
		return nil, nil
	}
	// p.returnStatement(prepared, stmt)
	defer rows.Close()
	// logCache.Debug("Return from PreparedQuery")
	return UnmarshalRows(rows)
}

// Login connects to the the postgres database cluster using the connection
// details provided in Connection configuration
func (con *Connection) Login() (err error) {
	if con.db != nil {
		return nil
	}

	conninfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		con.Host, con.Port, con.User, con.Password, con.DbName)

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return fmt.Errorf("Could not open connection to database %s, %s", con.DbName, err.Error())
	}
	con.db = db

	err = db.Ping()
	if err != nil {
		return err
	}

	// log.Infof("login successful")
	return nil
}

//Logout the database connection
func (con *Connection) Logout() (err error) {
	if con.db == nil {
		return nil
	}
	err = con.db.Close()
	// log.Infof("Logged out %s to %s", con.User, con.DbName)
	return
}

func (p *PgSharedConfigManager) printDBStats(db *sql.DB, log log.Logger) {
	log.Debug("Max Open Connections: " + strconv.Itoa(db.Stats().MaxOpenConnections))
	log.Debug("Number of Open Connections: " + strconv.Itoa(db.Stats().OpenConnections))
	log.Debug("In Use Connections: " + strconv.Itoa(db.Stats().InUse))
	log.Debug("Free Connections: " + strconv.Itoa(db.Stats().Idle))
	log.Debug("Max idle connection closed: " + strconv.FormatInt(db.Stats().MaxIdleClosed, 10))
	log.Debug("Max Lifetime connection closed: " + strconv.FormatInt(db.Stats().MaxLifetimeClosed, 10))
}
