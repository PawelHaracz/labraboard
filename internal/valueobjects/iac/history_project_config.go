package iac

import "labraboard/internal/valueobjects"

type HistoryProjectConfig struct {
	GitSha   string
	GitUrl   string
	GitPath  string
	Envs     []valueobjects.IaCEnv
	Variable map[string]string
}
