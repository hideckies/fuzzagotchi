package libhelpers

type Flags struct {
	Header    string
	PostData  string
	Threads   int8
	TimeDelay string
	Url       string
	Verbose   bool
	Wordlist  string
}

func NewFlags() Flags {
	var flags Flags
	flags.Header = ""
	flags.PostData = ""
	flags.Threads = 0
	flags.TimeDelay = ""
	flags.Url = ""
	flags.Verbose = false
	flags.Wordlist = ""

	return flags
}
