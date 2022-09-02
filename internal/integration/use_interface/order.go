package use_interface

//go:generate egen dao -src . -dst ./order_dao.go -type Order -import github.com/gotomicro/egen/internal/integration

type Order struct {
	// @PrimaryKey true
	UserId  uint32
	OrderId uint32
	Price   int
}
