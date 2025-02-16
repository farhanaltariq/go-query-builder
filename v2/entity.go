package new

type modes struct {
	mode string
}

const (
	MODE_SELECT = "SELECT"
	MODE_INSERT = "INSERT"
	MODE_UPDATE = "UPDATE"
	MODE_DELETE = "DELETE"
)

type direction struct {
	dir        string
	columnSort string
}

const (
	DIR_ASCENDING   = "ASC"
	DIR_DESCENDING = "DESC"
)
