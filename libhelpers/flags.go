package libhelpers

type Flags struct {
	Color         bool
	ContentLength string
	Cookie        string
	Header        string
	Method        string
	PostData      string
	Status        []int
	Threads       int
	TimeDelay     string
	Url           string
	Verbose       bool
	Wordlist      string
}

func NewFlags() Flags {
	var flags Flags
	flags.Color = true
	flags.ContentLength = "-1"
	flags.Cookie = ""
	flags.Header = ""
	flags.Method = ""
	flags.PostData = ""
	flags.Status = []int{}
	flags.Threads = 0
	flags.TimeDelay = ""
	flags.Url = ""
	flags.Verbose = false
	flags.Wordlist = ""

	return flags
}
