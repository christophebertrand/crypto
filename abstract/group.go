package abstract

import (
	"crypto/cipher"
	"github.com/dedis/crypto/suites"
)

// XXX consider renaming Secret to Scalar?

/*
A Secret abstractly represents a secret value by which
a Point (group element) may be encrypted to produce another Point.
This is an exponent in DSA-style groups,
in which security is based on the Discrete Logarithm assumption,
and a scalar multiplier in elliptic curve groups.
*/
type SecretInterface interface {
	Marshaling

	// Equality test for two Secrets derived from the same Group
	Equal(s2 *Secret) bool

	// Set equal to another Secret a
	Set(a *Secret) *Secret

	// Set to a small integer value
	SetInt64(v int64) *Secret

	// Set to the additive identity (0)
	Zero() *Secret

	// Set to the modular sum of secrets a and b
	Add(a, b *Secret) *Secret

	// Set to the modular difference a - b
	Sub(a, b *Secret) *Secret

	// Set to the modular negation of secret a
	Neg(a *Secret) *Secret

	// Set to the multiplicative identity (1)
	One() *Secret

	// Set to the modular product of secrets a and b
	Mul(a, b *Secret) *Secret

	// Set to the modular division of secret a by secret b
	Div(a, b *Secret) *Secret

	// Set to the modular inverse of secret a
	Inv(a *Secret) *Secret

	// Set to a fresh random or pseudo-random secret
	Pick(rand cipher.Stream) *Secret

	// GetSuite returns the suite used to create that secret
	GetSuite() Suite

	// SetSuite sets the suite used for that point - only used in
	// setting up that secret
	SetSuite(s Suite)
}

/*
Secrets carry around their suite their using, so we can easily unmarshal
into that structure, whereas unmarshalling into an interface needs the
interface to be initialised first.
*/

type Secret struct {
	SecretInterface
}

/*
A Point abstractly represents an element of a public-key cryptographic Group.
For example,
this is a number modulo the prime P in a DSA-style Schnorr group,
or an x,y point on an elliptic curve.
A Point can contain a Diffie-Hellman public key,
an ElGamal ciphertext, etc.
*/
type PointInterface interface {
	Marshaling

	// Equality test for two Points derived from the same Group
	Equal(s2 *Point) bool

	Null() *Point // Set to neutral identity element

	// Set to this group's standard base point.
	Base() *Point

	// Pick and set to a point that is at least partly [pseudo-]random,
	// and optionally so as to encode a limited amount of specified data.
	// If data is nil, the point is completely [pseudo]-random.
	// Returns this Point and a slice containing the remaining data
	// following the data that was successfully embedded in this point.
	Pick(data []byte, rand cipher.Stream) (*Point, []byte)

	// Maximum number of bytes that can be reliably embedded
	// in a single group element via Pick().
	PickLen() int

	// Extract data embedded in a point chosen via Embed().
	// Returns an error if doesn't represent valid embedded data.
	Data() ([]byte, error)

	// Add points so that their secrets add homomorphically
	Add(a, b *Point) *Point

	// Subtract points so that their secrets subtract homomorphically
	Sub(a, b *Point) *Point

	// Set to the negation of point a
	Neg(a *Point) *Point

	// Encrypt point p by multiplying with secret s.
	// If p == nil, encrypt the standard base point Base().
	Mul(p *Point, s *Secret) *Point

	// GetSuite returns the suite used to create that point
	GetSuite() Suite

	// SetSuite sets the suite used for that point - only used in
	// setting up that point
	SetSuite(s Suite)
}

/*
Point is a structure that holds also the suite used for easy unmarshalling.
*/

type Point struct {
	PointInterface
}

/*
Payload is used to store the suite and help thus with the unmarshalling
*/
type Payload struct {
	// Suite is stored as string
	Suite string
}

// GetSuite returns the suite used with this point/secret
func (pl *Payload) GetSuite() Suite {
	if pl.Suite != "" {
		return StringToSuite(pl.Suite)
	}
	return nil
}

// SetSuite sets the suite used with this point/secret
func (pl *Payload) SetSuite(s Suite) {
	pl.Suite = s.String()
}

// MakePoint creates a Point from a PointInterface
func (pl *Payload) MakePoint(p PointInterface) *Point {
	return &Point{p}
}

// MakeSecret creates a Secret from a SecretInterface
func (pl *Payload) MakeSecret(s SecretInterface) *Secret {
	return &Secret{s}
}

/*
This interface represents an abstract cryptographic group
usable for Diffie-Hellman key exchange, ElGamal encryption,
and the related body of public-key cryptographic algorithms
and zero-knowledge proof methods.
The Group interface is designed in particular to be a generic front-end
to both traditional DSA-style modular arithmetic groups
and ECDSA-style elliptic curves:
the caller of this interface's methods
need not know or care which specific mathematical construction
underlies the interface.

The Group interface is essentially just a "constructor" interface
enabling the caller to generate the two particular types of objects
relevant to DSA-style public-key cryptography;
we call these objects Points and Secrets.
The caller must explicitly initialize or set a new Point or Secret object
to some value before using it as an input to some other operation
involving Point and/or Secret objects.
For example, to compare a point P against the neutral (identity) element,
you might use P.Equal(suite.Point().Null()),
but not just P.Equal(suite.Point()).

It is expected that any implementation of this interface
should satisfy suitable hardness assumptions for the applicable group:
e.g., that it is cryptographically hard for an adversary to
take an encrypted Point and the known generator it was based on,
and derive the Secret with which the Point was encrypted.
Any implementation is also expected to satisfy
the standard homomorphism properties that Diffie-Hellman
and the associated body of public-key cryptography are based on.

XXX should probably delete the somewhat redundant ...Len() methods.
*/
type Group interface {
	String() string

	SecretLen() int  // Max len of secrets in bytes
	Secret() *Secret // Create new secret

	PointLen() int // Max len of point in bytes
	Point() *Point // Create new point

	PrimeOrder() bool // Returns true if group is prime-order
}
