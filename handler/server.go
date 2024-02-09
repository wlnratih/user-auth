package handler

import (
	"encoding/base64"

	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/golang-jwt/jwt/v4"
)

type Server struct {
	Repository repository.RepositoryInterface
	PrivateKey interface{}
	PublicKey  interface{}
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	PrivateKey string
	PublicKey  string
}

func NewServer(opts NewServerOptions) *Server {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(opts.PrivateKey)
	if err != nil {
		panic(err)
	}

	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	decodedPublicKey, err := base64.StdEncoding.DecodeString(opts.PublicKey)
	if err != nil {
		panic(err)
	}

	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	return &Server{
		Repository: opts.Repository,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}
