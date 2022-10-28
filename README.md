# Web3 - Random Transactions
Random transactions, to-and-fro the self and a peer.

## Instructions
Start by setting up a local, and preferably persistent, Ethereum blockchain.
```bash
npm ci
npx ganache -v -d --database.dbPath ./blk -i 344 -a 25
```

Choose from the output and add the relevant variables to a `.env` file in the root.
```env
# Required
WEB3_SELF_ADDRESS="<your-address>"
WEB3_SELF_KEY="<your-private-key>"

WEB3_PEER_ADDRESS="<another-address>"
WEB3_PEER_KEY="<their-private-key>"

# Optional
WEB3_PROVIDER="<ethereum-provider>"
```

Following are the instructions to run the application with different frameworks.

### [Web3.js](https://web3js.readthedocs.io/)
Web3.js is an Ethereum client for JavaScript.
```bash
npx ts-node index.ts
```

### [Web3.py](https://web3py.readthedocs.io/)
Web3.py is an Ethereum client for Python.
```bash
python -m venv env
. ./env/bin/activate
pip install -r requirements.txt
python app.py
```

### [Web3j](https://docs.web3j.io/)
Web3j is an Ethereum client for Java.
```bash
cd java/Web3Transactions
sdk env
mvn clean install
mvn package
java -jar target/Web3Transactions-1.0-SNAPSHOT-jar-with-dependencies.jar
```

### [Go-Ethereum](https://geth.ethereum.org/)
Go-Ethereum is an Ethereum client for Go.
```bash
go mod download
go run .
```

### [Rust-Web3](https://docs.rs/web3/latest/web3//)
Rust-Web3 is an Ethereum client for Rust.
```bash
cd rust
cargo run
```

# Made with ❤ by [Param](https://www.paramsid.com).