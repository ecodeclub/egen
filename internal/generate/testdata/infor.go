package code

type User struct {
	LoginTime string
	FirstName string
	UserId    uint32
	LastName  string
	Password  []byte
}

type Order struct {
	HasBuy    bool
	Price     float64
	OrderTime string
	OrderId   uint32
	UserId    uint32
	Seller    *int
}

type OrderDAO interface {
	// @select select * from order where name in #name and status=#status
	Hello(name [5]string, status int) (int64, error)
}
