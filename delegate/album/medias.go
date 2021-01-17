package album

//func ListMedias(filter *UpdateMediaFilter) (*MediaPage, error) {
//	return Repository.FindMedias(filter)
//}
//
//func InsertMedias(medias []Media) error {
//	return Repository.InsertMedias(medias)
//}
//
//func FindSignatures(signatures []MediaSignature) ([]MediaSignature, error) {
//	return Repository.FindExistingSignatures(signatures)
//}

type MoveMediaOperator interface {
	Copy(source, dest MediaLocation) error
	Remove(target MediaLocation) error

	UpdateStatus(done, total int) error
	Continue() bool
}

// RelocateMovedMedias moved data and update index with ability to interrupt
func RelocateMovedMedias(operator *MoveMediaOperator) (int, error) {
	return 0, nil
}

