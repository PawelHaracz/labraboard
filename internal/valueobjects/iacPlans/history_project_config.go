package iacPlans

type HistoryProjectConfig struct {
	GitSha   string
	GitUrl   string
	GitPath  string
	Envs     map[string]string //fixit todo: change to IaCEnv
	Variable map[string]string
}
