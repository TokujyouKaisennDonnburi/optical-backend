package option

type Option struct {
	Id   int
	Name string
	deprecated bool
}

// func NewOption(name string) (*Option, error) {
// 	id, err := uuid.NewV7()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if name == "" {
// 		return nil, errors.New("Option `Name` is empty")
// 	}
// 	return &Option{
// 		Id:   id,
// 		Name: name,
// 	}, nil
// }
