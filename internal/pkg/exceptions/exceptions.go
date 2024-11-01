package exceptions

import (
	"fmt"
)

type Exception struct {
	errCode uint
	errMsg  string
}

var (
	InvalidSessionStateException = Exception{101, "Invalid State parameter in Session"}
	AuthExchangeFailureException = Exception{102, "Failure to exchange an authorization code for a token"}
	IDTokenVerificationException = Exception{103, "Failed to verify ID Token"}
	SessionExpiredException      = Exception{104, "Session has Expired"}
)

func (err Exception) Error() string {
	return fmt.Sprintf("Error code %d: %v", err.errCode, err.errMsg)
}
