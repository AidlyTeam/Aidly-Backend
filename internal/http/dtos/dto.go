package dto

type IDTOManager interface {
	UserManager() *UserDTOManager
	CampaignManager() *CampaignDTOManager
}

type DTOManager struct {
	userDTOManager  *UserDTOManager
	campaignManager *CampaignDTOManager
}

func CreateNewDTOManager() IDTOManager {
	userDTOManager := NewUserDTOManager()
	campaingManager := NewCampaignDTOManager()

	return &DTOManager{
		userDTOManager:  &userDTOManager,
		campaignManager: &campaingManager,
	}
}

func (m *DTOManager) UserManager() *UserDTOManager {
	return m.userDTOManager
}

func (m *DTOManager) CampaignManager() *CampaignDTOManager {
	return m.campaignManager
}
