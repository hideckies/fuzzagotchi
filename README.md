![logo](https://raw.githubusercontent.com/hideckies/fuzzagotchi/main/img/logo.png?token=GHSAT0AAAAAABWMAINRLCHXU6PTWUWLPGV6YXZFOXQ)

# Fuzzagotchi

A fuzzing tool written in Go.

<br />

## Prerequisites

1. **Seclists**

    Fuzzagotchi uses [seclists](https://github.com/danielmiessler/SecLists) as default wordlist.

    ```sh
    sudo apt install seclists
    ```

<br />

## Usage

```sh
fuzzagotchi -w wordlist.txt -u https://fuzzagotchi.xxx/
```

<br />

## Compile

```sh
go build
```