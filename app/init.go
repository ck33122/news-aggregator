package app

var (
	appName string
)

// Init initialized instra'structure: config, log, etc.
// If something went wrong prints error message to stderr, and exists.
func Init(_appName string) {
	appName = _appName
	Must(InitConfig())
	Must(InitLog())
	Must(InitDB())
}

// Destroy prepares application to end, closes some important fd's.
func Destroy() {
	DestroyLog()
	DestroyDB()
}

func GetName() string {
	return appName
}
