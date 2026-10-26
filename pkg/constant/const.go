package constant

const (
	SendOTPSubject = "notification.send_opt"
	// BirthDateLayout is the Gregorian date format accepted for user birth dates.
	BirthDateLayout  = "2006-01-02"
	PhoneNumberRegex = `^\+?[(]?\d{3}[)]?[-\s.]?\d{3}[-\s.]?\d{4,6}$`
)
