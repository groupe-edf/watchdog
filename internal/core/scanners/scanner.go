package scanners

type Scanner interface {
	Name() string
	Kind() string
	Version() string
}

type BranchScanner interface {
	Scanner
	Scan()
}

type CommitScanner interface {
	Scanner
	Scan()
}

type FileScanner interface {
	Scanner
	Scan()
}

type TagScanner interface {
	Scanner
	Scan()
}
