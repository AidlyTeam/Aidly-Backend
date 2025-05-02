package dto

type IDTOManager interface {
	UserManager() *UserDTOManager
	CampaignManager() *CampaignDTOManager
	DonationManager() *DonationDTOManager
	CategoryManager() *CategoryDTOManager
	BadgeManager() *BadgeDTOManager
}

type DTOManager struct {
	userDTOManager  *UserDTOManager
	campaignManager *CampaignDTOManager
	donationManager *DonationDTOManager
	categoryManager *CategoryDTOManager
	badgeManager    *BadgeDTOManager
}

func CreateNewDTOManager() IDTOManager {
	userDTOManager := NewUserDTOManager()
	campaingManager := NewCampaignDTOManager()
	donationManager := NewDonationDTOManager()
	categoryManager := NewCategoryDTOManager()
	badgeManager := NewBadgeDTOManager()

	return &DTOManager{
		userDTOManager:  &userDTOManager,
		campaignManager: &campaingManager,
		donationManager: &donationManager,
		categoryManager: &categoryManager,
		badgeManager:    &badgeManager,
	}
}

func (m *DTOManager) UserManager() *UserDTOManager {
	return m.userDTOManager
}

func (m *DTOManager) CampaignManager() *CampaignDTOManager {
	return m.campaignManager
}

func (m *DTOManager) DonationManager() *DonationDTOManager {
	return m.donationManager
}

func (m *DTOManager) CategoryManager() *CategoryDTOManager {
	return m.categoryManager
}

func (m *DTOManager) BadgeManager() *BadgeDTOManager {
	return m.badgeManager
}
