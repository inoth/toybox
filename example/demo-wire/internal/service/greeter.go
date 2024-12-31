package service

type GreeterService struct{}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHi(name string) string {
	return "hello " + name
}
