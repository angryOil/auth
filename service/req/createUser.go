package req

type CreateUser struct {
	Email    string
	Password string
	Role     []string
}
