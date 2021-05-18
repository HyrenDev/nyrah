package providers

type IProvider interface {
	Prepare()

	Provide()
}