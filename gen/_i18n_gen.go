package i18n

type (
	i18nFunc func(args ...interface{}) string

	errorType struct {
		NotFirstLaunch i18nFunc
		InvalidToken   i18nFunc
	}

	viewType struct {
		Index indexType
	}

	indexType struct {
		Hello i18nFunc
	}
)

var (
	Error = errorType{
		InvalidToken: func(args ...interface{}) string {
			return call("error.InvalidToken", args)
		},
		NotFirstLaunch: func(args ...interface{}) string {
			return call(
				"error.NotFirstLaunch", args)
		},
	}

	View = viewType{
		Index: Index,
	}

	Index = indexType{Hello: func(args ...interface{}) string {
		return call("view.index.hello", args)
	}}
)

func call(key string, args ...interface{}) string {
	return GetInstance().Get(key, args)
}
