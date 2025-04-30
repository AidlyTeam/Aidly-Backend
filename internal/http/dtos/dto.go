package dto

type IDTOManager interface {
	UserManager() *UserDTOManager
	CampaignManager() *CampaignDTOManager
	DonationManager() *DonationDTOManager
	CategoryManager() *CategoryDTOManager
}

type DTOManager struct {
	userDTOManager  *UserDTOManager
	campaignManager *CampaignDTOManager
	donationManager *DonationDTOManager
	categoryManager *CategoryDTOManager
}

func CreateNewDTOManager() IDTOManager {
	userDTOManager := NewUserDTOManager()
	campaingManager := NewCampaignDTOManager()
	donationManager := NewDonationDTOManager()
	categoryManager := NewCategoryDTOManager()

	return &DTOManager{
		userDTOManager:  &userDTOManager,
		campaignManager: &campaingManager,
		donationManager: &donationManager,
		categoryManager: &categoryManager,
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
