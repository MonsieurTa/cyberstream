package transcoder

type Reader interface{}
type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Transcode(tp *TranscoderParams) error
}
