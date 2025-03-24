package cfg

type PSQLconfig struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Database      string
	MigrationPath string
}
type HTTPconfig struct {
	Host string
	Port string
}

type Config struct {
	DebugFilePath string
	InfoFilePath  string
	ConsoleOutput string
}
