package domain

import "time"

// User is a domain User
type User struct {
	id        int
	login     string
	password  string
	role      string
	token     string
	createdAt time.Time
	updatedAt time.Time
}

// NewUserData is a domain User
type NewUserData struct {
	ID        int
	Login     string
	Password  string
	Role      string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(data NewUserData) (User, error) {
	return User{
		id:        data.ID,
		login:     data.Login,
		password:  data.Password,
		role:      data.Role,
		token:     data.Token,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}, nil
}

func (u User) ID() int {
	return u.id
}
func (u User) Login() string {
	return u.login
}
func (u User) Password() string {
	return u.password
}
func (u User) Role() string {
	return u.role
}
func (u User) Token() string {
	return u.token
}
func (u User) CreratedAt() time.Time {
	return u.createdAt
}
func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

// // AddUser is ...
// type AddUser struct {
// 	Login    string `json:"login" db:"login" example:"cmd@cmd.ru"`
// 	Password string `json:"password" db:"password" example:"123456"`
// 	Role     string `json:"role" db:"role" example:"super"`
// }

// // LoginUser is
// type LoginUser struct {
// 	Login    string `json:"login" db:"login" example:"cmd@cmd.ru"`
// 	Password string `json:"password" db:"password" example:"123456"`
// }

// // NewUser is ...
// func NewUser() User {
// 	return User{}
// }

// // Validate проверяет правильность данных пользователя
// func (u *User) Validate() error {
// 	if u.Login == "" {
// 		return errors.New("login is required")
// 	}
// 	if !isValidEmail(u.Login) {
// 		return errors.New("invalid [login] format, must be email format")
// 	}
// 	if u.Password == "" {
// 		return errors.New("password is required")
// 	}
// 	if len(u.Password) < 6 {
// 		return errors.New("password must be at least 6 characters long")
// 	}
// 	if u.Role != "super" && u.Role != "regular" {
// 		return errors.New("invalid role")
// 	}
// 	return nil
// }

// // HashPassword хеширует пароль пользователя
// func (u *User) HashPassword() error {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	u.Password = string(hashedPassword)
// 	return nil
// }

// // isValidEmail проверяет формат электронной почты
// func isValidEmail(email string) bool {
// 	// Простая регулярка для проверки формата email
// 	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
// 	return re.MatchString(email)
// }
