package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Err_missing_authorization_header = errors.New("missing_authorization_header")
	Err_invalid_authorization_scheme = errors.New("invalid_authorization_scheme")
	Err_empty_bearer_token           = errors.New("empty_bearer_token")
	Err_unsupported_role_name        = errors.New("unsupported_role_name")
)

type Auth_claims struct {
	User_identifier string `json:"user_identifier,omitempty"`
	Role_name       string `json:"role_name"`
	Full_name       string `json:"full_name,omitempty"`
	jwt.RegisteredClaims
}

type Generate_access_token_command struct {
	User_identifier string
	Role_name       string
	Full_name       string
}

type Jwt_authenticator struct {
	jwt_signing_key    []byte
	jwt_issuer         string
	jwt_audience       string
	access_token_ttl   time.Duration
	allowed_algorithms []string
}

func New_jwt_authenticator(
	jwt_signing_key string,
	jwt_issuer string,
	jwt_audience string,
	access_token_ttl_minutes int,
) (*Jwt_authenticator, error) {
	normalized_signing_key := strings.TrimSpace(jwt_signing_key)
	if normalized_signing_key == "" {
		return nil, fmt.Errorf("jwt_signing_key_is_required")
	}

	normalized_jwt_issuer := strings.TrimSpace(jwt_issuer)
	if normalized_jwt_issuer == "" {
		return nil, fmt.Errorf("jwt_issuer_is_required")
	}

	normalized_jwt_audience := strings.TrimSpace(jwt_audience)
	if normalized_jwt_audience == "" {
		return nil, fmt.Errorf("jwt_audience_is_required")
	}

	if access_token_ttl_minutes <= 0 {
		return nil, fmt.Errorf("access_token_ttl_minutes_must_be_greater_than_zero")
	}

	return &Jwt_authenticator{
		jwt_signing_key:    []byte(normalized_signing_key),
		jwt_issuer:         normalized_jwt_issuer,
		jwt_audience:       normalized_jwt_audience,
		access_token_ttl:   time.Duration(access_token_ttl_minutes) * time.Minute,
		allowed_algorithms: []string{jwt.SigningMethodHS256.Alg()},
	}, nil
}

func (jwt_authenticator *Jwt_authenticator) Generate_access_token(
	generate_access_token_command Generate_access_token_command,
) (string, time.Time, error) {
	normalized_user_identifier := strings.TrimSpace(generate_access_token_command.User_identifier)
	if normalized_user_identifier == "" {
		return "", time.Time{}, fmt.Errorf("user_identifier_is_required")
	}

	normalized_role_name, normalize_role_name_error := normalize_role_name(generate_access_token_command.Role_name)
	if normalize_role_name_error != nil {
		return "", time.Time{}, normalize_role_name_error
	}

	token_issued_at := time.Now().UTC()
	token_expires_at := token_issued_at.Add(jwt_authenticator.access_token_ttl)

	auth_claims := Auth_claims{
		User_identifier: normalized_user_identifier,
		Role_name:       normalized_role_name,
		Full_name:       strings.TrimSpace(generate_access_token_command.Full_name),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   normalized_user_identifier,
			Issuer:    jwt_authenticator.jwt_issuer,
			Audience:  jwt.ClaimStrings{jwt_authenticator.jwt_audience},
			IssuedAt:  jwt.NewNumericDate(token_issued_at),
			NotBefore: jwt.NewNumericDate(token_issued_at.Add(-5 * time.Second)),
			ExpiresAt: jwt.NewNumericDate(token_expires_at),
		},
	}

	unsigned_token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth_claims)
	signed_access_token, sign_token_error := unsigned_token.SignedString(jwt_authenticator.jwt_signing_key)
	if sign_token_error != nil {
		return "", time.Time{}, fmt.Errorf("failed_to_sign_access_token: %w", sign_token_error)
	}

	return signed_access_token, token_expires_at, nil
}

func (jwt_authenticator *Jwt_authenticator) Parse_and_validate_access_token(
	access_token string,
) (*Auth_claims, error) {
	normalized_access_token := strings.TrimSpace(access_token)
	if normalized_access_token == "" {
		return nil, Err_empty_bearer_token
	}

	auth_claims := &Auth_claims{}
	parsed_token, parse_access_token_error := jwt.ParseWithClaims(
		normalized_access_token,
		auth_claims,
		func(token_to_validate *jwt.Token) (any, error) {
			return jwt_authenticator.jwt_signing_key, nil
		},
		jwt.WithValidMethods(jwt_authenticator.allowed_algorithms),
		jwt.WithIssuer(jwt_authenticator.jwt_issuer),
		jwt.WithAudience(jwt_authenticator.jwt_audience),
		jwt.WithLeeway(30*time.Second),
	)
	if parse_access_token_error != nil {
		return nil, fmt.Errorf("failed_to_parse_access_token: %w", parse_access_token_error)
	}

	if !parsed_token.Valid {
		return nil, fmt.Errorf("access_token_is_not_valid")
	}

	normalized_role_name, normalize_role_name_error := normalize_role_name(auth_claims.Role_name)
	if normalize_role_name_error != nil {
		return nil, normalize_role_name_error
	}
	auth_claims.Role_name = normalized_role_name

	if strings.TrimSpace(auth_claims.Subject) == "" {
		return nil, fmt.Errorf("subject_claim_is_required")
	}

	if strings.TrimSpace(auth_claims.User_identifier) == "" {
		auth_claims.User_identifier = strings.TrimSpace(auth_claims.Subject)
	}

	return auth_claims, nil
}

func Extract_bearer_token_from_request(request *http.Request) (string, error) {
	authorization_header_value := strings.TrimSpace(request.Header.Get("Authorization"))
	if authorization_header_value == "" {
		return "", Err_missing_authorization_header
	}

	authorization_header_parts := strings.SplitN(authorization_header_value, " ", 2)
	if len(authorization_header_parts) != 2 {
		return "", Err_invalid_authorization_scheme
	}

	authorization_scheme := strings.ToLower(strings.TrimSpace(authorization_header_parts[0]))
	if authorization_scheme != "bearer" {
		return "", Err_invalid_authorization_scheme
	}

	bearer_token := strings.TrimSpace(authorization_header_parts[1])
	if bearer_token == "" {
		return "", Err_empty_bearer_token
	}

	return bearer_token, nil
}

func normalize_role_name(raw_role_name string) (string, error) {
	normalized_role_name := strings.ToLower(strings.TrimSpace(raw_role_name))
	switch normalized_role_name {
	case "citizen", "driver", "admin", "condominium_admin":
		return normalized_role_name, nil
	default:
		return "", fmt.Errorf("%w: %s", Err_unsupported_role_name, raw_role_name)
	}
}
