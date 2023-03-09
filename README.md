![logo](img/banner.png)

# Fuzzagotchi

An automatic web fuzzer.

<br />

## Usage

![screenshot](img/screenshot.png)

### Automatic Fuzzing

This mode automatically fuzzes directories contains **`*.txt`**, **`*.php`**, **`.html`**, etc.

```sh
fuzzagotchi -u https://example.com -w wordlist.txt
```

### Specific Fuzzing

If you want to specify where to fuzz, you can put **"EGG"** keyword in URL, POST params, etc.

```sh
# Directories
fuzzagotchi -u https://example.com/EGG -w wordlist.txt
# Vhosts
fuzzagotchi -u https://example.com -H "Host: EGG.example.com" -w wordlist.txt
# Cookies
fuzzagotchi -u https://example.com -H "Cookie: key=EGG" -w wordlist.txt
```

### Using Built-in Wordlists

You can use **built-in wordlists** by specifying the special keywords as follow.

```sh
fuzzagotchi -u https://example.com/?id=EGG -w NUM_0_999
```

Below are the list of built-in wordlist.

```sh
# Alphabets (ALPHA_START_END)
ALPHA_A_Z
ALPHA_F_Q

# Numbers (NUM_START_END)
NUM_0_100
NUM_0000_9999
```

<br />

## Installation

The easiest way of installation is to install using **`go`** binary.  
Your system needs to have **`go`**.

### Install with Go

```sh
go install github.com/hideckies/fuzzagotchi@latest
```

### Clone This Repo & Build

Another way, you can clone this repository and build.

```sh
git clone https://github.com/hideckies/fuzzagotchi.git
cd fuzzagotchi
go get ; go build
```