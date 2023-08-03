package controller

type SingUpCustomerWithEmailBody struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}
type CustomerBody struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}
type ChangePasswordBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
type ChangeEmaildBody struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}

type EmailLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	Kind     string `json:"kind"`
}
