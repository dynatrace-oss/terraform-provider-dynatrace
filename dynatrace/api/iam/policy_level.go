package iam

type PolicyLevel string

var PolicyLevels = struct {
	Global      PolicyLevel
	Account     PolicyLevel
	Environment PolicyLevel
}{
	PolicyLevel("global"),
	PolicyLevel("account"),
	PolicyLevel("environment"),
}
