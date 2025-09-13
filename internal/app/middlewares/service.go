package middlewares

import jwtTypes "dengovie/internal/utils/jwt/types"

type service struct {
	jwt jwtTypes.Processor
}

func New(jwt jwtTypes.Processor) *service {
	return &service{
		jwt: jwt,
	}
}
