package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const key = "example_key"

type smartContract struct{}

func (sc smartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetArgs()
	fn := string(args[1])
	if len(args) > 2 {
		args = args[2:]
	} else {
		args = [][]byte{}
	}

	switch fn {
	case "set":
		if err := stub.PutState(key, args[0]); err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	case "get":
		data, err := stub.GetState(key)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(data)
	}
	return shim.Error("unknown method")
}

var SmartContract smartContract