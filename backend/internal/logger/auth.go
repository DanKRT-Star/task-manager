package logger

func AuthRegisterSuccess(userID uint, email string) {
	Info("auth_register_success", "user registered successfully",
		Field{Key: "user_id", Value: userID},
		Field{Key: "email", Value: email},
	)
}

func AuthRegisterEmailExists(email string) {
	Warn("auth_register_email_exists", "registration failed: email already registered",
		Field{Key: "email", Value: email},
	)
}

func AuthRegisterPasswordHashFailed(email string, err error) {
	Error("auth_register_password_hash_failed", "failed to hash password during registration",
		Field{Key: "email", Value: email},
		Field{Key: "error", Value: err},
	)
}

func AuthRegisterFailed(email string, err error) {
	Error("auth_register_failed", "failed to create user during registration",
		Field{Key: "email", Value: email},
		Field{Key: "error", Value: err},
	)
}

func AuthLoginSuccess(userID uint, email string) {
	Info("auth_login_success", "user logged in successfully",
		Field{Key: "user_id", Value: userID},
		Field{Key: "email", Value: email},
	)
}

func AuthLoginUserNotFound(email string) {
	Warn("auth_login_user_not_found", "login failed: invalid email or password",
		Field{Key: "email", Value: email},
	)
}

func AuthLoginInvalidPassword(userID uint, email string) {
	Warn("auth_login_invalid_password", "login failed: invalid email or password",
		Field{Key: "user_id", Value: userID},
		Field{Key: "email", Value: email},
	)
}

func AuthLoginTokenGenerationFailed(userID uint, err error) {
	Error("auth_login_token_generation_failed", "failed to generate jwt token",
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}

func AuthGetMeFailed(userID uint, err error) {
	Error("auth_get_me_failed", "failed to fetch current user",
		Field{Key: "user_id", Value: userID},
		Field{Key: "error", Value: err},
	)
}