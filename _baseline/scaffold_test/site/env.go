package site

import (
	"database/sql"
	"log"
)

// Env is used to pass site-wide configuration and resources to handlers
type Env struct {
	// settings
	HTTPS bool

	// loggers
	ErrorLog *log.Logger
	WarnLog  *log.Logger
	InfoLog  *log.Logger

	// resource pool
	DB *sql.DB

	// EDITME: add your site-wide stuff here
}
