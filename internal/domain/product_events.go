package domain

type ProductCreatedEvent struct {
	Product *Product
}

func (p *ProductCreatedEvent) IsNotification() {

}
