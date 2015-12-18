package runner

type Runner interface {
	Run() error
	Setup() error
}
