package subsplease

type SubsPlease interface {
	Latest() ([]Episode, error)
}
