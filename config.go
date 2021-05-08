package main

import "os"

var (
	igSessionId       = os.Getenv("IG_SESSION_ID1")
	YarbDBIp          = os.Getenv("YARB_DB_IP")
	YarbDBPort        = os.Getenv("YARB_DB_PORT")
	YarbDBDomain      = os.Getenv("YARB_DB_DMN")
	YarbMakabaApiURL  = os.Getenv("YARB_MAKABA_API_URL")
	yarbBasicAuthUser = os.Getenv("YARB_BASIC_AUTH_USER")
	yarbBasicAuthPass = os.Getenv("YARB_BASIC_AUTH_PASS")
)