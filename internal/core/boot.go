package core

func Boot() {
	InitializeLogger()
	InitializeDatabase()
	InitializeRedis()
	InitializeAzureProxy()
}
