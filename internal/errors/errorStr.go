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
	ErrUserNotFound        = "USER_NOT_FOUND"
	ErrDefaultUserNotFound = "DEFAULT_USER_NOT_FOUND"
	ErrRoleNotFound        = "ROLE_NOT_FOUND"
	ErrCampaignNotFound    = "CAMPAIGN_NOT_FOUND"
	ErrDonationNotFound    = "DONATION_NOT_FOUND"
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
	ErrFilteringRole      = "ERROR_FILTERING_ROLES"
	ErrFilteringUsers     = "ERROR_FILTERING_USERS"
	ErrFilteringCampaigns = "ERROR_FILTERING_CAMPAIGNS"
	ErrFilteringDonation  = "ERROR_FILTERING_DONATIONS"

	ErrCreatingUser      = "ERROR_CREATE_USER"
	ErrCreatingCampaigns = "ERROR_CREATE_CAMPAINGS"
	ErrCreatingDontaions = "ERROR_CREATE_DONATIONS"

	ErrUpdatingUserRole  = "ERROR_UPDATING_USER_ROLE"
	ErrUpdatingCampaigns = "ERROR_UPDATING_CAMPAINGS"
	ErrUpdatingUsers     = "ERROR_UPDATING_USERS"

	ErrDeletingCampaigns = "ERROR_DELETING_CAMPAIGNS"
	ErrDeletingDonations = "ERROR_DELETING_DONATIONS"

	ErrCountDonations = "ERROR_COUNT_DONATIONS"

	ErrWalletVerificationError = "WALLET_VERIFICATION_ERROR"
)
