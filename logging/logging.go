package logging

import (
	"os"

	"github.com/op/go-logging"
)

const logformat = "%{color}%{time:20060102150405} [%{level:.8s}] [%{program}>%{callpath}]%{color:reset} %{message}"

var Log = logging.MustGetLogger("acmetool_dns_hooks")

func Init() {
	// Init Logging with Color, i love colors :D
	format := logging.MustStringFormatter(logformat)
	logbackend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(logbackend, format)
	logging.SetBackend(backendFormatter)
}
