package controller

type SingUpCustomerWithEmailBody struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}
type CustomerBody struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Address    string `json:"address"`
	Profile    string `json:"profile"`
}
type ChangePasswordBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
type ChangeEmaildBody struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}
