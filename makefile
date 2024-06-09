t:
	cd contracts && npx hardhat test
d:
	cd contracts && npx hardhat deploy --network sepolia
compile:
	cd cmd && go run main.go -mode compile
prove:
	cd cmd && go run main.go -mode prove -tx 8b805e46758497c6b32d0bf3cad3b3b435afeb0adb649857f24e424f75b79e46 -rpc https://eth-sepolia.g.alchemy.com/v2/69kYd2tLO_zx83HNzokjUqSojjCI3tw1
p:
	cd cmd && go run main.go -mode prove -tx 0x78a194fd6f445008a7e462c7b6b9e73d0017301ac1504f8c900f083f43d0d905 -rpc https://eth-sepolia.g.alchemy.com/v2/69kYd2tLO_zx83HNzokjUqSojjCI3tw1
p1:
	cd cmd && go run main.go -mode prove -tx 0x75fe114866f1633b537b8bde33b187eba6359207e7f7f4ad690366fd33e9f395 -rpc https://eth-sepolia.g.alchemy.com/v2/69kYd2tLO_zx83HNzokjUqSojjCI3tw1
p2:
	cd cmd && go run main.go -mode prove -tx 0x75fe114866f1633b537b8bde33b187eba6359207e7f7f4ad690366fd33e9f395

# My aacount initial transaction 0x78a194fd6f445008a7e462c7b6b9e73d0017301ac1504f8c900f083f43d0d905
# Refundee 0x58b529F9084D7eAA598EB3477Fe36064C5B7bbC1
# _callback Oxef1B4B164Fd3b7933bfaDb042373560e715Ec5D6


									 Oxef1B4B164Fd3b7933bfaDb042373560e715Ec5D6
https://sepolia.etherscan.io/address/0xef1B4B164Fd3b7933bfaDb042373560e715Ec5D6#code