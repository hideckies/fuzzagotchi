package libgotchi

import "github.com/hideckies/fuzzagotchi/libhelpers"

type Conf struct {
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

func NewConf(flags libhelpers.Flags) Conf {
	var conf Conf
	conf.Color = flags.Color
	conf.ContentLength = flags.ContentLength
	conf.Cookie = flags.Cookie
	conf.Header = flags.Header
	conf.Method = flags.Method
	conf.PostData = flags.PostData
	conf.Status = flags.Status
	conf.Threads = flags.Threads
	conf.TimeDelay = flags.TimeDelay
	conf.Url = flags.Url
	conf.Verbose = flags.Verbose
	conf.Wordlist = flags.Wordlist

	return conf
}