package main

import (
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
)

func Test_Loadable_ChainCode(t *testing.T) {
	data, err := ioutil.ReadFile("smart-contract.so")
	if err != nil {
		t.Fatal(err)
	}
	stub := shimtest.NewMockStub("loadable", new(mainChainCode))

	txId := uuid.New().String()
	r1 := stub.MockInvoke(txId, [][]byte{
		[]byte("upload"),
		data,
	})
	assert.Equal(t, http.StatusOK, int(r1.Status), r1.Message)

	smartContractAddr := txId
	randomData := make([]byte, 1024)
	rand.Read(randomData)

	r2 := stub.MockInvoke(uuid.New().String(), [][]byte{
		[]byte(smartContractAddr),
		[]byte("set"), randomData,
	})

	assert.Equal(t, http.StatusOK, int(r2.Status), r2.Message)
	r3 := stub.MockInvoke(uuid.New().String(), [][]byte{
		[]byte(smartContractAddr),
		[]byte("get"),
	})
	assert.Equal(t, http.StatusOK, int(r3.Status), r3.Message)
	assert.Equal(t, randomData, r3.Payload)
}