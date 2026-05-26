package levelfour

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	rootpkg "github.com/LevelFourAI/levelfour-go"
	"github.com/LevelFourAI/levelfour-go/core"
	"github.com/stretchr/testify/assert"
)

func TestStatusCode(t *testing.T) {
	t.Run("raw APIError", func(t *testing.T) {
		err := core.NewAPIError(404, nil, errors.New("not found"))
		assert.Equal(t, 404, StatusCode(err))
	})

	t.Run("wrapped APIError", func(t *testing.T) {
		apiErr := core.NewAPIError(403, nil, errors.New("forbidden"))
		wrapped := fmt.Errorf("request failed: %w", apiErr)
		assert.Equal(t, 403, StatusCode(wrapped))
	})

	t.Run("non-API error", func(t *testing.T) {
		err := errors.New("network timeout")
		assert.Equal(t, 0, StatusCode(err))
	})

	t.Run("nil error", func(t *testing.T) {
		assert.Equal(t, 0, StatusCode(nil))
	})
}

func TestErrorBody(t *testing.T) {
	t.Run("extracts body from APIError", func(t *testing.T) {
		err := core.NewAPIError(400, nil, errors.New(`{"error":"bad request"}`))
		assert.Equal(t, `{"error":"bad request"}`, ErrorBody(err))
	})

	t.Run("extracts body from wrapped APIError", func(t *testing.T) {
		apiErr := core.NewAPIError(500, nil, errors.New("internal server error"))
		wrapped := fmt.Errorf("call failed: %w", apiErr)
		assert.Equal(t, "internal server error", ErrorBody(wrapped))
	})

	t.Run("extracts body from typed error", func(t *testing.T) {
		body := `{"success":false,"error":{"code":"FORBIDDEN","message":"access denied"}}`
		apiErr := core.NewAPIError(403, nil, errors.New(body))
		forbidden := &rootpkg.ForbiddenError{APIError: apiErr}
		assert.Equal(t, body, ErrorBody(forbidden))
	})

	t.Run("empty for non-API error", func(t *testing.T) {
		assert.Equal(t, "", ErrorBody(errors.New("network error")))
	})

	t.Run("empty for nil", func(t *testing.T) {
		assert.Equal(t, "", ErrorBody(nil))
	})

	t.Run("empty for APIError with nil inner", func(t *testing.T) {
		err := core.NewAPIError(404, nil, nil)
		assert.Equal(t, "", ErrorBody(err))
	})
}

func TestErrorMessage(t *testing.T) {
	t.Run("extracts from ForbiddenError", func(t *testing.T) {
		apiErr := core.NewAPIError(403, nil, errors.New("raw"))
		forbidden := &rootpkg.ForbiddenError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "FORBIDDEN",
					Message: "access denied",
				},
			},
		}
		assert.Equal(t, "access denied", ErrorMessage(forbidden))
	})

	t.Run("extracts from UnauthorizedError", func(t *testing.T) {
		apiErr := core.NewAPIError(401, nil, errors.New("raw"))
		unauthorized := &rootpkg.UnauthorizedError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "UNAUTHORIZED",
					Message: "invalid token",
				},
			},
		}
		assert.Equal(t, "invalid token", ErrorMessage(unauthorized))
	})

	t.Run("extracts from UnprocessableEntityError via JSON fallback", func(t *testing.T) {
		body := `{"error":{"code":"VALIDATION_ERROR","message":"field is required"}}`
		apiErr := core.NewAPIError(422, nil, errors.New(body))
		unprocessable := &rootpkg.UnprocessableEntityError{APIError: apiErr}
		assert.Equal(t, "field is required", ErrorMessage(unprocessable))
	})

	t.Run("extracts from BadRequestError", func(t *testing.T) {
		apiErr := core.NewAPIError(400, nil, errors.New("raw"))
		badReq := &rootpkg.BadRequestError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "BAD_REQUEST",
					Message: "invalid parameter",
				},
			},
		}
		assert.Equal(t, "invalid parameter", ErrorMessage(badReq))
	})

	t.Run("extracts from NotFoundError", func(t *testing.T) {
		apiErr := core.NewAPIError(404, nil, errors.New("raw"))
		notFound := &rootpkg.NotFoundError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			},
		}
		assert.Equal(t, "resource not found", ErrorMessage(notFound))
	})

	t.Run("extracts from ConflictError", func(t *testing.T) {
		apiErr := core.NewAPIError(409, nil, errors.New("raw"))
		conflict := &rootpkg.ConflictError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "CONFLICT",
					Message: "already exists",
				},
			},
		}
		assert.Equal(t, "already exists", ErrorMessage(conflict))
	})

	t.Run("extracts from TooManyRequestsError", func(t *testing.T) {
		apiErr := core.NewAPIError(429, nil, errors.New("raw"))
		rateLimit := &rootpkg.TooManyRequestsError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "RATE_LIMITED",
					Message: "slow down",
				},
			},
		}
		assert.Equal(t, "slow down", ErrorMessage(rateLimit))
	})

	t.Run("extracts from InternalServerError", func(t *testing.T) {
		apiErr := core.NewAPIError(500, nil, errors.New("raw"))
		serverErr := &rootpkg.InternalServerError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "INTERNAL_ERROR",
					Message: "something went wrong",
				},
			},
		}
		assert.Equal(t, "something went wrong", ErrorMessage(serverErr))
	})

	t.Run("falls back to JSON body parsing", func(t *testing.T) {
		body := `{"error":{"code":"NOT_FOUND","message":"resource not found"}}`
		err := core.NewAPIError(404, nil, errors.New(body))
		assert.Equal(t, "resource not found", ErrorMessage(err))
	})

	t.Run("empty for non-JSON body", func(t *testing.T) {
		err := core.NewAPIError(500, nil, errors.New("Internal Server Error"))
		assert.Equal(t, "", ErrorMessage(err))
	})

	t.Run("empty for nil", func(t *testing.T) {
		assert.Equal(t, "", ErrorMessage(nil))
	})
}

func TestErrorCode(t *testing.T) {
	t.Run("extracts from ForbiddenError", func(t *testing.T) {
		apiErr := core.NewAPIError(403, nil, errors.New("raw"))
		forbidden := &rootpkg.ForbiddenError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "FORBIDDEN",
					Message: "access denied",
				},
			},
		}
		assert.Equal(t, "FORBIDDEN", ErrorCode(forbidden))
	})

	t.Run("extracts from BadRequestError", func(t *testing.T) {
		apiErr := core.NewAPIError(400, nil, errors.New("raw"))
		badReq := &rootpkg.BadRequestError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "BAD_REQUEST",
					Message: "invalid parameter",
				},
			},
		}
		assert.Equal(t, "BAD_REQUEST", ErrorCode(badReq))
	})

	t.Run("extracts from NotFoundError", func(t *testing.T) {
		apiErr := core.NewAPIError(404, nil, errors.New("raw"))
		notFound := &rootpkg.NotFoundError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			},
		}
		assert.Equal(t, "NOT_FOUND", ErrorCode(notFound))
	})

	t.Run("extracts from TooManyRequestsError", func(t *testing.T) {
		apiErr := core.NewAPIError(429, nil, errors.New("raw"))
		rateLimit := &rootpkg.TooManyRequestsError{
			APIError: apiErr,
			Body: &rootpkg.ErrorResponse{
				Error: &rootpkg.ErrorResponseError{
					Code:    "RATE_LIMITED",
					Message: "too many requests",
				},
			},
		}
		assert.Equal(t, "RATE_LIMITED", ErrorCode(rateLimit))
	})

	t.Run("falls back to JSON body parsing", func(t *testing.T) {
		body := `{"error":{"code":"RATE_LIMITED","message":"too many requests"}}`
		err := core.NewAPIError(429, nil, errors.New(body))
		assert.Equal(t, "RATE_LIMITED", ErrorCode(err))
	})

	t.Run("empty for non-typed error without JSON", func(t *testing.T) {
		assert.Equal(t, "", ErrorCode(core.NewAPIError(500, nil, errors.New("err"))))
	})

	t.Run("empty for nil", func(t *testing.T) {
		assert.Equal(t, "", ErrorCode(nil))
	})
}

func TestIsBadRequest(t *testing.T) {
	assert.True(t, IsBadRequest(core.NewAPIError(http.StatusBadRequest, nil, errors.New("bad"))))
	assert.False(t, IsBadRequest(core.NewAPIError(http.StatusNotFound, nil, errors.New("nope"))))
}

func TestIsUnauthorized(t *testing.T) {
	assert.True(t, IsUnauthorized(core.NewAPIError(http.StatusUnauthorized, nil, errors.New("unauth"))))
	assert.False(t, IsUnauthorized(core.NewAPIError(http.StatusOK, nil, errors.New("ok"))))
}

func TestIsForbidden(t *testing.T) {
	assert.True(t, IsForbidden(core.NewAPIError(http.StatusForbidden, nil, errors.New("forbidden"))))
	assert.False(t, IsForbidden(core.NewAPIError(http.StatusUnauthorized, nil, errors.New("unauth"))))
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(core.NewAPIError(http.StatusNotFound, nil, errors.New("missing"))))
	assert.False(t, IsNotFound(core.NewAPIError(http.StatusOK, nil, errors.New("ok"))))
	assert.False(t, IsNotFound(errors.New("not an API error")))
}

func TestIsTimeout(t *testing.T) {
	assert.True(t, IsTimeout(core.NewAPIError(http.StatusRequestTimeout, nil, errors.New("timeout"))))
	assert.False(t, IsTimeout(core.NewAPIError(http.StatusOK, nil, errors.New("ok"))))
}

func TestIsConflict(t *testing.T) {
	assert.True(t, IsConflict(core.NewAPIError(http.StatusConflict, nil, errors.New("conflict"))))
	assert.False(t, IsConflict(core.NewAPIError(http.StatusOK, nil, errors.New("ok"))))
}

func TestIsRateLimited(t *testing.T) {
	assert.True(t, IsRateLimited(core.NewAPIError(http.StatusTooManyRequests, nil, errors.New("slow down"))))
	assert.False(t, IsRateLimited(core.NewAPIError(http.StatusOK, nil, errors.New("ok"))))
}

func TestIsValidationError(t *testing.T) {
	assert.True(t, IsValidationError(core.NewAPIError(http.StatusUnprocessableEntity, nil, errors.New("invalid"))))
	assert.False(t, IsValidationError(core.NewAPIError(http.StatusBadRequest, nil, errors.New("bad"))))
}

func TestIsServerError(t *testing.T) {
	for _, code := range []int{500, 502, 503, 504} {
		assert.True(t, IsServerError(core.NewAPIError(code, nil, errors.New("server"))), "expected true for %d", code)
	}
	assert.False(t, IsServerError(core.NewAPIError(499, nil, errors.New("client"))))
	assert.False(t, IsServerError(core.NewAPIError(400, nil, errors.New("bad"))))
	assert.False(t, IsServerError(errors.New("not API")))
}
