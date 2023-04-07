package global

// PaperType [4 开 一页纸 闭] [4 上下] [4 中末] [4 AB]
type PaperType uint16

// AB default A
func (pt PaperType) AB() byte {
	switch pt & 0x0f {
	case 1:
		return 'A'
	case 2:
		return 'B'
	default:
		return 'A'
	}
}

func (pt PaperType) SetAB(x byte) PaperType {
	n := PaperType(0)
	switch x {
	case 'A':
		n = 1
	case 'B':
		n = 2
	}
	return pt | n
}

// MiddleFinal default 平时
func (pt PaperType) MiddleFinal() string {
	switch (pt & 0xf0) >> 4 {
	case 1:
		return "期中"
	case 2:
		return "期末"
	default:
		return "平时"
	}
}

func (pt PaperType) SetMiddleFinal(x string) PaperType {
	n := PaperType(0)
	switch x {
	case "中":
		n = 1 << 4
	case "末":
		n = 2 << 4
	}
	return pt | n
}

// FirstSecond default is 年度
func (pt PaperType) FirstSecond() string {
	switch (pt & 0x0f00) >> 8 {
	case 1:
		return "第1学期"
	case 2:
		return "第2学期"
	default:
		return "年度"
	}
}

func (pt PaperType) SetFirstSecond(x byte) PaperType {
	n := PaperType(0)
	switch x {
	case '1':
		n = 1 << 8
	case '2':
		n = 2 << 8
	}
	return pt | n
}

// OpenClose default 闭卷
func (pt PaperType) OpenClose() string {
	switch (pt & 0xf000) >> 12 {
	case 1:
		return "开卷"
	case 2:
		return "一页纸开卷"
	case 3:
		return "闭卷"
	default:
		return "闭卷"
	}
}

func (pt PaperType) SetOpenClose(x string) PaperType {
	n := PaperType(0)
	switch x {
	case "开卷":
		n = 1 << 12
	case "一页纸开卷":
		n = 2 << 12
	case "闭卷":
		n = 3 << 12
	}
	return pt | n
}
