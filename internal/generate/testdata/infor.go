package code

type User struct {
	LoginTime string
	FirstName string
	UserId    uint32
	LastName  string
}

type Order struct {
	OrderTime string
	OrderId   uint32
	UserId    uint32
}
