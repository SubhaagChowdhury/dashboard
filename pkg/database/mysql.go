package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

/*
mysql structure
*/
type Object struct {
	host     string
	port     int
	userName string
	password string
	Database string

	sslFlag               bool
	clientCert            string
	clientKey             string
	caCert                string
	tlsOnlyCaCert         bool
	tlsInsecureSKipVerify bool

	idleConns int
	openConns int
	maxLife   int
	timeout   string

	retryCount int
	delayTime  int

	InterruptFlag bool

	client *sql.DB
}

/*
mysql interface declaration
*/
type Interface interface {
	Reconnect() error
	Check() error
	Disconnect() error
	Connect() error
	connectWithoutSSL() error
	connectWithSSL() error
}

/*
initialization function
*/
func (msql *Object) LoadConfigurations(host, username, database, password, clientCert, clientKey, caCert string, port, idleConns, openConns, maxLife, retryCount, delayTime int, sslFlag, tlsOnlyCaCert, tlsInsecureSKipVerify bool, timeout string) error {
	var err error
	// panic handler
	defer func() {
		if err := recover(); err != nil {
			fmt.Println()
			if e, ok := err.(error); ok {
				fmt.Println(e)
			}
		}
	}()

	// data initialization
	msql.InterruptFlag = false

	msql.host = host
	msql.port = port
	msql.userName = username
	msql.Database = database
	msql.password = password

	if err != nil {
		return fmt.Errorf("decompress error %v", err)
	}

	msql.sslFlag = sslFlag
	msql.clientCert = clientCert
	msql.clientKey = clientKey
	msql.caCert = caCert
	msql.tlsOnlyCaCert = tlsOnlyCaCert
	msql.tlsInsecureSKipVerify = tlsInsecureSKipVerify

	msql.idleConns = idleConns
	msql.openConns = openConns
	msql.maxLife = maxLife
	msql.timeout = timeout

	msql.retryCount = retryCount
	msql.delayTime = delayTime

	return nil
}

/*
connection V2 function
*/
func (msql *Object) Connect() error {
	var err error
	if !msql.sslFlag {
		err = msql.connectWithoutSSL()
	} else {
		err = msql.connectWithSSL()
	}

	if err == nil {
		// check whether mysql is connected
		err := msql.Check()
		if err != nil {
			// mysql connection was not active
			return fmt.Errorf("mysql not connected %v", err)
		}
		// allow idle connections to exist
		msql.client.SetMaxIdleConns(msql.idleConns)
		// mysql connection active
	}

	return err
}

func (msql *Object) connectWithoutSSL() error {
	// create connections
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&timeout=%s&readTimeout=%s&writeTimeout=%s", msql.userName, msql.password, msql.host, msql.port, msql.Database, msql.timeout, msql.timeout, msql.timeout))
	if err != nil {
		// connection could not be activated
		return fmt.Errorf("connect failed. %v", err)
	}

	// connect and store client
	if err := conn.Ping(); err != nil {
		return err
	}
	msql.client = conn

	return nil
}

func (msql *Object) connectWithSSL() error {
	/* This method makes connection to mysql using SSL certificates
	return type: void */
	err := msql.sslCheck()
	if err != nil {
		return err
	}

	// create connections
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&timeout=%s&readTimeout=%s&writeTimeout=%s&tls=custom", msql.userName, msql.password, msql.host, msql.port, msql.Database, msql.timeout, msql.timeout, msql.timeout))
	if err != nil {
		// connection could not be activated
		return fmt.Errorf("ssl connect failed. %v", err)
	}

	// connect and store client
	if err := conn.Ping(); err != nil {
		return err
	}
	msql.client = conn

	return nil
}

func (msql *Object) sslCheck() error {
	rootCertPool := x509.NewCertPool()
	clientCert := make([]tls.Certificate, 0, 1)

	file, err := os.Open(msql.caCert)
	if err != nil { // handling error
		return fmt.Errorf("certificate load failed.%v", err)
	}
	defer file.Close()

	pem, err := io.ReadAll(file)
	if err != nil { // handling error
		return fmt.Errorf("certificate read failed.%v", err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append PEM.. %v", err)
	}

	if !msql.tlsOnlyCaCert {
		var certs tls.Certificate
		certs, err = tls.LoadX509KeyPair(msql.clientCert, msql.clientKey)
		if err != nil { // handling error
			return fmt.Errorf("ssl connect failed.%v", err)
		}

		clientCert = append(clientCert, certs)
	}

	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: msql.tlsInsecureSKipVerify,
	})
	return err
}

/*
disconnect function
*/

func (msql *Object) Disconnect() error {
	// check status
	if msql.Check() == nil {
		// close client if connected
		err := msql.client.Close()
		if err != nil {
			// closure could not be done normally
			return fmt.Errorf("connection closure failed")
		}
	}
	return nil
}

/*
status check function
*/
func (msql *Object) Check() error {
	// ping mysql for connection check
	err := msql.client.Ping()

	if err != nil {
		// connection inactive
		return fmt.Errorf("client inactive")
	}
	return nil
}

/*
mysql reconnect attempt function
*/
func (msql *Object) Reconnect() error {
	// attempt reconnect
	var retry int
	msql.Disconnect()
	// retry until attempts are exhausted or connection is created
	for retry = 1; retry <= msql.retryCount; retry++ {
		// create connection
		if msql.Connect() != nil {
			break
		}

		fmt.Printf("MySQL ReConnect failed. Retry count:%d\n", retry)
		// delay until next attempt
		time.Sleep(time.Duration(msql.delayTime) * time.Second)
	}

	// if reconnection failed multiple times close application
	if retry > msql.retryCount {
		// reconnection was unsuccessful after repeated attempts
		return fmt.Errorf("reconnect failed")
	}

	return nil
}

func (msql *Object) PrepareStatement(query string) (*sql.Stmt, error) {
	statement, err := msql.client.Prepare(query)
	if err != nil {
		return statement, err
	}
	return statement, nil
}

func (msql *Object) PrepareAndSelectRowQuery(query string, args []any) (*sql.Row, error) {
	statement, err := msql.client.Prepare(query)
	if err != nil {
		return &sql.Row{}, err
	}
	result := statement.QueryRow(args)
	return result, nil
}
func (msql *Object) PrepareAndSelectQuery(query string, args []any) (*sql.Rows, error) {
	statement, err := msql.client.Prepare(query)
	if err != nil {
		return &sql.Rows{}, err
	}
	result, err := statement.Query(args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (msql *Object) SelectRowQuery(statement *sql.Stmt, args []any) ([]any, error) {
	result := statement.QueryRow(args)
	var data []any
	err := result.Scan(&data)
	if err == sql.ErrNoRows {
		return data, err
	} else if err != nil {
		return data, err
	}
	return data, nil
}

func (msql *Object) SelectQuery(statement *sql.Stmt, args []any) ([][]any, error) {
	result, err := statement.Query(args...)
	if err != nil {
		return [][]any{}, err
	}
	defer result.Close()

	values := make([][]any, 0)

	cols, _ := result.Columns()
	for result.Next() {
		var vals = make([]interface{}, len(cols))

		// Create a slice to hold the scanned values for each row
		var rowVal = make([]interface{}, len(cols))
		for i := range vals {
			rowVal[i] = new(interface{})
		}

		err = result.Scan(rowVal...)
		if err == sql.ErrNoRows {
			return values, nil
		} else if err != nil {
			return values, err
		}

		for i, v := range rowVal {
			switch vv := (*v.(*interface{})).(type) {
			case int64:
				rowVal[i] = vv
			case float64:
				rowVal[i] = vv
			case string:
				rowVal[i] = vv
			case []byte:
				rowVal[i] = string(vv)
			case nil:
				rowVal[i] = nil
			default:
				return values, fmt.Errorf("unsupported data type in scan")
			}
		}

		values = append(values, rowVal)
	}
	return values, nil
}

func (msql *Object) InsertORUpdate(statement *sql.Stmt, args []any) (int, int64, error) {
	tx, err := msql.client.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, 0, fmt.Errorf("insert/update tx begin failed %v", err)
	}
	defer tx.Rollback()

	result, err := tx.StmtContext(context.Background(), statement).ExecContext(context.Background(), args)

	if err != nil {
		return 0, 0, fmt.Errorf("insert/update execution failed %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, 0, fmt.Errorf("insert/update rows affected failed %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, fmt.Errorf("insert/update failed to commit %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return int(rows), id, err
	}

	return int(rows), id, nil
}
