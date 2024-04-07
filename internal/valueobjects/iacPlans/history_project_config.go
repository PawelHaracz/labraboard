package iacPlans

type HistoryProjectConfig struct {
	GitSha   string
	GitUrl   string
	GitPath  string
	Envs     map[string]string //fixit : change to IaCEnv
	Variable []string
}
