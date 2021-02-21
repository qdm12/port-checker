package health

func MakeIsHealthy() func() error {
	return isHealthy
}

func isHealthy() (err error) {
	return nil
}
