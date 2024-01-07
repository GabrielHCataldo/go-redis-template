package option

type SetMode string

const (
	// SetModeDefault sets the value regardless of whether it already exists or not
	SetModeDefault SetMode = ""
	// SetModeNx only set the key if it does not already exist.
	SetModeNx SetMode = "NX"
	// SetModeXx Only set the key if it already exists.
	SetModeXx SetMode = "XX"
)

func (s SetMode) String() string {
	return string(s)
}
