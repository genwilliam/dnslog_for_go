package response

const (
	CodeOK              = "ok"
	CodeBadRequest      = "bad_request"
	CodeUnauthorized    = "unauthorized"
	CodeMissingAPIKey   = "missing_key"
	CodeNotFound        = "not_found"
	CodeInvalidAPIKey   = "invalid_api_key"
	CodeInvalidKey      = "invalid_key"
	CodeDisabledAPIKey  = "disabled_key"
	CodeRateLimited     = "rate_limited"
	CodeRateLimitError  = "rate_limit_error"
	CodeRateLimitOff    = "rate_limit_unavailable"
	CodeTokenNotFound   = "token_not_found"
	CodeInternalError   = "internal_error"
	CodeForbidden       = "forbidden"
	CodeSystemPaused    = "system_paused"
	CodeWebhookSecretKeyRequired = "webhook_secret_key_required"
	CodeConflict        = "conflict"
	CodeAPIKeyAlreadyInitialized = "api_key_already_initialized"
)
