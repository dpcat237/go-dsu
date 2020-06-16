package output

const (
	ModeProd Mode = iota + 1
	ModeDev
)

//Mode is mode of execution
type Mode uint16

//IsProduction checks if mode is production
func (md Mode) IsProduction() bool {
	return md == ModeProd
}
