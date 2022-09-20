package libhelpers

type Flags struct {
	Color           bool
	ContentLength   string
	NoContentLength string
	Cookie          string
	Header          string
	Method          string
	PostData        string
	Status          []int
	NoStatus        []int
	Threads         int
	TimeDelay       string
	Timeout         int
	Url             string
	Verbose         bool
	Wordlist        string
}

func NewFlags() Flags {
	var flags Flags
	flags.Color = true
	flags.ContentLength = "-1"
	flags.NoContentLength = "-1"
	flags.Cookie = ""
	flags.Header = ""
	flags.Method = ""
	flags.PostData = ""
	flags.Status = []int{}
	flags.NoStatus = []int{}
	flags.Threads = 0
	flags.TimeDelay = ""
	flags.Timeout = 0
	flags.Url = ""
	flags.Verbose = false
	flags.Wordlist = ""

	return flags
}
