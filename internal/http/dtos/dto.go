package dto

type IDTOManager interface {
	UserManager() *UserDTOManager
}

type DTOManager struct {
	userDTOManager *UserDTOManager
}

func CreateNewDTOManager() IDTOManager {
	userDTOManager := NewUserDTOManager()

	return &DTOManager{
		userDTOManager: &userDTOManager,
	}
}

func (m *DTOManager) UserManager() *UserDTOManager {
	return m.userDTOManager
}
