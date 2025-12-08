package utils

// Paginate 对任何 slice 进行分页（传入数据、页数、页大小）
func Paginate[T any](data []T, page, pageSize int) []T {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	total := len(data)
	start := (page - 1) * pageSize
	if start >= total {
		return []T{}
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return data[start:end]
}
