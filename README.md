![logo](img/logo.png)

# Fuzzagotchi

A fuzzing tool written in Go. It helps to discover directories, files, parameters, etc. in web applications.  

<br />

## Status

**This is under development. So only some features work for now.**

<br />

## Prerequisites

1. **Seclists**

    Fuzzagotchi uses [seclists](https://github.com/danielmiessler/SecLists) as default wordlist.

    ```sh
    sudo apt install seclists
    ```

<br />

## Example

Fuzzagotchi uses "EGG" keyword for fuzzing.  
As you may have noticed, it imitates **FFUF**ツ

```sh
fuzzagotchi -u https://example.com/EGG -w wordlist.txt
```

<br />

## Usage

```
Usage:
  fuzzagotchi [flags]

Examples:

  [Content Discovery]
        fuzzagotchi -u https://example.com/EGG -w wordlist.txt
        fuzzagotchi -u https://example.com/EGG -w wordlist.txt --status-codes 200,301
        fuzzagotchi -u https://example.com/EGG -w wordlist.txt --content-length 175
        fuzzagotchi -u https://example.com/EGG -w wordlist.txt -H "Cookie: isGotchi=true"

        fuzzagotchi -u https://example.com/EGG.php -w wordlist.txt
        fuzzagotchi -u https://example.com/?q=EGG -w wordlist.txt

  [Brute Force POST Data] *Unser development so unavailable currently.
        fuzzagotchi -u https://example.com/login -w passwords.txt --post-data "username=admin&password=EGG"
        fuzzagotchi -u https://example.com/login -w passwords.txt --post-data "{username:admin, password: EGG}"

  [Subdomain Scan] *Under development so unavailable currently.
        fuzzagotchi -u https://EGG.example.com -w wordlist.txt


Flags:
      --content-length int   Display the specific content length only (default -1)
      --delay string         Time delay per requests e.g. 500ms. Or random delay e.g. 500ms-700ms (default "100-200")
  -H, --header string        Custom header e.g. "Authorization: Bearer <token>; Cookie: key=value"
  -h, --help                 help for fuzzagotchi
      --post-data string     POST request with data e.g. "username=admin&password=EGG"
  -s, --status-codes ints    Display the specific status codes only (default [200,204,301,302,307,401,403])
  -t, --threads int8         Number of concurrent threads. (default 10)
  -u, --url string           Target URL (required)
  -v, --verbose              Verbose mode
      --version              version for fuzzagotchi
  -w, --wordlist string      Wordlist for fuzzing (default "/usr/share/seclists/Discovery/Web-Content/common.txt")
```

<br />

## Compile

```sh
go get ; go build
```