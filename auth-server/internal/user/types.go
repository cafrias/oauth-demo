package user

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (f *loginForm) Validate() map[string]string {
	errors := make(map[string]string)

	if len(f.Email) == 0 {
		errors["email"] = "Email is required"
	}
	// TODO: email must be valid email address

	if len(f.Password) == 0 {
		errors["password"] = "Password is required"
	}
	// TODO: validate for strong password

	return errors
}
