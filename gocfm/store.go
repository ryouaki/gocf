package gocfm

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/ryouaki/koa/log"
	"github.com/ryouaki/koa/session"
)

var SessionStore session.MemStore = *session.NewMemStore()
var sqliteDB *sql.DB

func init() {
	sql.Register("gocf", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return nil
		},
	})

	srcDb, err := sql.Open("gocf", "./gocfm.db")
	if err != nil {
		log.Error("DB gocf open failed")
		return
	}

	initTables(srcDb)

	sqliteDB = srcDb
}

func initTables(db *sql.DB) {
	initUser(db)
	initScripts(db)
	initClient(db)
}

func initUser(db *sql.DB) {
	var err error
	md5Str := md5.Sum([]byte("123456"))
	tbUser := fmt.Sprintf("create table if not exists Users(UserName varhcar(20) default '',  Password varchar(32) default '%s' )", md5Str)
	_, err = db.Exec(tbUser)
	if err != nil {
		log.Error("Create Table User failed")
		return
	}

	var hasUser *sql.Rows
	hasUser, err = db.Query("select UserName, Password from Users where UserName='admin'")
	defer hasUser.Close()
	if err != nil {
		log.Error("Select Admin failed")
		return
	}

	if !hasUser.Next() {
		_, err = db.Exec(fmt.Sprintf("insert into Users(UserName, Password) values('admin', '%s')", md5Str))
		if err != nil {
			log.Error("Init Admin failed")
			return
		}
	}
}

func initScripts(db *sql.DB) {
	tbScripts := fmt.Sprintf(`create table if not exists Scripts(
		_id INTEGER PRIMARY KEY ASC,
		name varhcar(20) default '',
		mode varchar(1) default '0',
		path varchar(128) default '',
		module_name varchar(32) default '',
		ver varchar(32) default '',
		active varchar(1) default '0',
		script TEXT default ''
	)`)

	_, err := db.Exec(tbScripts)
	if err != nil {
		log.Error("Create Table Scripts failed")
		return
	}
}

func initClient(db *sql.DB) {
	tbClients := fmt.Sprintf(`create table if not exists Clients(
		_id INTEGER PRIMARY KEY ASC,
		ip varhcar(20) default '',
		last_update varchar(32) default CURRENT_TIMESTAMP,
		cpu varchar(10) default '',
		mem varchar(20) default '',
		status varchar(1) default '0',
		listen varchar(1) default '0'
	)`)

	_, err := db.Exec(tbClients)
	if err != nil {
		log.Error("Create Table Scripts failed")
		return
	}
}

func GetDB() *sql.DB {
	return sqliteDB
}
