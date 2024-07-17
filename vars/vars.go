package vars

import "os"

var VERSION string

var HELP_MESSAGE string

var ENV string

var DB_FILE string

func init() {
	ENV = os.Getenv("SIZR_ENV")
}
