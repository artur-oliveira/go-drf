package auth

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Group{},
		&Permission{},
	}
}
