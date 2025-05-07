package dto

import (
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
)

type BadgeDTOManager struct{}

// NewBadgeDTOManager returns a new instance of BadgeDTOManager
func NewBadgeDTOManager() BadgeDTOManager {
	return BadgeDTOManager{}
}

// BadgeView represents the structure of a badge in responses
type BadgeView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconPath    string `json:"iconPath"`
	Threshold   int32  `json:"threshold"`
}

type BadgeViews struct {
	Badges     []BadgeView `json:"badges"`
	TotalCount int         `json:"totalCount,omitempty"`
}

// ToBadgeView transforms domain badge to BadgeView
func (BadgeDTOManager) ToBadgeView(badge *repo.TBadge) *BadgeView {
	if badge == nil {
		return nil
	}

	return &BadgeView{
		ID:          badge.ID.String(),
		Name:        badge.Name,
		Description: badge.Description.String,
		IconPath:    badge.IconPath.String,
		Threshold:   badge.DonationThreshold,
	}
}

// ToBadgeViews transforms list of badges to BadgeViews
func (m *BadgeDTOManager) ToBadgeViews(badges []repo.TBadge, count int64) *BadgeViews {
	var badgeViews []BadgeView
	for _, badge := range badges {
		badgeViews = append(badgeViews, *m.ToBadgeView(&badge))
	}
	return &BadgeViews{Badges: badgeViews, TotalCount: int(count)}
}

// BadgeCreateDTO is used to create a new badge
type BadgeCreateDTO struct {
	Name              string `json:"name" validate:"required,min=3,max=20"`
	Description       string `json:"description" validate:"required,min=5,max=100"`
	DonationThreshold int32  `json:"donationThreshold"`
}

// BadgeUpdateDTO is used to update an existing badge
type BadgeUpdateDTO struct {
	Name              string `json:"name" validate:"omitempty,min=3,max=20"`
	Description       string `json:"description" validate:"omitempty,min=5,max=100"`
	DonationThreshold int32  `json:"donationThreshold" validate:"omitempty"`
}
