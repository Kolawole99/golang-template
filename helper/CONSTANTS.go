package helper

import (
	"os"
	"strconv"
	"time"
)

const COULD_NOT_PROCESS_REQUEST = "Could not process request"
const CURRENT_USER = "Current-User-Id"
const ERROR_PROCESSING_REQUEST = "Error occurred processing request"
const JWT_ISSUER = "JWT_ISSUER"

var envTime, _ = strconv.ParseInt(os.Getenv("DB_CONNECTION_TIMEOUT"), 10, 64)
var DATABASE_TIMEOUT = time.Duration(envTime) * time.Second
