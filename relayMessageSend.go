package main

import (
	"fmt"
	"time"
	"crypto/sha256"
	// "github.com/OpenBazaar/openbazaar-go/ipfs"
	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	// "github.com/golang/protobuf/ptypes/any"
	mh "github.com/multiformats/go-multihash"
  // "github.com/ipfs/go-ipfs/core"
  // libp2p "github.com/libp2p/go-libp2p-crypto"
  "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"

  "golang.org/x/crypto/nacl/box"
  "crypto/rand"
  "encoding/base64"

)

func main() {

	// 1.	Create the pb.Envelope protobuf object.
	// message Envelope {
	// 		Message message = 1;
	// 		bytes pubkey    = 2;
	// 		bytes signature = 3;
	// }


	// First we create the `pb.Chat` object
  chatMessage := "hey you bastard"
  subject := ""

  // The chat message contains a timestamp. This must be in the protobuf `Timestamp` format.
  timestamp, _ := ptypes.TimestampProto(time.Now())

  // The messageID is derived from the message data. In this case it's the hash of the message,
  // subject, and timestamp which is then multihash encoded.
  idBytes := sha256.Sum256([]byte(chatMessage + subject + ptypes.TimestampString(timestamp)))
  encoded, _ := mh.Encode(idBytes[:], mh.SHA2_256)
  msgID, _ := mh.Cast(encoded)

  chatPb := &pb.Chat{
    MessageId: msgID.B58String(),
    Subject:   subject,
    Message:   chatMessage,
    Timestamp: timestamp,
    Flag:      pb.Chat_MESSAGE,
  }

  // Now we wrap it in a `pb.Message` object
  payload, _ := ptypes.MarshalAny(chatPb)
  m := &pb.Message {
    MessageType: pb.Message_CHAT,
    Payload:     payload,
  }

	fmt.Printf("decoded: %s\n", m)

  // Now we wrap it in the envelop object
  // pubKeyBytes, _ := n.IpfsNode.PrivateKey.GetPublic().Bytes()

	// For now, I'm hard coding these pubKeyBytes from those I gathered from the server running my node.
	var pubKeyBytes = []byte{8, 1, 18, 32, 161, 156, 159, 166, 18, 55, 146, 203, 166, 196, 177, 109, 197, 156, 96, 232, 212, 21, 230, 139, 216, 18, 234, 89, 138, 44, 0, 220, 226, 141, 237, 245}


  // Use the protobuf serialize function to convert the object to a serialized byte array
  serializedMessage, _ := proto.Marshal(m)
  fmt.Printf("serializedMessage: %v\n", serializedMessage)
  // Sign the serializedMessage with the private key
  // signature, _ := n.IpfsNode.PrivateKey.Sign(serializedMessage)
  var signature = []byte{243, 91, 90, 62, 183, 105, 151, 42, 92, 131, 65, 227, 187, 53, 119, 63, 255, 135, 18, 123, 12, 26, 239, 231, 225, 166, 213, 196, 185, 139, 193, 108, 203, 210, 219, 12, 21, 194, 203, 240, 111, 235, 181, 187, 181, 191, 124, 249, 59, 186, 234, 190, 188, 9, 74, 74, 64, 98, 193, 35, 206, 88, 53, 10}


  // Create the envelope
  env := pb.Envelope{
    Message: m,
    Pubkey: pubKeyBytes,
    Signature: signature,
  }


  // ----- STEP 2: Encrypt the serialized envelope using the recipient's public key. For this you
  // will need to use an nacl library. NOTE for this you will need the recipient's public key.
  // We will have to create a server endpoint to get the pubkey. Technically I think the gateway
  // already has one but we may need to improve it for this purpose. The public key is also found
  // inside a listing so if you're looking at a listing you should already have it.

  // Serialize the envelope
  serializedEnvelope, _ := proto.Marshal(&env)

  // Get the public key
  // recipientPublicKey := getPublicKeyFromGatewayOrListing()
  var recipientPublicKeyB64 = "CAESIJBMYMhIa7B+rZQdfRpf92y4g2Izmv35C6tyuflhR8P8"
  

  recipientPublicKey, err := base64.StdEncoding.DecodeString(recipientPublicKeyB64)
  if err != nil {
    fmt.Printf("That didn't decode properly...")
  }

  // Generate an ephemeral key pair
  ephemPub, ephemPriv, _ := box.GenerateKey(rand.Reader)

  // Extra thing Chris provided in Slack
  pub, _ := crypto.UnmarshalPublicKey(recipientPublicKey)
  edPubkey := pub.(*crypto.Ed25519PublicKey)

  cPubkey, _ := edPubkey.ToCurve25519()

  // Convert recipient's key into curve25519
  // pk, _ := recipientPublicKey.ToCurve25519()

  // Encrypt with nacl

  // Nonce must be a random 24 bytes
  var nonce [24]byte
  n := make([]byte, 24)
  rand.Read(n)
  for i := 0; i < 24; i++ {
    nonce[i] = n[i]
  }

  var ciphertext2 []byte
  // Encrypt
  ciphertext := box.Seal(ciphertext2, serializedEnvelope, &nonce, cPubkey, ephemPriv)

  // Prepend the ephemeral public key to the ciphertext
  ciphertext = append(ephemPub[:], ciphertext...)

  // Prepend nonce to the ephemPubkey+ciphertext
  ciphertext = append(nonce[:], ciphertext...)

  // Base64 encode
  encodedCipherText := base64.StdEncoding.EncodeToString(ciphertext)

  fmt.Printf("encodedCipherText: %s\n", encodedCipherText)


}