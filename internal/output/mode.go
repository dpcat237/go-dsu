package output

const (
	//ModeProd execute command without debugging logs
	ModeProd Mode = iota + 1
	//ModeDev execute command with debugging logs
	ModeDev
)

//Mode is mode of execution
type Mode uint16

//IsProduction checks if mode is production
func (md Mode) IsProduction() bool {
	return md == ModeProd
}
