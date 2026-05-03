package domain

// "go-ca/internal/domain/product"

type Item struct {
	BaseEntity[uint]
	ProductID   uint     `gorm:"not null;index" json:"product_id"`
	Name        string   `gorm:"type:varchar(255);not null" json:"name"`
	Description *string  `gorm:"type:varchar(255)" json:"description"`
	Product     *Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"product,omitempty"`
}

func NewItem(name string, description *string) Item {
	return Item{
		Name:        name,
		Description: description,
	}
}

func (i *Item) SetProductId(value uint) {
	i.ProductID = value
}

func (i *Item) SetName(value string) {
	i.Name = value
}

func (i *Item) SetDescription(value *string) {
	i.Description = value
}
