

<div align="center">



   # ⚡ Blazingly FAST Line Counter
   
   ![GitHub License](https://img.shields.io/github/license/Moderrek/lines)
   ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Moderrek/lines/total)


</div>


Fast command-line [concurrent](https://en.wikipedia.org/wiki/Concurrent_computing) **non-blank** line counter implemented in [GO](https://go.dev/) using [lightweight execution threads](https://go.dev/tour/concurrency/1).

## ⚙️ Usage

```shell
lines            # Prints file with the most lines at current directory
lines --dir      # Path to the analysis folder
lines --top N    # Prints the top N files
lines --hidden   # Allow to analyze hidden files & dirs
lines --version  # Prints installed version
lines --help     # Prints help
lines --no-color # Disables colored standard output
```

### 📈 Example output

```bat
lines --dir C:\Users\Moderr\dev --top 5
```

```out
Analyzing.. C:\Users\Moderr\dev

.java | Lines of code: 24409
.json | Lines of code: 8828
.yaml | Lines of code: 4980
.tsx | Lines of code: 4357
.yml | Lines of code: 1122

Time taken: 27.157ms to analyze 79 635 files
```

## 📸 Screenshots

![Example Usage](/images/ss.png)

## 🖥️ Quick Start

Requires

- Installed [Git](https://www.git-scm.com/downloads)
- Installed [GO](https://go.dev/doc/install)

Steps
1. Clone repository
   ```shell
   git clone https://github.com/Moderrek/lines
   ```
2. Run
   ```shell
   go run main.go
   ```

## © License

```license
MIT License

Copyright (c) 2024 Tymon Woźniak

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
