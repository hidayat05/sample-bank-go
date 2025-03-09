

### preparation setup

- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


### Sample usage of RPC service

```
grpcurl -plaintext -import-path ./proto -proto bank.proto  -d '{"account_no":""}' '[::1]:PORT' bank.BankService/GetBalance
grpcurl -plaintext -import-path ./proto -proto bank.proto  -d '{"from_account_number":"", "to_account_number":"","amount":0 }' '[::1]:PORT' bank.BankService/TransferFunds
```
