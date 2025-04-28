package models

type Defaults struct {
	Roles    RoleDefaults
	Language LanguageDefaults
	Limits   LimitDefaults
}

type RoleDefaults struct {
	DefaultRole string
	FirstRole   string
	RoleAdmin   string
	RolePublic  string
}

type LanguageDefaults struct {
	DefaultLanguage string
}

type LimitDefaults struct {
	DefaultLanguageLimit int
	DefaultUserLimit     int
}

// Initialize Defaults
func NewDefaults() Defaults {
	return Defaults{
		Roles: RoleDefaults{
			DefaultRole: "user",
			FirstRole:   "first",
			RoleAdmin:   "admin",
			RolePublic:  "public",
		},
		Language: LanguageDefaults{
			DefaultLanguage: "EN",
		},
		Limits: LimitDefaults{
			DefaultLanguageLimit: 10,
			DefaultUserLimit:     10,
		},
	}
}
