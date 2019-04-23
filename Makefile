default:
	GODEBUG=http2debug=2  go run src/main/mywrapper.go -id $(CLIENT_ID) -key $(KEY)
