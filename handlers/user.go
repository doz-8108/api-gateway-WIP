package handlers

import (
	"fmt"
	"regexp"

	"github.com/doz-8108/api-gateway/config"
	"github.com/doz-8108/api-gateway/storage"
	"github.com/doz-8108/api-gateway/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserHandlers struct {
		Storage *storage.UserStorage
		EnvVars config.EnvVars
		Utils   utils.Utils
	}
	UserSignUpReqBody struct {
		UserName string `json:"user_name" validate:"required,min=6,max=25,username"`
		Gender   string `json:"gender" validate:"oneof=M F ''"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=20,password"`
	}
	UserVerifyReqQueries struct {
		RedirectUrl string `json:"redirect_url"`
		Token       string `json:"token"`
	}
	UserSignInReqBody struct {
		Email    string
		Password string
	}
)

func NewUserHandlers(storage *storage.UserStorage, envVars config.EnvVars, utls utils.Utils) *UserHandlers {
	return &UserHandlers{
		Storage: storage,
		EnvVars: envVars,
		Utils:   utls,
	}
}

func (u *UserHandlers) SignInUser(f fiber.Ctx) error {
	userSignInReqBody := new(UserSignInReqBody)
	if err := f.Bind().Body(userSignInReqBody); err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	user := u.Storage.FindUserByEmail(userSignInReqBody.Email)
	if user == nil {
		u.Utils.CatchError(nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password"))
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userSignInReqBody.Password))
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password"))
	}

	token, err := u.Utils.CreateToken(user.Id, string(f.Request().Host()))
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return f.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (u *UserHandlers) VerifyUserEmail(f fiber.Ctx) error {
	userVerifyReqQueries := new(UserVerifyReqQueries)
	if err := f.Bind().Query(userVerifyReqQueries); err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusBadRequest, "Invalid query params"))
	}

	email, err := u.Storage.FindUnverifiedEmailByToken(userVerifyReqQueries.Token)
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusGone, "The registration link is expired or invalidated"))
	}

	user := u.Storage.FindUnverifiedUserByEmail(email)
	if user == nil {
		u.Utils.CatchError(nil, fiber.NewError(fiber.StatusGone, "The registration link is expired or invalidated"))
	}

	err = u.Storage.CreateUser(storage.UserDataPending{
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, "Internal server error"))
	}

	return f.Redirect().Status(fiber.StatusSeeOther).To(userVerifyReqQueries.RedirectUrl)
}

func (u *UserHandlers) SignUpUser(f fiber.Ctx) error {
	userSignUpReqBody := new(UserSignUpReqBody)

	if err := f.Bind().Body(userSignUpReqBody); err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	if user := u.Storage.FindUnverifiedUserByEmail(userSignUpReqBody.Email); user != nil {
		u.Utils.CatchError(fmt.Sprintf("Email already sent: %s", user.Email), fiber.NewError(fiber.StatusBadRequest, "Please activate your account via email"))
	}

	// field validation
	validate := validator.New()
	containsAcceptedCharsOnly := regexp.MustCompile("^[a-zA-Z0-9!@#~$%^&*()_+|<>-?:{}]+$").MatchString
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
		hasLowercase := regexp.MustCompile(`[a-z]`).MatchString
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString
		hasSymbol := regexp.MustCompile(`[!@#~$%^&*()_+|<>?:{}]`).MatchString
		return containsAcceptedCharsOnly(password) && hasLowercase(password) && hasUppercase(password) && hasNumber(password) && hasSymbol(password)
	})
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return containsAcceptedCharsOnly(fl.Field().String())
	})

	if err := validate.Struct(userSignUpReqBody); err != nil {
		errMessage := "Invalid request body"
		errs := err.(validator.ValidationErrors)
		if len(errs) > 0 {
			errMap := map[string](map[string]string){
				"UserName": {
					"min":      "Username should be at least 6 characters long",
					"max":      "Username should be at most 25 characters long",
					"username": "Username only accepts uppercase letters, lowercase letters, numbers, and symbols from (!@#~$%^&*()_+|<>-?:{})",
					"required": "Username is missing",
				},
				"Email": {
					"email":    "Invalid email address",
					"required": "Email is missing",
				},
				"Password": {
					"min":      "Password should be at least 8 characters long",
					"max":      "Password should be at most 20 characters long",
					"password": "Password must include at least one uppercase letter, one lowercase letter, one number, and one symbol from (!@#~$%^&*()_+|<>-?:{})",
					"required": "Password is missing",
				},
				"Gender": {
					"oneof": "Accepted values of gender: M or F (if specified)",
				},
			}
			tag := errs[0].Tag()
			if val, ok := errMap[errs[0].Field()][tag]; ok {
				errMessage = val
			}
			return fiber.NewError(fiber.StatusUnprocessableEntity, errMessage)
		}
	}

	// check email's existence
	if user := u.Storage.FindUserByEmail(userSignUpReqBody.Email); user != nil {
		u.Utils.CatchError(fmt.Sprintf("Duplicate email: %s", user.Email), fiber.NewError(fiber.StatusConflict, "Email already exists"))
	}

	// store user data
	bytes, err := bcrypt.GenerateFromPassword([]byte(userSignUpReqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	err = u.Storage.CreateUnverifiedUser(storage.UserDataPending{
		UserName: userSignUpReqBody.UserName,
		Email:    userSignUpReqBody.Email,
		Password: string(bytes),
	})
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	token, err := gonanoid.New()
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	err = u.Storage.CreateVerificationToken(userSignUpReqBody.Email, token)
	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = u.Utils.SendSignUpEmail(
		string(f.Request().Host()),
		userSignUpReqBody.UserName,
		userSignUpReqBody.Email,
		token,
		string(f.Request().Header.Referer()),
	)

	if err != nil {
		u.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return f.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "An verification email has been sent to your mailbox",
	})
}
