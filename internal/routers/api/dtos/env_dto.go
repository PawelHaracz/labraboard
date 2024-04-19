package dtos

type AddEnvDto struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

type AddVariableDto struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
