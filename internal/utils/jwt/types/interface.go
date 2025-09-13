package types

import "dengovie/internal/web"

type Processor interface {
	VerifyJWT(string) (map[web.JWTKey]any, error)
	Sign(data ...any) (string, error)
}
