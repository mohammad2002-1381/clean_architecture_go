package domain

type Product struct {
	BaseEntity[uint]
	Name        string  `gorm:"column:name;not null"`
	Description *string `gorm:"column:description"`
	Price       float64 `gorm:"column:price;not null"`
	UserID      uint    `gorm:"column:user_id;not null"`
	Items       []*Item `gorm:"foreignKey:ProductID"`
	User        *User   `gorm:"foreignKey:UserID"`
}

func NewProduct(name string, description *string, price float64, userID uint, userRegisterdEvent Notification) Product {
	p := Product{
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
		Items:       make([]*Item, 0),
	}

	p.AddNotification(userRegisterdEvent)

	return p
}

func (p *Product) SetName(value string) {
	p.Name = value
}

func (p *Product) SetDescription(value *string) {
	p.Description = value
}

func (p *Product) SetPrice(value float64) {
	p.Price = value
}

func (p *Product) SetUserID(userID uint) {
	p.UserID = userID
}

func (p *Product) SetUser(user *User) {
	p.User = user
	if user != nil {
		p.UserID = user.ID
	}
}

func (p *Product) SetItems(items []*Item) {
	p.Items = items
}

func (p *Product) AddItem(item *Item) {
	p.Items = append(p.Items, item)
}
