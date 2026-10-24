package main

import (
	"os"

	"github.com/android-sms-gateway/server/internal/health"
	smsgateway "github.com/android-sms-gateway/server/internal/sms-gateway"
	"github.com/android-sms-gateway/server/internal/worker"
)

const (
	cmdWorker = "worker"
	cmdHealth = "health"
)

//	@securitydefinitions.basic	ApiAuth
//	@description				User authentication

//	@securitydefinitions.apikey	JWTAuth
//	@in							header
//	@name						Authorization
//	@description				JWT authentication

//	@securitydefinitions.apikey	UserCode
//	@in							header
//	@name						Authorization
//	@description				User one-time code authentication

//	@securitydefinitions.apikey	MobileToken
//	@in							header
//	@name						Authorization
//	@description				Mobile device token

//	@securitydefinitions.apikey	ServerKey
//	@in							header
//	@name						Authorization
//	@description				Private server authentication

//	@title			AndroGate Server API
//	@version		1.0.0
//	@description	Private outbound SMS transport API for CONNECT ME and AndroGate. Features include enqueuing SMS, tracking message delivery, device management, and system health checks.

//	@license.name	Apache 2.0
//	@license.url	https://www.apache.org/licenses/LICENSE-2.0

//	@host		localhost:3000/api
//	@schemes	https http
//
// AndroGate Server Backend.
func main() {
	args := os.Args[1:]
	cmd := "start"
	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case cmdHealth:
		health.Run()
		return
	case cmdWorker:
		worker.Run()
		return
	}

	smsgateway.Run()
}
