# Cyclone's Hashes.com API Escrow Tool

![image](https://i.imgur.com/0W4O7Ex.png)

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

### Compile from source code info:
- https://github.com/cyclone-github/scripts/blob/main/intro_to_go.txt
- *example of compiling on windows from src directory:* `go build  -ldflags="-s -w" -o hashes_tool.exe .`