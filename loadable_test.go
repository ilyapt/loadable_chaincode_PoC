package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
)

const SenderIdentity = "0a0a61746f6d797a654d535012d7062d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943536a434341664367417749424167495241496b514e37444f456b6836686f52425057633157495577436759494b6f5a497a6a304541774977675963780a437a414a42674e5642415954416c56544d524d77455159445651514945777044595778705a6d3979626d6c684d525977464159445651514845773154595734670a526e4a68626d4e7063324e764d534d77495159445651514b45787068644739746558706c4c6e56686443356b624851755958527662586c365a53356a6144456d0a4d4351474131554541784d64593245755958527662586c365a533531595851755a4778304c6d463062323135656d5575593267774868634e4d6a41784d44457a0a4d4467314e6a41775768634e4d7a41784d4445784d4467314e6a4177576a42324d517377435159445651514745774a56557a45544d4245474131554543424d4b0a5132467361575a76636d3570595445574d4251474131554542784d4e5532467549455a795957356a61584e6a627a45504d4130474131554543784d47593278700a5a5735304d536b774a7759445651514444434256633256794d554268644739746558706c4c6e56686443356b624851755958527662586c365a53356a6144425a0a4d424d4742797147534d34394167454743437147534d3439417745484130494142427266315057484d51674d736e786263465a346f3579774b476e677830594e0a504b6270494335423761446f6a46747932576e4871416b5656723270697853502b4668497634434c634935633162473963365a375738616a5454424c4d4134470a41315564447745422f775145417749486744414d42674e5648524d4241663845416a41414d437347413155644977516b4d434b4149464b2f5335356c6f4865700a6137384441363173364e6f7433727a4367436f435356386f71462b37585172344d416f4743437147534d343942414d43413067414d4555434951436e6870476d0a58515664754b632b634266554d6b31494a6835354444726b3335436d436c4d657041533353674967596b634d6e5a6b385a42727179796953544d6466526248740a5a32506837364e656d536b62345651706230553d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a"

func Test_Loadable_ChainCode(t *testing.T) {
	data1, err := ioutil.ReadFile("smart-contract1.so")
	if err != nil {
		t.Fatal(err)
	}
	stub := shimtest.NewMockStub("loadable", new(mainChainCode))

	identity, err := hex.DecodeString(SenderIdentity)
	assert.NoError(t, err)
	stub.Creator = identity

	txId := uuid.New().String()
	r1 := stub.MockInvoke(txId, [][]byte{
		[]byte("upload"),
		data1,
	})
	assert.Equal(t, http.StatusOK, int(r1.Status), r1.Message)

	hashed1 := sha3.Sum256([]byte(txId))
	smartContractAddr := base58.CheckEncode(hashed1[1:], hashed1[0])
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

	data2, err := ioutil.ReadFile("smart-contract2.so")
	if err != nil {
		t.Fatal(err)
	}

	txId2 := uuid.New().String()
	r4 := stub.MockInvoke(txId2, [][]byte{
		[]byte("upload"),
		data2,
	})
	assert.Equal(t, http.StatusOK, int(r4.Status), r4.Message)

	hashed := sha3.Sum256([]byte(txId2))
	r5 := stub.MockInvoke(uuid.New().String(), [][]byte{
		[]byte(base58.CheckEncode(hashed[1:], hashed[0])),
		[]byte("call"),
		[]byte(smartContractAddr),

	})
	assert.Equal(t, http.StatusOK, int(r5.Status), r5.Message)
	assert.Equal(t, randomData, r5.Payload)
}