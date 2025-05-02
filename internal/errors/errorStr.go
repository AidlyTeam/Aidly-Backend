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
	ErrUserNotFound             = "USER_NOT_FOUND"
	ErrDefaultUserNotFound      = "DEFAULT_USER_NOT_FOUND"
	ErrRoleNotFound             = "ROLE_NOT_FOUND"
	ErrCampaignNotFound         = "CAMPAIGN_NOT_FOUND"
	ErrCampaignCategoryNotFound = "CAMPAIGN_CATEGORY_NOT_FOUND"
	ErrDonationNotFound         = "DONATION_NOT_FOUND"
	ErrCategoryNotFound         = "CATEGORY_NOT_FOUND"
	ErrBadgeNotFound            = "BADGE_NOT_FOUND"
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

	// Campaign
	ErrInvalidCampaignStatus  = "INVALID_CAMPAIGN_STATUS"
	ErrCategoryAlreadyAdded   = "CATEGORY_ALREADY_ADDED_INTO_CAMPAIGN"
	ErrUserBadgeAlreadyExists = "USER_BADGE_ALREADY_EXISTS"
)

// 500
const (
	ErrFilteringRole               = "ERROR_FILTERING_ROLES"
	ErrFilteringUsers              = "ERROR_FILTERING_USERS"
	ErrFilteringCampaigns          = "ERROR_FILTERING_CAMPAIGNS"
	ErrFilteringDonation           = "ERROR_FILTERING_DONATIONS"
	ErrFilteringCategories         = "ERROR_FILTERING_CATEGORIES"
	ErrFilteringCampaignCategories = "ERROR_FILTERING_CATEGORIES"
	ErrFilteringBadge              = "ERROR_FILTERING_BADGE"

	ErrCreatingUser             = "ERROR_CREATE_USER"
	ErrCreatingCampaigns        = "ERROR_CREATE_CAMPAINGS"
	ErrCreatingDontaions        = "ERROR_CREATE_DONATIONS"
	ErrCreatingCategories       = "ERROR_CREATING_CATEGORIES"
	ErrCreatingCampaingCategory = "ERROR_CREATING_CAMPAIGN_CATEGORY"
	ErrCreatingBadge            = "ERROR_CREATEING_BADGE"
	ErrCreatingUserBadge        = "ERROR_CREATING_USER_BADGE"

	ErrUpdatingUserRole   = "ERROR_UPDATING_USER_ROLE"
	ErrUpdatingCampaigns  = "ERROR_UPDATING_CAMPAINGS"
	ErrUpdatingUsers      = "ERROR_UPDATING_USERS"
	ErrUpdatingCategories = "ERROR_UPDATING_CATEGORIES"
	ErrUpdatingBadge      = "ERROR_UPDATING_BADGE"

	ErrDeletingCampaigns        = "ERROR_DELETING_CAMPAIGNS"
	ErrDeletingDonations        = "ERROR_DELETING_DONATIONS"
	ErrDeletingCategories       = "ERROR_DELETING_CATEGORIES"
	ErrDeletingCampaignCategory = "ERROR_DELETING_CAMPAING_CATEGORY"
	ErrDeletingBadge            = "ERROR_DELETING_BADGE"

	ErrCountDonations  = "ERROR_COUNT_DONATIONS"
	ErrCountCategories = "ERROR_COUNT_CATEGORIES"
	ErrCountBadge      = "ERROR_COUNT_BADGE"

	ErrWalletVerificationError = "WALLET_VERIFICATION_ERROR"

	ErrDecimalConvertionError = "DECIMAL_CONVERTION_ERROR"
	ErrCommittingTx           = "ERROR_COMMIT_TX"
)
