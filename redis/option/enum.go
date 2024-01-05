package option

type SetMode string

const (
	SetModeDefault SetMode = ""
	SetModeNx      SetMode = "NX"
	SetModeXx      SetMode = "XX"
)

func (s SetMode) String() string {
	return string(s)
}
