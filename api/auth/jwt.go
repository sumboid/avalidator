package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthRole int

const (
	StudentRole AuthRole = iota
	AdminRole
)

var authRoleNames = map[AuthRole]string{
	StudentRole: "student",
	AdminRole:   "admin",
}

var authRolesByName = map[string]AuthRole{}

func init() {
	for role, name := range authRoleNames {
		authRolesByName[name] = role
	}
}

func StringToAuthRole(roleName string) (AuthRole, bool) {
	if role, ok := authRolesByName[roleName]; ok {
		return role, true
	}

	return -1, false
}

func (r AuthRole) String() string {
	return authRoleNames[r]
}

func (r AuthRole) AllowedRoles() []AuthRole {
	if r == StudentRole {
		return []AuthRole{StudentRole}
	}

	if r == AdminRole {
		return []AuthRole{StudentRole, AdminRole}
	}

	return []AuthRole{}
}

func (r AuthRole) AllowedRolesString() []string {
	allowedRoles := r.AllowedRoles()
	res := make([]string, len(allowedRoles))

	for i := range allowedRoles {
		res[i] = allowedRoles[i].String()
	}

	return res
}

type hasuraClaimT struct {
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	DefaultRole  string   `json:"x-hasura-default-role"`
	UserID       string   `json:"x-hasura-user-id"`
}

type claimsT struct {
	Hasura hasuraClaimT `json:"hasura"`
	jwt.StandardClaims
}

type JWTService interface {
	CreateToken(string, string) (string, string, error)
	RefreshToken(string) (string, string, error)
	RemoveToken(string) error
}

type jwtService struct {
	client *redis.Client
}

func NewJWTService(client *redis.Client) JWTService {
	return &jwtService{
		client: client,
	}
}

func (s *jwtService) CreateToken(id, roleStr string) (string, string, error) {
	role, ok := StringToAuthRole(roleStr)
	if !ok {
		return "", "", fmt.Errorf("Failed to convert role: %s", roleStr)
	}

	authClaims := &claimsT{
		Hasura: hasuraClaimT{
			UserID:       id,
			DefaultRole:  roleStr,
			AllowedRoles: role.AllowedRolesString(),
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(Config.JWT.AuthTTL) * time.Minute).Unix(),
		},
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims)
	authTokenString, err := authToken.SignedString([]byte(Config.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return "", "", err
	}

	refreshTokenString := refreshToken.String()
	if err != nil {
		return "", "", err
	}

	s.client.Set(refreshTokenString, authTokenString, time.Duration(Config.JWT.RefreshTTL)*time.Minute)

	return authTokenString, refreshTokenString, nil
}

func (s *jwtService) RefreshToken(refreshToken string) (string, string, error) {
	authToken, err := s.client.Get(refreshToken).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", fmt.Errorf("Failed to lookup refresh token")
		}

		return "", "", err
	}

	s.client.Expire(refreshToken, time.Minute)

	claims := &claimsT{}

	_, err = jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWT.Secret), nil
	})

	if err != nil {
		return "", "", err
	}

	id := claims.Hasura.UserID
	role := claims.Hasura.DefaultRole

	return s.CreateToken(id, role)
}

func (s *jwtService) RemoveToken(refreshToken string) error {
	_, err := s.client.Del(refreshToken).Result()

	return err
}
