# Hashes.com API Escrow Tool
```
 ######################################################################
#              Cyclone's Hashes.com API Escrow Tool v0.1.3             #
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
n.  Enter New API
r.  Remove API Key
c.  Clear Screen
q.  Quit
```
Tool written in Go for interacting with https://hashes.com escrow's API. Currently supports all known API calls from hashes.com.

Crypto / USD prices provided by: https://api.kraken.com

Inspiration from Plum's python3 script:
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
- Saves API key locally with AES encrypted key file

### Compile from source:
- If you want the latest features, compiling from source is the best option since the release version may run several revisions behind the source code.
- This assumes you have Go and Git installed
  - `git clone https://github.com/cyclone-github/hashes.com-escrow-tool.git`
  - `cd hashes.com-escrow-tool/src`
  - `go mod init escrow_tool`
  - `go mod tidy`
  - `go build -ldflags="-s -w" .`
- Compile from source code how-to:
  - https://github.com/cyclone-github/scripts/blob/main/intro_to_go.txt
