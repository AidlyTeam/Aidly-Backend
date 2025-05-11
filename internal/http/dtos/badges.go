package dto

import (
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
)

type BadgeDTOManager struct{}

// NewBadgeDTOManager returns a new instance of BadgeDTOManager
func NewBadgeDTOManager() BadgeDTOManager {
	return BadgeDTOManager{}
}

// MetadataAttributeView represents an individual attribute in NFT metadata
type MetadataAttributeView struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

// MetadataView represents the structure of an NFT metadata in responses
type MetadataView struct {
	Name        string                  `json:"name"`
	Symbol      string                  `json:"symbol"`
	Description string                  `json:"description"`
	Image       string                  `json:"image"`
	URI         string                  `json:"uri"`
	Attributes  []MetadataAttributeView `json:"attributes,omitempty"`
	SellerFee   int                     `json:"seller_fee_basis_points"`
}

func (BadgeDTOManager) ToMetadataView(badge *repo.TBadge) *MetadataView {
	if badge == nil || !badge.IsNft {
		return nil
	}

	return &MetadataView{
		Name:        badge.Name,
		Symbol:      badge.Symbol.String,
		Description: badge.Description.String,
		Image:       badge.IconPath.String,
		URI:         badge.Uri.String,
		SellerFee:   int(badge.SellerFee.Int32),
	}
}

// BadgeView represents the structure of a badge in responses
type BadgeView struct {
	ID          string `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconPath    string `json:"iconPath"`
	Uri         string `json:"uri"`
	SellerFee   int32  `json:"sellerFee"`
	Threshold   int32  `json:"threshold"`
	IsNft       bool   `json:"isNft"`
	IsMinted    bool   `json:"isMinted"`
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
		Symbol:      badge.Symbol.String,
		Name:        badge.Name,
		Description: badge.Description.String,
		SellerFee:   badge.SellerFee.Int32,
		Uri:         badge.Uri.String,
		IsNft:       badge.IsNft,
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

// ToBadgeView transforms domain badge to BadgeView
func (BadgeDTOManager) ToUserBadgeView(badge *repo.GetUserBadgesRow) *BadgeView {
	if badge == nil {
		return nil
	}

	return &BadgeView{
		ID:          badge.ID.String(),
		Symbol:      badge.Symbol.String,
		Name:        badge.Name,
		Description: badge.Description.String,
		SellerFee:   badge.SellerFee.Int32,
		Uri:         badge.Uri.String,
		IsNft:       badge.IsNft,
		IconPath:    badge.IconPath.String,
		Threshold:   badge.DonationThreshold,
		IsMinted:    badge.IsMinted,
	}
}

// ToBadgeViews transforms list of badges to BadgeViews
func (m *BadgeDTOManager) ToUserBadgeViews(badges []repo.GetUserBadgesRow, count int64) *BadgeViews {
	var badgeViews []BadgeView
	for _, badge := range badges {
		badgeViews = append(badgeViews, *m.ToUserBadgeView(&badge))
	}
	return &BadgeViews{Badges: badgeViews, TotalCount: int(count)}
}

// BadgeCreateDTO is used to create a new badge
type BadgeCreateDTO struct {
	Symbol            string `json:"symbol"`
	Name              string `json:"name" validate:"required,min=3,max=20"`
	Description       string `json:"description" validate:"required,min=5,max=100"`
	SellerFee         int32  `json:"sellerFee"`
	DonationThreshold int32  `json:"donationThreshold"`
	IsNft             bool   `json:"isNft"`
}

// BadgeUpdateDTO is used to update an existing badge
type BadgeUpdateDTO struct {
	Symbol            string `json:"symbol"`
	Name              string `json:"name" validate:"omitempty,min=3,max=20"`
	Description       string `json:"description" validate:"omitempty,min=5,max=100"`
	SellerFee         int32  `json:"sellerFee"`
	DonationThreshold int32  `json:"donationThreshold" validate:"omitempty"`
}
