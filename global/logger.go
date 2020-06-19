package global

import (
	"github.com/sirupsen/logrus"
)

var RequestLogger = &logrus.Entry{}



func init()  {
	// TODO: requestLogger log format
	// RequestLogger = logrus.WithFields(logrus.Fields{"request_id": request_id, "user_ip": user_ip})
}