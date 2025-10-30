package dto

type IPatchDTO interface {
	// ToPatchMap Converte o DTO de ponteiros em um mapa
	// apenas com os campos que não são nil.
	ToPatchMap() map[string]interface{}
}
