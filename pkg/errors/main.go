package apiErrors

type ErrorDescriptionT struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// Errors used to make http errors more understandable
var Errors = []ErrorDescriptionT{
	{0, "Invalid data"},
	{1, "User with this email already exists"},
	{2, "User don't exists"},
	{3, "Wrong password"},
	{4, "JWT token is invalid"},
	{5, "Group don't exist"},
	{6, "User doesn't own this m_group"},
}
