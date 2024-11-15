package exceltypes

type ServiceFieldType uint8

const (
	Single ServiceFieldType = iota
	OneDArray
	TwoDArray
	OneDMap
	TwoDMap
	TwoDArrayMap
)

type ServiceFieldMode uint8

const (
	Input ServiceFieldMode = iota
	Output
)

type ServiceFieldOrient uint8

const (
	Horizontal ServiceFieldOrient = iota
	Vertical
)
