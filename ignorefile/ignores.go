package ignorefile

var defaultIgnores = []string{
	"**/.git/",
	"**/.gitignore",
	"**/.gitattributes",
	"**/.gitkeep",
	"**/node_modules/",
	"**/*cache*/",
	"**/*vscode*/",
	"**/go/",
	"**/.bash_history",
}
