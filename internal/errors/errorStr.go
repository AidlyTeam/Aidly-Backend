package serviceErrors

// HTTP Status Codes
const (
	SatatusOK                 = 200
	StatusNotFound            = 404
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

// 404
const (
	ErrLanguageNotFound        = "LANGUAGE_NOT_FOUND"
	ErrLanguageDefaultNotFound = "DEFAULT_LANGUAGE_NOT_FOUND"

	ErrUserProfileNotFound = "USER_PROFILE_NOT_FOUND"
	ErrUserNotFound        = "USER_NOT_FOUND"
	ErrDefaultUserNotFound = "DEFAULT_USER_NOT_FOUND"
	ErrRoleNotFound        = "ROLE_NOT_FOUND"
)

// 400
const (

	// Authentication
	ErrUsernameBeingUsed   = "USERNAME_ALREADY_BEING_USED"
	ErrInvalidAuth         = "USERNAME_OR_PASSWORD_WRONG"
	ErrEmailBeingUsed      = "EMAIL_ALREADY_BEING_USED"
	ErrPasswordsDoNotMatch = "PASSWORDS_DO_NOT_MATCH"

	// General
	ErrInvalidID               = "INVALID_ID"
	ErrInvalidBoolean          = "INVALID_BOOLEAN"
	ErrInvalidWalletConnection = "INVALID_WALLET_CONNECTION"
	ErrInvalidFileType         = "INVALID_FILE_TYPE"
)

// 500
const (
	ErrFilteringUserProfile = "ERROR_FILTERING_USER_PROFILE"
	ErrFilteringRole        = "ERROR_FILTERING_ROLES"
	ErrFilteringUsers       = "ERROR_FILTERING_USERS"
	ErrFilteringLanguages   = "ERROR_FILTERING_LANGUAGES"
	ErrCreatingUser         = "ERROR_CREATE_USER"

	ErrCreatingUsers       = "ERROR_CREATING_USERS"
	ErrUpdatingUsers       = "ERROR_UPDATING_USERS"
	ErrChangingRole        = "ERROR_CHANGING_USER_ROLE"
	ErrCreatingUserProfile = "ERROR_CREATING_USER_PROFILE"
	ErrUpdatingUserProfile = "ERROR_UPDATING_USER_PROFILE"

	ErrComparingPassword = "ERROR_COMPARING_PASSWORDS"
	ErrHashing           = "ERROR_HASHING"
	ErrTransactionError  = "ERROR_COMMITING"

	ErrWalletVerificationError = "WALLET_VERIFICATION_ERROR"
)
