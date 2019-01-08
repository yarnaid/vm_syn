package vminstance

type vmController interface {
	CurrentAddr() cpuArch
}
