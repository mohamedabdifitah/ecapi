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
type SingUpMerchantWithPhoneBody struct {
	BusinessPhone string `json:"business_phone"`
	Password      string `json:"password"`
}
type MerchantBody struct {
	Location           []float64 `json:"location"`
	Address            string    `json:"address"`
	BusinessEmail      string    `json:"business_email"`
	TimeOperationStart int       `json:"time_operation_start"`
	TimeOperationEnd   int       `json:"time_operation_end"`
	BusinessName       string    `json:"business_name"`
}
type ChangePhonedBody struct {
	NewPhone string `json:"new_phone"`
	OldPhone string `json:"old_phone"`
}
type PhoneLoginBody struct {
	BusinessPhone string `json:"business_phone"`
	Password      string `json:"password"`
}
