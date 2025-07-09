package animal

type Animal struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Age         int    `db:"age" json:"age"`
	Description string `db:"description" json:"description"`
}

type AninalCreateRequest struct {
	Name        string
	Age         int
	Description string
}

type AnimalUpdateRequest struct {
	Name        string
	Age         int
	Description string
}
