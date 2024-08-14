package vars

import "os"

var VERSION string

var HELP_MESSAGE string

var ENV string

var DB_FILE string

var BASEDIR string

func init() {
	ENV = os.Getenv("SIZR_ENV")
	BASEDIR = os.Getenv("BASEDIR") // TODO: should be passed with ldflags
}
