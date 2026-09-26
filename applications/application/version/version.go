package version

var Type = "dev"
var SemVer = "unknown"
var Name = "unknown"

type Version struct {
	Type   string `json:"type"`
	SemVer string `json:"semver"`
	Name   string `json:"name"`
}
