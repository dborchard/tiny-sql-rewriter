package database

import (
	"bytes"
	"database/sql"
	"io"
)

// Connector database connection basic objects
type Connector struct {
	Addr     string
	User     string
	Pass     string
	Database string
	Charset  string
	Conn     *sql.DB
}

// QueryResult database query return value
type QueryResult struct {
	Rows      *sql.Rows
	Error     error
	Warning   *sql.Rows
	QueryCost float64
}

// NewConnector 创建新连接
func NewConnector(dsn string) (*Connector, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	connector := &Connector{
		Conn: conn,
	}
	return connector, err
}

// Query execute sql
func (db *Connector) Query(sql string, params ...interface{}) (QueryResult, error) {
	var res QueryResult
	return res, nil
}

// stringEscape mysql_escape_string
// https://github.com/liule/golang_escape
func stringEscape(source string) string {
	var j int
	if source == "" {
		return source
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i, b := range tempStr {
		flag := false
		var escape byte
		switch b {
		case '\000':
			flag = true
			escape = '\000'
		case '\r':
			flag = true
			escape = '\r'
		case '\n':
			flag = true
			escape = '\n'
		case '\\':
			flag = true
			escape = '\\'
		case '\'':
			flag = true
			escape = '\''
		case '"':
			flag = true
			escape = '"'
		case '\032':
			flag = true
			escape = 'Z'
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

// quoteEscape sql_mode=no_backslash_escapes
func quoteEscape(source string) string {
	var buf bytes.Buffer
	last := 0
	for ii, bb := range source {
		if bb == '\'' {
			_, err := io.WriteString(&buf, source[last:ii])
			if err != nil {
				panic(err)
			}
			_, err = io.WriteString(&buf, `''`)
			if err != nil {
				panic(err)
			}
			last = ii + 1
		}
	}
	_, err := io.WriteString(&buf, source[last:])
	if err != nil {
		panic(err)
	}
	return buf.String()
}
