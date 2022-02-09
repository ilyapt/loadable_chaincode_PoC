package common

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type LoadableStubInterface interface {
	GetStub() shim.ChaincodeStubInterface
	MyAddress() string
	Sender() string
	CallContract(string, string, [][]byte) peer.Response
}
