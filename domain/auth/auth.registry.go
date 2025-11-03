package auth

func GetModels() []interface{} {
	return []interface{}{
		&Permission{},
		&Group{},
		&User{},
	}
}
