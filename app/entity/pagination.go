package entity

type Pagination int

func (e Pagination) GetOffset(limit int) int {
	offset := (int(e) - 1) * limit
	if offset < 0 {
		return 0
	}
	return offset
}
