package simulation

// Writer is an interface to implement to write SimData to a persistent storage
type Writer interface {
	Write(s SimData)
}
