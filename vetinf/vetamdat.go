package vetinf

type VetamdatRecord struct {
	ClientID string
	AnimalID string
	Index    int
	Data     string
}

type record struct {
	client [6]byte
	animal [6]byte
	index  [3]byte
	data   [(89 - 15)]byte
}