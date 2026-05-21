package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"ecochitas/internal/auth"
)

type auth_claims_context_key struct{}

func (api_handler *Api_handler) with_required_roles(
	allowed_role_name_list []string,
	next_handler_function http.HandlerFunc,
) http.HandlerFunc {
	allowed_role_name_set := make(map[string]struct{}, len(allowed_role_name_list))
	for _, raw_role_name := range allowed_role_name_list {
		normalized_role_name := strings.ToLower(strings.TrimSpace(raw_role_name))
		if normalized_role_name == "" {
			continue
		}
		allowed_role_name_set[normalized_role_name] = struct{}{}
	}

	return func(response_writer http.ResponseWriter, request *http.Request) {
		if api_handler.jwt_authenticator == nil {
			write_json_error(response_writer, http.StatusInternalServerError, "authenticator_not_configured")
			return
		}

		access_token, extract_access_token_error := auth.Extract_bearer_token_from_request(request)
		if extract_access_token_error != nil {
			if errors.Is(extract_access_token_error, auth.Err_missing_authorization_header) ||
				errors.Is(extract_access_token_error, auth.Err_invalid_authorization_scheme) ||
				errors.Is(extract_access_token_error, auth.Err_empty_bearer_token) {
				write_json_error(response_writer, http.StatusUnauthorized, "authorization_header_is_required")
				return
			}

			write_json_error(response_writer, http.StatusUnauthorized, "invalid_authorization_header")
			return
		}

		auth_claims, parse_access_token_error := api_handler.jwt_authenticator.Parse_and_validate_access_token(access_token)
		if parse_access_token_error != nil {
			write_json_error(response_writer, http.StatusUnauthorized, "invalid_or_expired_access_token")
			return
		}

		if len(allowed_role_name_set) > 0 {
			_, role_name_is_allowed := allowed_role_name_set[strings.ToLower(auth_claims.Role_name)]
			if !role_name_is_allowed {
				write_json_error(response_writer, http.StatusForbidden, "insufficient_permissions")
				return
			}
		}

		request_context := context.WithValue(request.Context(), auth_claims_context_key{}, auth_claims)
		next_handler_function(response_writer, request.WithContext(request_context))
	}
}

func get_auth_claims_from_request(request *http.Request) (*auth.Auth_claims, bool) {
	auth_claims, has_auth_claims := request.Context().Value(auth_claims_context_key{}).(*auth.Auth_claims)
	if !has_auth_claims || auth_claims == nil {
		return nil, false
	}

	return auth_claims, true
}
