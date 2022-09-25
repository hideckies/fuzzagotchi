package libgotchi

import "github.com/hideckies/fuzzagotchi/libhelpers"

type Conf struct {
	libhelpers.Flags
}

func NewConf(flags libhelpers.Flags) Conf {
	var conf Conf
	conf.Color = flags.Color
	conf.ContentLength = flags.ContentLength
	conf.NoContentLength = flags.NoContentLength
	conf.Cookie = flags.Cookie
	conf.Header = flags.Header
	conf.Method = flags.Method
	conf.PostData = flags.PostData
	conf.Rate = flags.Rate
	conf.Status = flags.Status
	conf.Threads = flags.Threads
	conf.Timeout = flags.Timeout
	conf.Url = flags.Url
	conf.Verbose = flags.Verbose
	conf.Wordlist = flags.Wordlist

	return conf
}
