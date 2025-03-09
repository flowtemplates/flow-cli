package fs

type File struct {
	Name   string
	Path   string
	Source string
}

type Dir struct {
	Name  string
	Path  string
	Files []File
	Dirs  []Dir
}
