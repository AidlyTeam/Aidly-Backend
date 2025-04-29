package dto

type IDTOManager interface {
	UserManager() *UserDTOManager
	CampaignManager() *CampaignDTOManager
	DonationManager() *DonationDTOManager
}

type DTOManager struct {
	userDTOManager  *UserDTOManager
	campaignManager *CampaignDTOManager
	donationManager *DonationDTOManager
}

func CreateNewDTOManager() IDTOManager {
	userDTOManager := NewUserDTOManager()
	campaingManager := NewCampaignDTOManager()
	donationManager := NewDonationDTOManager()

	return &DTOManager{
		userDTOManager:  &userDTOManager,
		campaignManager: &campaingManager,
		donationManager: &donationManager,
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
