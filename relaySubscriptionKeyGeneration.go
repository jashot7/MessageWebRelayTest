package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/multiformats/go-multihash"
	// "github.com/libp2p/go-libp2p-crypto"
	"gx/ipfs/QmWFAMPqsEyUX7gDUsRVmMWz59FxSpJ1b2v6bJ1yYzo7jY/go-base58-fast/base58"
)

/**
	Input: QmaSAmPPynrWfz1R8XvRm1GX6ghzPze6XSZCov6fWWUzSg
	Output: QmaGLQjHHdeZ3wKtKqHS9etMwSUDnckHnAYS6eqvAgp2Hf
**/
func main() {

	// This is a test peerID from our docs on this.
	// peerIDString := "QmaSAmPPynrWfz1R8XvRm1GX6ghzPze6XSZCov6fWWUzSg"

	// This is my BTC4002 store
	peerIDString := "QmdonEatSBPBw35MUb6vUjkxboirX1SVFduPDmMjZw37MD"
	

	// Step 1: Convert the PeerID string to a Multihash object
	peerIDMultihash, _ := multihash.FromB58String(peerIDString)
	fmt.Printf("peerIDMultihash: %s\n", peerIDMultihash.B58String())

	// Step 2: Decode the Multihash to extract the digest
	decoded, _ := multihash.Decode(peerIDMultihash)
	fmt.Printf("decoded: %s\n", decoded)
	digest := decoded.Digest
	fmt.Printf("digest: %v\n", digest)

	// Step 3: Grab the first 8 bytes of the digest byte array
	prefix := digest[:8]
	fmt.Printf("prefix: %v\n", prefix)

	// Step 4: Bit shift the prefix to the right by 48 places.

	// If the prefix is:
	// 11111111 10101010 01010101 00001111 11110000 11001100 00110011 10110010

	// After the shift you should have:
	// 00000000 00000000 00000000 00000000 00000000 00000000 11111111 10101010

	// In Go this is done by first converting the byte array to a uint64
	prefix64 := binary.BigEndian.Uint64(prefix)
	fmt.Printf("prefix64: %v\n", prefix64)

	// Then shifting
	shiftedPrefix64 := prefix64>>uint(48)

	// Then converting back to a byte array
	shiftedBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(shiftedBytes, shiftedPrefix64)
	fmt.Printf("shiftedBytes: %v\n", shiftedBytes)

	// Step 5: Hash the shifted prefix with SHA256
	hashedShiftedPrefix := sha256.Sum256(shiftedBytes)
	fmt.Printf("hashedShiftedPrefix: %+v\n", hashedShiftedPrefix)
	fmt.Printf("hashedShiftedPrefixHEX: %X\n", hashedShiftedPrefix)

	// Step 6: Re-encode as a multihash to get your SubscriptionKey
	SubcriptionKey, _ := multihash.Encode(hashedShiftedPrefix[:], multihash.SHA2_256)
	// fmt.Println(SubcriptionKey.B58String()) // QmPZ9uGyBAJXjE7GXRHE9NG8CcuKei2z5PczmmuQnStaMu

	base58SubScriptionKey := base58.Encode(SubcriptionKey)
	fmt.Println(base58SubScriptionKey)


	// mh, err := multihash.FromB58String(SubcriptionKey)
	// if err != nil {
	// 	t.Error(err)
	// }
	// for _, test := range tests {
	// 	key := CreatePointerKey(mh, test.PrefixLen)
	// 	if key.B58String() != test.ExpectedValue {
	// 		t.Error("Returned incorrect pointer key")
	// 	}
	// }

}