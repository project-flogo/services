package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "postgres.connection")
var factory = &PgFactory{}

// var db *sql.DB
type pgConnection struct {
	DatabaseURL          string `md:"databaseUrl"`
	Host                 string `md:"host"`
	Port                 int    `md:"port"`
	User                 string `md:"user"`
	Password             string `md:"password"`
	DbName               string `md:"databaseName"`
	MaxOpenConnections   int    `md:"maxopenconnection"`
	MaxIdleConnections   int    `md:"maxidleconnection"`
	MaxConnLifetime      string `md:"connmaxlifetime"`
	MaxConnRetryAttempts int    `md:"maxconnectattempts"`
	ConnRetryDelay       int    `md:"connectionretrydelay"`
	TLSEnable            bool   `md:"tlsEnable"`
	TLSMode              string `md:"tlsMode"`
	CACert               string `md:"cacert"`
	ClientCert           string `md:"clientcert"`
	ClientKey            string `md:"clientkey"`
}

// PgFactory for postgres connection
type PgFactory struct {
}

// Type PgFactory
func (*PgFactory) Type() string {
	return "Postgres"
}

func decodeTLSParam(tlsparm string) string {
	switch tlsparm {
	case "VerifyCA":
		return "verify-ca"
	case "VerifyFull":
		return "verify-full"
	default:
		return ""
	}
}

func NewDB(settings map[string]interface{}) (*sql.DB, error) {
	var err error

	s := &pgConnection{}
	err = metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}
	cHost := s.Host
	if cHost == "" {
		return nil, errors.New("Required Parameter Host Name is missing")
	}
	cPort := s.Port
	if cPort == 0 {
		return nil, errors.New("Required Parameter Port is missing")
	}
	cDbName := s.DbName
	if cDbName == "" {
		return nil, errors.New("Required Parameter Database name is missing")
	}
	cUser := s.User
	if cUser == "" {
		return nil, errors.New("Required Parameter User is missing")
	}
	cPassword := s.Password
	if cPassword == "" {
		return nil, errors.New("Required Parameter Password is missing")
	}

	cMaxOpenConn := s.MaxOpenConnections
	if cMaxOpenConn == 0 {
		logCache.Debug("Default value of Max open connection chosen. Default is 0")
	}
	if cMaxOpenConn < 0 {
		logCache.Debugf("Max open connection value received is %d, it will be defaulted to 0", cMaxOpenConn)
	}
	cMaxIdleConn := s.MaxIdleConnections
	if cMaxIdleConn == 2 {
		logCache.Debug("Default value of Max idle connection chosen. Default is 2")
	}
	if cMaxIdleConn < 0 {
		logCache.Debugf("Max idle connection value received is %d, it will be defaulted to 0", cMaxIdleConn)
	}
	cMaxConnLifetime := s.MaxConnLifetime
	if cMaxConnLifetime == "0" {
		logCache.Debug("Default value of Max lifetime of connection chosen. Default is 0")
	}
	if strings.HasPrefix(cMaxConnLifetime, "-") {
		logCache.Debugf("Max lifetime connection value received is %s, it will be defaulted to 0", cMaxConnLifetime)
	}
	var lifetimeDuration time.Duration
	if cMaxConnLifetime != "" {
		lifetimeDuration, err = time.ParseDuration(cMaxConnLifetime)
		if err != nil {
			return nil, fmt.Errorf("Could not parse connection lifetime duration")
		}
	}
	cMaxConnRetryAttempts := s.MaxConnRetryAttempts
	if cMaxConnRetryAttempts == 0 {
		logCache.Debug("Maximum connection retry attempt is 0, no retry attempts will be made")
	}
	if cMaxConnRetryAttempts < 0 {
		logCache.Debugf("Max connection retry attempts received is %d", cMaxConnRetryAttempts)
		return nil, fmt.Errorf("Max connection retry attempts cannot be a negative number")
	}
	cConnRetryDelay := s.ConnRetryDelay
	if cConnRetryDelay == 0 {
		logCache.Debug("Connection Retry Delay is 0, no delays between the retry attempts")
	}
	if cConnRetryDelay < 0 {
		logCache.Debugf("Connection retry delay value received is %d", cConnRetryDelay)
		return nil, fmt.Errorf("Connection retry delay cannot be a negative number")
	}
	cConnTimeout := 10 // conn timeout
	logCache.Debugf("Connection timeout value configured is %d", cConnTimeout)

	cTLSConfig := s.TLSEnable
	var conninfo string
	if cTLSConfig == false {
		logCache.Debugf("Login attempting plain connection")
		conninfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d ", cHost, cPort, cUser, cPassword, cDbName, cConnTimeout)
	} else {
		logCache.Debugf("Login attempting SSL connection")
		cTLSMode := s.TLSMode
		conninfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d ",
			cHost, cPort, cUser, cPassword, cDbName, decodeTLSParam(cTLSMode), cConnTimeout)
		//create temp file
		pwd, err := os.Getwd()
		if err != nil {
			logCache.Errorf("could not get working dir due to %s", err.Error())
			return nil, fmt.Errorf("could not get working dir due to %s", err.Error())
		}
		if s.CACert != "" {
			pathCACert := filepath.Join(pwd, "caCert.pem")
			err = os.WriteFile(pathCACert, []byte(s.CACert), 0644)
			if err != nil {
				logCache.Errorf("could not create CA cert file due to %s", err.Error())
				return nil, fmt.Errorf("could not create CA cert file due to %s", err.Error())
			}
			conninfo = conninfo + fmt.Sprintf("sslrootcert=%s ", pathCACert)
		}
		if s.ClientCert != "" {
			pathClientCert := filepath.Join(pwd, "clientCert.pem")
			err = os.WriteFile(pathClientCert, []byte(s.ClientCert), 0644)
			if err != nil {
				logCache.Errorf("could not create client cert file due to %s", err.Error())
				return nil, fmt.Errorf("could not create client cert file due to %s", err.Error())
			}
			conninfo = conninfo + fmt.Sprintf("sslcert=%s ", pathClientCert)
		}
		if s.ClientKey != "" {
			pathClientKey := filepath.Join(pwd, "cacert.pem")
			err = os.WriteFile(pathClientKey, []byte(s.ClientKey), 0644)
			if err != nil {
				logCache.Errorf("could not create client key file due to %s", err.Error())
				return nil, fmt.Errorf("could not create client key file due to %s", err.Error())
			}
			conninfo = conninfo + fmt.Sprintf("sslkey=%s ", pathClientKey)
		}
	}
	// add connection delay
	// check for bad err connecton and then only do retry, dont do it for invalid creds types of errors
	// move the code outside of the for loop for attempt=0 since we are treating 0 as no reattempts
	// do sql.open outside of loop and then check for the errors based on error type and if reattempts!=0 then do retry
	//
	var db *sql.DB
	dbConnected := 0
	if cMaxConnRetryAttempts == 0 {
		logCache.Info("No connection retry selected, connection attempt will be tried only once...")
		db, err = sql.Open("postgres", conninfo)
		if err != nil {
			return nil, fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
		} else {
			err = db.Ping()
			if err != nil {
				return nil, fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
			}
			dbConnected = 1
		}
	} else if cMaxConnRetryAttempts > 0 {
		logCache.Debugf("Maximum connection retry attempts allowed - %d", cMaxConnRetryAttempts)
		logCache.Debugf("Connection retry delay - %d", cConnRetryDelay)
		db, err = sql.Open("postgres", conninfo)
		if err != nil {
			// return nil, dont retry
			return nil, fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
		}
		// sql returned db handle in first attempt
		if db != nil && err == nil {
			logCache.Info("Trying to ping the database server...")
			err = db.Ping()
			// retry attempt on ping only for conn refused and driver bad conn
			if err != nil {
				if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
					strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
					strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
					strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
					logCache.Info("Failed to ping the database server, trying again...")
					for i := 1; i <= cMaxConnRetryAttempts; i++ {
						logCache.Infof("Connecting to database server... Attempt-[%d]", i)
						// retry delay
						time.Sleep(time.Duration(cConnRetryDelay) * time.Second)
						err = db.Ping()
						if err != nil {
							if err == driver.ErrBadConn || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "network is unreachable") ||
								strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp: lookup") ||
								strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "timedout") ||
								strings.Contains(err.Error(), "timed out") || strings.Contains(err.Error(), "net.Error") || strings.Contains(err.Error(), "i/o timeout") {
								continue
							} else {
								return nil, fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
							}
						} else {
							// ping succesful
							dbConnected = 1
							logCache.Infof("Successfully connected to database server in attempt-[%d]", i)
							break
						}
					}
					if dbConnected == 0 {
						logCache.Errorf("Could not connect to database server even after %d number of attempts", cMaxConnRetryAttempts)
						return nil, fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
					}
				} else {
					return nil, fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
				}
			} else {
				logCache.Info("ping to database server is successful...")
			}
			if dbConnected != 0 {
				logCache.Info("Successfully connected to database server")
			}
		}
	}
	db.SetMaxOpenConns(cMaxOpenConn)
	db.SetMaxIdleConns(cMaxIdleConn)
	db.SetConnMaxLifetime(lifetimeDuration)
	logCache.Debug("----------- DB Stats after setting extra connection configs -----------")
	logCache.Debug("Max Open Connections: " + strconv.Itoa(db.Stats().MaxOpenConnections))
	logCache.Debug("Number of Open Connections: " + strconv.Itoa(db.Stats().OpenConnections))
	logCache.Debug("In Use Connections: " + strconv.Itoa(db.Stats().InUse))
	logCache.Debug("Free Connections: " + strconv.Itoa(db.Stats().Idle))
	logCache.Debug("Max idle connection closed: " + strconv.FormatInt(db.Stats().MaxIdleClosed, 10))
	logCache.Debug("Max Lifetime connection closed: " + strconv.FormatInt(db.Stats().MaxLifetimeClosed, 10))

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// copyCertToTempFile creates temp mssql.pem file for running app in container
// and sqlserver needs filepath for ssl cert so can not pass byte array which we get from connection tile
func copyCertToTempFile(certdata []byte, name string) (string, error) {
	var path = name + "_" + "postgresql.pem"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return "", fmt.Errorf("Could not create file %s %s ", path, err.Error())
		}
		if err := os.Chmod(path, 0600); err != nil {
			return "", fmt.Errorf("Could not give permissions file %s %s ", path, err.Error())
		}
		defer file.Close()
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return "", fmt.Errorf("Could not open file %s %s", path, err.Error())
	}
	_, err = file.Write(certdata)
	if err != nil {
		return "", fmt.Errorf("Could not write data to file %s %s", path, err.Error())
	}
	return path, nil
}

//Get cacert file from connection
func getByteCertDataForPemFile(cacert string) ([]byte, error) {
	//case when you provide file having less permission on Ui but as onUI there is no check as per now for cert and key
	// we need to handle it
	if cacert == "" {
		return nil, nil
	}

	cacertMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(cacert), &cacertMap)
	if err != nil {
		return nil, nil // as cert is not mandatory
	}

	base64ContentStr := ""
	for key, value := range cacertMap {
		if key == "content" {
			if value == nil || value == "" {
				return nil, nil // as cert is not mandatory
			}
			base64ContentStr = value.(string)
			break
		}
	}

	if base64ContentStr != "" {
		index := strings.IndexAny(base64ContentStr, ",")
		if index > -1 {
			base64ContentStr = base64ContentStr[index+1:]
		}

		byteCertData, err := base64.StdEncoding.DecodeString(base64ContentStr)
		if err != nil {
			return nil, fmt.Errorf("Error while getting bytes from cert file %s", err.Error())
		}
		return byteCertData, nil
	}
	return nil, fmt.Errorf("%s", err.Error())
}
