package subsplease

type SubsPlease interface {
	Latests() ([]Episode, error)
}
