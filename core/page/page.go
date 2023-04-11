package page

type Page struct {
	PageNum   int64 `json:"pageNum"`
	PageSize  int64 `json:"pageSize"`
	TotalPage int64 `json:"totalPage"`
	Count     int64 `json:"count"`
}

func NewPage(pageNum, pageSize int64) *Page {
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return &Page{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func (p *Page) SetCount(count int64) {
	p.Count = count
	page := p.Count / p.PageSize
	if p.Count%p.PageSize == 0 {
		p.TotalPage = page
		return
	}
	p.TotalPage = page + 1
}
