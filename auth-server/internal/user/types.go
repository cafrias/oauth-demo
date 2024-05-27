package user

type User struct {
	ID    string
	Email string
}

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

type signupForm struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	PasswordConfirmation string `form:"password_confirmation"`
}

func (f *signupForm) Validate() map[string]string {
	errors := make(map[string]string)

	if len(f.Email) == 0 {
		errors["email"] = "Email is required"
	}
	// TODO: validate email

	if len(f.Password) == 0 {
		errors["password"] = "Password is required"
	}

	if f.Password != f.PasswordConfirmation {
		errors["password_confirmation"] = "Passwords do not match"
	}

	return errors
}
