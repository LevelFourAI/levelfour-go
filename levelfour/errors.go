package levelfour

import (
	"encoding/json"
	"errors"
	"net/http"

	rootpkg "github.com/LevelFourAI/levelfour-go"
	"github.com/LevelFourAI/levelfour-go/core"
)

// StatusCode extracts the HTTP status code from an API error.
// Returns 0 if the error does not wrap a *core.APIError.
func StatusCode(err error) int {
	var apiErr *core.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return 0
}

// ErrorBody extracts the raw response body from an API error.
// Returns an empty string if the error does not wrap a *core.APIError
// or has no body content.
func ErrorBody(err error) string {
	var apiErr *core.APIError
	if errors.As(err, &apiErr) {
		if inner := apiErr.Unwrap(); inner != nil {
			return inner.Error()
		}
	}
	return ""
}

// ErrorMessage extracts the structured error message from an API error.
// First checks typed errors, then falls back to parsing the raw response
// body JSON for an "error.message" field.
func ErrorMessage(err error) string {
	if resp := errorResponse(err); resp != nil && resp.Error != nil {
		return resp.Error.Message
	}
	if msg := jsonField(ErrorBody(err), "error", "message"); msg != "" {
		return msg
	}
	return ""
}

// ErrorCode extracts the structured error code from an API error.
// First checks typed errors, then falls back to parsing the raw response
// body JSON for an "error.code" field.
func ErrorCode(err error) string {
	if resp := errorResponse(err); resp != nil && resp.Error != nil {
		return resp.Error.Code
	}
	if code := jsonField(ErrorBody(err), "error", "code"); code != "" {
		return code
	}
	return ""
}

func errorResponse(err error) *rootpkg.ErrorResponse {
	var badRequest *rootpkg.BadRequestError
	if errors.As(err, &badRequest) && badRequest.Body != nil {
		return badRequest.Body
	}
	var unauthorized *rootpkg.UnauthorizedError
	if errors.As(err, &unauthorized) && unauthorized.Body != nil {
		return unauthorized.Body
	}
	var forbidden *rootpkg.ForbiddenError
	if errors.As(err, &forbidden) && forbidden.Body != nil {
		return forbidden.Body
	}
	var notFound *rootpkg.NotFoundError
	if errors.As(err, &notFound) && notFound.Body != nil {
		return notFound.Body
	}
	var conflict *rootpkg.ConflictError
	if errors.As(err, &conflict) && conflict.Body != nil {
		return conflict.Body
	}
	var tooManyRequests *rootpkg.TooManyRequestsError
	if errors.As(err, &tooManyRequests) && tooManyRequests.Body != nil {
		return tooManyRequests.Body
	}
	var serverError *rootpkg.InternalServerError
	if errors.As(err, &serverError) && serverError.Body != nil {
		return serverError.Body
	}
	return nil
}

func jsonField(body string, keys ...string) string {
	if body == "" {
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return ""
	}
	current := data
	for i, key := range keys {
		val, ok := current[key]
		if !ok {
			return ""
		}
		if i == len(keys)-1 {
			if s, ok := val.(string); ok {
				return s
			}
			return ""
		}
		if m, ok := val.(map[string]interface{}); ok {
			current = m
		} else {
			return ""
		}
	}
	return ""
}

// IsBadRequest reports whether err is a 400 Bad Request API error.
func IsBadRequest(err error) bool {
	return StatusCode(err) == http.StatusBadRequest
}

// IsUnauthorized reports whether err is a 401 Unauthorized API error.
func IsUnauthorized(err error) bool {
	return StatusCode(err) == http.StatusUnauthorized
}

// IsForbidden reports whether err is a 403 Forbidden API error.
func IsForbidden(err error) bool {
	return StatusCode(err) == http.StatusForbidden
}

// IsNotFound reports whether err is a 404 Not Found API error.
func IsNotFound(err error) bool {
	return StatusCode(err) == http.StatusNotFound
}

// IsTimeout reports whether err is a 408 Request Timeout API error.
func IsTimeout(err error) bool {
	return StatusCode(err) == http.StatusRequestTimeout
}

// IsConflict reports whether err is a 409 Conflict API error.
func IsConflict(err error) bool {
	return StatusCode(err) == http.StatusConflict
}

// IsRateLimited reports whether err is a 429 Too Many Requests API error.
func IsRateLimited(err error) bool {
	return StatusCode(err) == http.StatusTooManyRequests
}

// IsValidationError reports whether err is a 422 Unprocessable Entity API error.
func IsValidationError(err error) bool {
	return StatusCode(err) == http.StatusUnprocessableEntity
}

// IsServerError reports whether err is a 5xx server error.
func IsServerError(err error) bool {
	code := StatusCode(err)
	return code >= 500 && code < 600
}
