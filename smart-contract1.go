package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"loadable/common"
)

const key = "example_key1"

type smartContract struct{}

func (sc smartContract) Invoke(stub common.LoadableStubInterface, fn string, args [][]byte) peer.Response {
	fmt.Println("I AM", stub.MyAddress(), "SENDER", stub.Sender())

	switch fn {
	case "set":
		if err := stub.GetStub().PutState(key, args[0]); err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "get":
		data, err := stub.GetStub().GetState(key)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	case "call":
		return stub.CallContract(string(args[0]), "get", nil)
	}
	return shim.Error("unknown method")
}

var SmartContract smartContract