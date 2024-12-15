package job

type HandleableJob interface {
	Handle() chan bool
}
