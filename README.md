![logo](https://raw.githubusercontent.com/hideckies/fuzzagotchi/main/img/logo.png?token=GHSAT0AAAAAABWMAINRLCHXU6PTWUWLPGV6YXZFOXQ)

# Fuzzagotchi

A fuzzing tool.

## Prerequisites

1. **Seclists**

    Fuzzagotchi uses [seclists](https://github.com/danielmiessler/SecLists) as default wordlist.

    ```sh
    sudo apt install seclists
    ```

## Usage

```sh
fuzzagotchi -w wordlist.txt -u https://fuzzagotchi.xxx/
```
