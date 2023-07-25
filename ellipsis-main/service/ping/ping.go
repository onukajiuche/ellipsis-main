package ping

type HealthService interface {
	ReturnTrue() bool
}

type healthService struct{}

func NewHealthService() HealthService {
	return healthService{}
}

func (healthService) ReturnTrue() bool {
	return true
}
