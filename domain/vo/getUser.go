package vo

type GetUser struct {
	Id       int
	Email    string
	Password string
	Role     []string
}
