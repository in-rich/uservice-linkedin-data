package models

type UpsertUser struct {
	FirstName string `json:"firstName" validate:"required,max=255"`
	LastName  string `json:"lastName" validate:"required,max=255"`
	// In the hypothetical case that 1 character = 1 Bit, this allows for a maximum of 20MB.
	// Since LinkedIn upload size is capped at 8Mb, and there is some expansion due to base64
	// encoding at roughly 135%, this gives us some head, while preventing abuse.
	ProfilePicture string `json:"profilePicture" validate:"max=20971520"`
}
