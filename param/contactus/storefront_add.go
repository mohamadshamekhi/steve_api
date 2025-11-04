package contactusparam

type AddNewContactRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type AddNewContactResponse struct {
}
