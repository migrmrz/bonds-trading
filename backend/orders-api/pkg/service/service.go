package service

type Service struct {
	storeClient storeClient
	redisClient redisClient
}

func New(sc storeClient, rc redisClient) *Service {
	return &Service{
		storeClient: sc,
		redisClient: rc,
	}
}
