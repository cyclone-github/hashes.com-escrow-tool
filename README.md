[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username=cyclone-github&repo=hashes.com-escrow-tool&theme=gruvbox)](https://github.com/cyclone-github/hashes.com-escrow-tool/)

<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/cyclone-github/hashes.com-escrow-tool)](https://goreportcard.com/report/github.com/cyclone-github/hashes.com-escrow-tool) -->
[![GitHub issues](https://img.shields.io/github/issues/cyclone-github/hashes.com-escrow-tool.svg)](https://github.com/cyclone-github/hashes.com-escrow-tool/issues)
[![License](https://img.shields.io/github/license/cyclone-github/hashes.com-escrow-tool.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/cyclone-github/hashes.com-escrow-tool.svg)](https://github.com/cyclone-github/hashes.com-escrow-tool/releases)

# Hashes.com API Escrow Tool
```
 ######################################################################
#              Cyclone's Hashes.com API Escrow Tool v1.1.2             #
#           github.com/cyclone-github/hashes.com-escrow-tool           #
#            This tool requires an API key from hashes.com             #
#                   'Search Hashes' requires credits                   #
#                     See hashes.com for more info                     #
 ######################################################################

API key verified

Select an option:
1.  Upload Founds
2.  Upload History
3.  Download Left Lists
4.  Search Hashes
5.  Hash Identifier
6.  Wallet Balance
7.  Show Profit
8.  Withdrawal History
9.  Monitor Escrow WebSocket
n.  Enter New API
r.  Remove API Key
c.  Clear Screen
q.  Quit
```
Tool written in Go for interacting with https://hashes.com escrow's API. Currently supports all known API calls from hashes.com.

Inspiration from Plum's python3 CLI tool:
https://github.com/PlumLulz/hashes.com-cli
 
### Features:
- Upload Founds
  - ![image](https://i.imgur.com/GzRN3lE.png)
  - Select from file list
  - Paste custom file path
  - Paste hash:plaintext
- Show Upload History
- Download List Lists
- Search Hashes (requires credits from hashes.com)
- Hash Identifier
- Wallet Balance
- Show Profit
- Withdraw History
- Monitor Escrow WebSocket
- Saves API key locally with AES encrypted key file

### Usage:
- Run tool normally
  - `escrow_tool`
- Start tool in websocket escrow monitor mode (can be used to monitor escrow and/or pipe escrow job json into external program)
  - `escrow_tool -websocket`
  - `escrow_tool -websocket | jq`
  - `escrow_tool -websocket | escrow_cracking_bot`

### Install latest release:
```
go install github.com/cyclone-github/hashes.com-escrow-tool@latest
```
### Install from latest source code (bleeding edge):
```
go install github.com/cyclone-github/hashes.com-escrow-tool@main
```
### Compile from source:
- If you want the latest features, compiling from source is the best option since the release version may run several revisions behind the source code.
- This assumes you have Go and Git installed
  - `git clone https://github.com/cyclone-github/hashes.com-escrow-tool.git`  # clone repo
  - `cd hashes.com-escrow-tool`                                               # enter project directory
  - `go build -ldflags="-s -w" .`                              # compile binary in current directory
  - `go install -ldflags="-s -w" .`                            # compile binary and install to $GOPATH
- Compile from source code how-to:
  - https://github.com/cyclone-github/scripts/blob/main/intro_to_go.txt

### Changelog:
- https://github.com/cyclone-github/hashes.com-escrow-tool/blob/main/CHANGELOG.md

### Mentions:
- hashes.com offical docs: https://hashes.com/docs
- hashpwn forum: https://forum.hashpwn.net/post/76

### Antivirus False Positives:
- Several antivirus programs on VirusTotal incorrectly detect compiled Go binaries as a false positive. This issue primarily affects the Windows executable binary, but is not limited to it. If this concerns you, I recommend carefully reviewing the source code, then proceed to compile the binary yourself.
- Uploading your compiled binaries to https://virustotal.com and leaving an up-vote or a comment would be helpful as well.
