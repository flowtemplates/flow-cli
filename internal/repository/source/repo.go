package source

type SourceRepo struct {
	baseDir string
}

func New(baseDir string) *SourceRepo {
	return &SourceRepo{
		baseDir: baseDir,
	}
}
