package utils

func CalcPages(total, pageSize int64) int64 {
	var pages int64
	if total > pageSize {
		if total%pageSize > 0 {
			pages = (total / pageSize) + 1
		} else {
			pages = total / pageSize
		}
	} else {
		pages = 1
	}
	return pages
}
