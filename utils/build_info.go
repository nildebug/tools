package utils

import "fmt"

var buildTime string
var buildDev string = "true"

//CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
//-ldflags="-w -s -X 'github.com/nildebug/tools/utils.buildTime=$(date +"%Y-%m-%d %H:%M:%S")' -X 'github.com/nildebug/tools/utils.buildDev=false'"

func IsDev() bool {
	return buildDev == "true"
}

func InputBuildInfo() string {
	return fmt.Sprintf("Build Time: %s  Build Dev: %s", buildTime, buildDev)
}
