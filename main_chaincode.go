package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"os"
	"path"
	"plugin"
)

const scKey = "smart-contracts"

type sc interface{
	Invoke(stub shim.ChaincodeStubInterface) peer.Response
}

type mainChainCode struct {}

func (mcc *mainChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (mcc *mainChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetArgs()
	if string(args[0]) == "upload" {
		key, err := stub.CreateCompositeKey(scKey, []string{ stub.GetTxID() })
		if err != nil {
			return shim.Error(err.Error())
		}
		if err := stub.PutState(key, args[1]); err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	}

	key, err := stub.CreateCompositeKey(scKey, []string{ string(args[0]) })
	if err != nil {
		return shim.Error(err.Error())
	}
	data, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(data) == 0 {
		return shim.Error("smart-contract doesn't exist")
	}
	file := path.Join("/tmp", string(args[0]))
	f, err := os.Create(file)
	if err != nil {
		return shim.Error(err.Error())
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return shim.Error(err.Error())
	}
	if err := f.Close(); err != nil {
		return shim.Error(err.Error())
	}
	plug, err := plugin.Open(file)
	if err != nil {
		return shim.Error(err.Error())
	}
	smartContract, err := plug.Lookup("SmartContract")
	if err != nil {
		return shim.Error(err.Error())
	}
	sc, ok := smartContract.(sc)
	if !ok {
		return shim.Error("incompatible smart contract")
	}
	return sc.Invoke(stub)
}
