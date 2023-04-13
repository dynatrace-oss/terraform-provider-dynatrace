package httpcache

type Finisher interface {
	Finish(v any) error
}
