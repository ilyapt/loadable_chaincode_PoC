package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"golang.org/x/crypto/sha3"
	"loadable/common"
	"os"
	"path"
	"plugin"
)

const scKey = "smart-contracts"

type sc interface{
	Invoke(common.LoadableStubInterface, string, [][]byte) peer.Response
}

type LoadableStub struct {
	stub shim.ChaincodeStubInterface
	myAddress string
	sender string
}

func (ls *LoadableStub) GetStub() shim.ChaincodeStubInterface {
	return ls.stub
}

func (ls *LoadableStub) MyAddress() string {
	return ls.myAddress
}

func (ls *LoadableStub) Sender() string {
	return ls.sender
}

func (ls *LoadableStub) CallContract(addr, fn string, args [][]byte) peer.Response {
	return callSmartContract(2, addr, ls.stub, ls.myAddress, fn, args)
}

type mainChainCode struct {}

func (mcc *mainChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (mcc *mainChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetArgs()
	if string(args[0]) == "upload" {
		hashed := sha3.Sum256([]byte(stub.GetTxID()))
		key, err := stub.CreateCompositeKey(scKey, []string{ base58.CheckEncode(hashed[1:], hashed[0]) })
		if err != nil {
			return shim.Error(err.Error())
		}
		if err := stub.PutState(key, args[1]); err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	}

	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}
	var identity msp.SerializedIdentity
	if err := proto.Unmarshal(creator, &identity); err != nil {
		return shim.Error(err.Error())
	}
	b, _ := pem.Decode(identity.IdBytes)
	parsed, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	pk := parsed.PublicKey.(*ecdsa.PublicKey)
	ski := sha3.Sum256(elliptic.Marshal(pk.Curve, pk.X, pk.Y))

	return callSmartContract(1, string(args[0]), stub, base58.CheckEncode(ski[1:], ski[0]), string(args[1]), args[2:])
}

func callSmartContract(step int, key string, stub shim.ChaincodeStubInterface, sender string, fn string, args [][]byte) peer.Response {
	compositeKey, err := stub.CreateCompositeKey(scKey, []string{ key })
	if err != nil {
		return shim.Error(err.Error())
	}

	data, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(data) == 0 {
		return shim.Error("smart-contract doesn't exist")
	}
	file := path.Join("/tmp", key)
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
		fmt.Println("HERE", step)
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

	return sc.Invoke(&LoadableStub{
		stub: stub,
		sender: sender,
		myAddress: key,
	}, fn, args)
}
