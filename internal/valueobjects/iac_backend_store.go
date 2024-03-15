package valueobjects

type IacBackendStore struct {
}

func (store IacBackendStore) GetEnvs() map[string]string {
	return map[string]string{
		"ARM_TENANT_ID":       "4c83ec3e-26b4-444f-afb7-8b171cd1b420",
		"ARM_CLIENT_ID":       "99cc9476-40fd-48b6-813f-e79e0ff830fc",
		"ARM_CLIENT_SECRET":   "fixit",
		"ARM_SUBSCRIPTION_ID": "cb5863b1-784d-4813-b2c7-e87919081ecb",
	}
}
