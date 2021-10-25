package postgres

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//DBDiagnostic ..
type DBDiagnostic struct {
	State       string `json:"State"`
	NativeError string `json:"NativeError"`
	Messge      string `json:"Message"`
}

//ActionError ..
type ActionError struct {
	APIName string         `json:"APIName"`
	Diag    []DBDiagnostic `json:"Diag,omitempty"`
}

//New ..
func New(actionerr ActionError) error {
	return &ActionError{
		APIName: actionerr.APIName,
		Diag:    actionerr.Diag,
	}
}

func (acterr *ActionError) Error() string {
	out, err := json.Marshal(acterr)
	if err != nil {
		return "Error JSON Marshal Fialed"
	}
	return string(out)
}

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_."
const terminators = " ;),<>+-*%/"
const starters = " =,(<>+-*%/"
const quotes = "\"'`"

func alphaOnly(s string, keyset string) bool {
	for _, char := range s {
		if !strings.Contains(keyset, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

//QueryBasicCheck ..
func QueryBasicCheck(s string) error {
	if alphaOnly(string(s[0]), quotes) || alphaOnly(string(s[len(s)-1]), quotes) {
		var dignostic []DBDiagnostic
		QueryError := ActionError{
			APIName: "HandleQuery",
			Diag: append(dignostic, DBDiagnostic{
				State:       "42000",
				NativeError: "1064",
				Messge:      "Query Syntax Error: Query should not be quoted, syntax error at or near " + string(s[0]) + " index:0",
			}),
		}
		return New(QueryError)
	}
	return nil
}

// EvaluateQuery ... consider query:insert into users(?name, ?id) values(?name, ?id), inputData.values=[('john', 'id')]
func EvaluateQuery(query string, inputData Input) (string, []interface{}, []string, error) {
	re := regexp.MustCompile("\\n")
	query = re.ReplaceAllString(query, " ")
	re = regexp.MustCompile("\\t")
	query = re.ReplaceAllString(query, " ")
	query = strings.TrimSpace(query)
	if query[len(query)-1:] != ";" {
		query += ";"
	}
	odbcquery, param := "", ""
	endIndex := 0
	dq, sq, bt := "\"", "'", "`"
	dqm, sqm, btm := false, false, false
	var paramsarray []string
	paramMarker := false
	paramCounter := 0
	err := QueryBasicCheck(string(query))
	if err != nil {
		return "", nil, nil, err
	}
	paramIndex := 1
	prepared := query
	reducedLength := 0
	inputArgs := []interface{}{}
	for i, val := range query {
		chr := string(val)
		if !dqm || !sqm || !btm {
			if chr == dq && string(query[i-1]) != "\\" && !sqm && !btm {
				dqm = !dqm
				continue
			}
			if chr == sq && string(query[i-1]) != "\\" && !dqm && !btm {
				sqm = !sqm
				continue
			}
			if chr == bt && string(query[i-1]) != "\\" && !dqm && !sqm {
				btm = !btm
				continue
			}
			if !dqm && !sqm && !btm {
				if chr == "?" && alphaOnly(string(query[i-1]), starters) {
					if paramMarker {
						paramMarker = false
						param = ""
						continue
					}
					paramMarker = true
					param = ""
					continue
				}
				if paramMarker {
					if !alphaOnly(chr, alpha) && chr != "\n" {
						paramMarker = false
						if param == "" && string(query[i-1]) == "?" && alphaOnly(chr, terminators) {
							return "", nil, nil, fmt.Errorf("Parameters can not be unnamed, hint: ?paramname")
						}
						if param == "" && string(query[i-1]) == "?" && !alphaOnly(chr, terminators) {
							continue
						}
						if alphaOnly(chr, terminators) {
							paramsarray = append(paramsarray, param)
							paramLength := len(param)
							substituteStartIndex := i - paramLength - 1 - reducedLength
							substitute := "$" + strconv.Itoa(paramIndex)
							//substitution of paramIndex-> insert into users(?name, ?id) values($1, ?id)
							prepared = prepared[:substituteStartIndex] + substitute + prepared[i-reducedLength:]
							paramIndex++
							reducedLength += paramLength - len(substitute) + 1

							substitution, error := getSubstutionValue(param, inputData)
							if error != nil {
								return "", nil, nil, error
							}
							// contains actual values :'john', 'id' mapped with $index in prepared query
							inputArgs = append(inputArgs, substitution)

							odbcquery += query[endIndex : i-len(param)]
							endIndex = i
							paramCounter++
							param = ""
							continue
						}
					}
					param = param + chr
				}
			}
		}
	}
	odbcquery += query[endIndex:]
	return prepared, inputArgs, paramsarray, nil
}

//ex- will extract 'john' from inputdata.Values or inputdata.Parameters for parameter 'name'
func getSubstutionValue(parameter string, inputData Input) (interface{}, error) {
	substitution, ok := inputData.Parameters[parameter]
	if !ok {
		for _, values := range inputData.Values {
			substitution, ok = values[parameter]
			if ok {
				break
			}
		}
		if !ok {
			return nil, fmt.Errorf("missing substitution for: %s", parameter)
		}
	}
	// parameterType, ok := schema[parameter]
	// if ok && parameterType == "BYTEA" {
	// 	substitution = decodeBlob(substitution.(string))
	// }

	return substitution, nil
}
