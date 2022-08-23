package libhelpers

type Flags struct {
	ContentLength int
	Header        string
	PostData      string
	StatusCodes   []int
	Threads       int8
	TimeDelay     string
	Url           string
	Verbose       bool
	Wordlist      string
}

func NewFlags() Flags {
	var flags Flags
	flags.ContentLength = -1
	flags.Header = ""
	flags.PostData = ""
	flags.StatusCodes = []int{}
	flags.Threads = 0
	flags.TimeDelay = ""
	flags.Url = ""
	flags.Verbose = false
	flags.Wordlist = ""

	return flags
}
