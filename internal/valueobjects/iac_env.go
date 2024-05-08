package valueobjects

const SECRET_VALUE_HASH = "***"

type IaCEnv struct {
	Name      string
	Value     string
	HasSecret bool
}
