# go-bcwe3

## Cài đặt solc

#### **Trên macOS (dùng Homebrew)**

```bash
brew update
brew install solidity
```

#### **Trên Ubuntu / Debian**

```bash
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt update
sudo apt install solc
```

#### **Hoặc cài từ source hoặc binary chính thức**

Tải binary từ: https://github.com/ethereum/solidity/releases

Ví dụ:

```bash
wget https://github.com/ethereum/solidity/releases/download/v0.8.24/solc-static-linux
chmod +x solc-static-linux
sudo mv solc-static-linux /usr/local/bin/solc
```

Sau đó kiểm tra lại:

```bash
solc --version
```

## Cài đặt abigen

go install github.com/ethereum/go-ethereum/cmd/abigen

abigen --version

# example

abigen --abi=ERC721.abi --pkg=abi --type=ERC721 --out=./abi/erc721.go
