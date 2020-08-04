// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goff (v0.3.1) DO NOT EDIT

// Package fp contains field arithmetic operations
package fp

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math/big"
	"math/bits"
	"strconv"
	"sync"
	"unsafe"
)

// Element represents a field element stored on 12 words (uint64)
// Element are assumed to be in Montgomery form in all methods
// field modulus q =
//
// 6891450384315732539396789682275657542479668912536150109513790160209623422243491736087683183289411687640864567753786613451161759120554247759349511699125301598951605099378508850372543631423596795951899700429969112842764913119068299
type Element [12]uint64

// Limbs number of 64 bits words needed to represent Element
const Limbs = 12

// Bits number bits needed to represent Element
const Bits = 761

// field modulus stored as big.Int
var _modulus big.Int
var onceModulus sync.Once

// Modulus returns q as a big.Int
// q =
//
// 6891450384315732539396789682275657542479668912536150109513790160209623422243491736087683183289411687640864567753786613451161759120554247759349511699125301598951605099378508850372543631423596795951899700429969112842764913119068299
func Modulus() *big.Int {
	onceModulus.Do(func() {
		_modulus.SetString("6891450384315732539396789682275657542479668912536150109513790160209623422243491736087683183289411687640864567753786613451161759120554247759349511699125301598951605099378508850372543631423596795951899700429969112842764913119068299", 10)
	})
	return &_modulus
}

// q (modulus)
var qElement = Element{
	17626244516597989515,
	16614129118623039618,
	1588918198704579639,
	10998096788944562424,
	8204665564953313070,
	9694500593442880912,
	274362232328168196,
	8105254717682411801,
	5945444129596489281,
	13341377791855249032,
	15098257552581525310,
	81882988782276106,
}

// q'[0], see montgommery multiplication algorithm
var qElementInv0 uint64 = 744663313386281181

// rSquare
var rSquare = Element{
	14305184132582319705,
	8868935336694416555,
	9196887162930508889,
	15486798265448570248,
	5402985275949444416,
	10893197322525159598,
	3204916688966998390,
	12417238192559061753,
	12426306557607898622,
	1305582522441154384,
	10311846026977660324,
	48736111365249031,
}

// Bytes returns the regular (non montgomery) value
// of z as a big-endian byte slice.
func (z *Element) Bytes() []byte {
	var _z Element
	_z.Set(z).FromMont()
	res := make([]byte, Limbs*8)
	binary.BigEndian.PutUint64(res[(Limbs-1)*8:], _z[0])
	for i := Limbs - 2; i >= 0; i-- {
		binary.BigEndian.PutUint64(res[i*8:(i+1)*8], _z[Limbs-1-i])
	}
	return res
}

// SetBytes interprets e as the bytes of a big-endian unsigned integer,
// sets z to that value (in Montgomery form), and returns z.
func (z *Element) SetBytes(e []byte) *Element {
	var tmp big.Int
	tmp.SetBytes(e)
	z.SetBigInt(&tmp)
	return z
}

// SetUint64 z = v, sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
func (z *Element) SetUint64(v uint64) *Element {
	z[0] = v
	z[1] = 0
	z[2] = 0
	z[3] = 0
	z[4] = 0
	z[5] = 0
	z[6] = 0
	z[7] = 0
	z[8] = 0
	z[9] = 0
	z[10] = 0
	z[11] = 0
	return z.ToMont()
}

// Set z = x
func (z *Element) Set(x *Element) *Element {
	z[0] = x[0]
	z[1] = x[1]
	z[2] = x[2]
	z[3] = x[3]
	z[4] = x[4]
	z[5] = x[5]
	z[6] = x[6]
	z[7] = x[7]
	z[8] = x[8]
	z[9] = x[9]
	z[10] = x[10]
	z[11] = x[11]
	return z
}

// SetInterface converts i1 from uint64, int, string, or Element, big.Int into Element
// panic if provided type is not supported
func (z *Element) SetInterface(i1 interface{}) *Element {
	switch c1 := i1.(type) {
	case Element:
		return z.Set(&c1)
	case *Element:
		return z.Set(c1)
	case uint64:
		return z.SetUint64(c1)
	case int:
		return z.SetString(strconv.Itoa(c1))
	case string:
		return z.SetString(c1)
	case *big.Int:
		return z.SetBigInt(c1)
	case big.Int:
		return z.SetBigInt(&c1)
	case []byte:
		return z.SetBytes(c1)
	default:
		panic("invalid type")
	}
}

// SetZero z = 0
func (z *Element) SetZero() *Element {
	z[0] = 0
	z[1] = 0
	z[2] = 0
	z[3] = 0
	z[4] = 0
	z[5] = 0
	z[6] = 0
	z[7] = 0
	z[8] = 0
	z[9] = 0
	z[10] = 0
	z[11] = 0
	return z
}

// SetOne z = 1 (in Montgomery form)
func (z *Element) SetOne() *Element {
	z[0] = 144959613005956565
	z[1] = 6509995272855063783
	z[2] = 11428286765660613342
	z[3] = 15738672438262922740
	z[4] = 17071399330169272331
	z[5] = 13899911246788437003
	z[6] = 12055474021000362245
	z[7] = 2545351818702954755
	z[8] = 8887388221587179644
	z[9] = 5009280847225881135
	z[10] = 15539704305423854047
	z[11] = 23071597697427581
	return z
}

// Div z = x*y^-1 mod q
func (z *Element) Div(x, y *Element) *Element {
	var yInv Element
	yInv.Inverse(y)
	z.Mul(x, &yInv)
	return z
}

// Equal returns z == x
func (z *Element) Equal(x *Element) bool {
	return (z[11] == x[11]) && (z[10] == x[10]) && (z[9] == x[9]) && (z[8] == x[8]) && (z[7] == x[7]) && (z[6] == x[6]) && (z[5] == x[5]) && (z[4] == x[4]) && (z[3] == x[3]) && (z[2] == x[2]) && (z[1] == x[1]) && (z[0] == x[0])
}

// IsZero returns z == 0
func (z *Element) IsZero() bool {
	return z[11] == 0 && z[10] == 0 && z[9] == 0 && z[8] == 0 && z[7] == 0 && z[6] == 0 && z[5] == 0 && z[4] == 0 && z[3] == 0 && z[2] == 0 && z[1] == 0 && z[0] == 0
}

// SetRandom sets z to a random element < q
func (z *Element) SetRandom() *Element {
	bytes := make([]byte, 96)
	io.ReadFull(rand.Reader, bytes)
	z[0] = binary.BigEndian.Uint64(bytes[0:8])
	z[1] = binary.BigEndian.Uint64(bytes[8:16])
	z[2] = binary.BigEndian.Uint64(bytes[16:24])
	z[3] = binary.BigEndian.Uint64(bytes[24:32])
	z[4] = binary.BigEndian.Uint64(bytes[32:40])
	z[5] = binary.BigEndian.Uint64(bytes[40:48])
	z[6] = binary.BigEndian.Uint64(bytes[48:56])
	z[7] = binary.BigEndian.Uint64(bytes[56:64])
	z[8] = binary.BigEndian.Uint64(bytes[64:72])
	z[9] = binary.BigEndian.Uint64(bytes[72:80])
	z[10] = binary.BigEndian.Uint64(bytes[80:88])
	z[11] = binary.BigEndian.Uint64(bytes[88:96])
	z[11] %= 81882988782276106

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[11] < 81882988782276106 || (z[11] == 81882988782276106 && (z[10] < 15098257552581525310 || (z[10] == 15098257552581525310 && (z[9] < 13341377791855249032 || (z[9] == 13341377791855249032 && (z[8] < 5945444129596489281 || (z[8] == 5945444129596489281 && (z[7] < 8105254717682411801 || (z[7] == 8105254717682411801 && (z[6] < 274362232328168196 || (z[6] == 274362232328168196 && (z[5] < 9694500593442880912 || (z[5] == 9694500593442880912 && (z[4] < 8204665564953313070 || (z[4] == 8204665564953313070 && (z[3] < 10998096788944562424 || (z[3] == 10998096788944562424 && (z[2] < 1588918198704579639 || (z[2] == 1588918198704579639 && (z[1] < 16614129118623039618 || (z[1] == 16614129118623039618 && (z[0] < 17626244516597989515))))))))))))))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 17626244516597989515, 0)
		z[1], b = bits.Sub64(z[1], 16614129118623039618, b)
		z[2], b = bits.Sub64(z[2], 1588918198704579639, b)
		z[3], b = bits.Sub64(z[3], 10998096788944562424, b)
		z[4], b = bits.Sub64(z[4], 8204665564953313070, b)
		z[5], b = bits.Sub64(z[5], 9694500593442880912, b)
		z[6], b = bits.Sub64(z[6], 274362232328168196, b)
		z[7], b = bits.Sub64(z[7], 8105254717682411801, b)
		z[8], b = bits.Sub64(z[8], 5945444129596489281, b)
		z[9], b = bits.Sub64(z[9], 13341377791855249032, b)
		z[10], b = bits.Sub64(z[10], 15098257552581525310, b)
		z[11], _ = bits.Sub64(z[11], 81882988782276106, b)
	}

	return z
}

// One returns 1 (in montgommery form)
func One() Element {
	var one Element
	one.SetOne()
	return one
}

// MulAssign is deprecated
// Deprecated: use Mul instead
func (z *Element) MulAssign(x *Element) *Element {
	return z.Mul(z, x)
}

// AddAssign is deprecated
// Deprecated: use Add instead
func (z *Element) AddAssign(x *Element) *Element {
	return z.Add(z, x)
}

// SubAssign is deprecated
// Deprecated: use Sub instead
func (z *Element) SubAssign(x *Element) *Element {
	return z.Sub(z, x)
}

// API with assembly impl

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Mul(x, y *Element) *Element {
	mul(z, x, y)
	return z
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Square(x *Element) *Element {
	square(z, x)
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *Element) FromMont() *Element {
	fromMont(z)
	return z
}

// Add z = x + y mod q
func (z *Element) Add(x, y *Element) *Element {
	add(z, x, y)
	return z
}

// Double z = x + x mod q, aka Lsh 1
func (z *Element) Double(x *Element) *Element {
	double(z, x)
	return z
}

// Sub  z = x - y mod q
func (z *Element) Sub(x, y *Element) *Element {
	sub(z, x, y)
	return z
}

// Neg z = q - x
func (z *Element) Neg(x *Element) *Element {
	neg(z, x)
	return z
}

// Exp z = x^exponent mod q
func (z *Element) Exp(x Element, exponent *big.Int) *Element {
	var bZero big.Int
	if exponent.Cmp(&bZero) == 0 {
		return z.SetOne()
	}

	z.Set(&x)

	for i := exponent.BitLen() - 2; i >= 0; i-- {
		z.Square(z)
		if exponent.Bit(i) == 1 {
			z.Mul(z, &x)
		}
	}

	return z
}

// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *Element) ToMont() *Element {
	return z.Mul(z, &rSquare)
}

// ToRegular returns z in regular form (doesn't mutate z)
func (z Element) ToRegular() Element {
	return *z.FromMont()
}

// String returns the string form of an Element in Montgomery form
func (z *Element) String() string {
	var _z big.Int
	return z.ToBigIntRegular(&_z).String()
}

// ToBigInt returns z as a big.Int in Montgomery form
func (z *Element) ToBigInt(res *big.Int) *big.Int {
	bits := (*[12]big.Word)(unsafe.Pointer(z))
	return res.SetBits(bits[:])
}

// ToBigIntRegular returns z as a big.Int in regular form
func (z Element) ToBigIntRegular(res *big.Int) *big.Int {
	z.FromMont()
	bits := (*[12]big.Word)(unsafe.Pointer(&z))
	return res.SetBits(bits[:])
}

// SetBigInt sets z to v (regular form) and returns z in Montgomery form
func (z *Element) SetBigInt(v *big.Int) *Element {
	z.SetZero()

	zero := big.NewInt(0)
	q := Modulus()

	// fast path
	c := v.Cmp(q)
	if c == 0 {
		return z
	} else if c != 1 && v.Cmp(zero) != -1 {
		// v should
		vBits := v.Bits()
		for i := 0; i < len(vBits); i++ {
			z[i] = uint64(vBits[i])
		}
		return z.ToMont()
	}

	// copy input
	vv := new(big.Int).Set(v)
	vv.Mod(v, q)

	// v should
	vBits := vv.Bits()
	for i := 0; i < len(vBits); i++ {
		z[i] = uint64(vBits[i])
	}
	return z.ToMont()
}

// SetString creates a big.Int with s (in base 10) and calls SetBigInt on z
func (z *Element) SetString(s string) *Element {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("Element.SetString failed -> can't parse number in base10 into a big.Int")
	}
	return z.SetBigInt(x)
}

var (
	_bLegendreExponentElement *big.Int
	_bSqrtExponentElement     *big.Int
)

func init() {
	_bLegendreExponentElement, _ = new(big.Int).SetString("9174127dc1e70568c3e4a0027d7f9f5c930c3540e8a34429413af7c043df20b83dd31c72c2748c81e75d7f92da11824344e476897cfec838ee69ee39f5ff974c508b612b33d47c0b067c577578521bf3489f34380000417a4e800000000045", 16)
	const sqrtExponentElement = "48ba093ee0f382b461f250013ebfcfae49861aa07451a214a09d7be021ef905c1ee98e39613a4640f3aebfc96d08c121a2723b44be7f641c7734f71cfaffcba62845b09599ea3e05833e2bbabc290df9a44f9a1c000020bd27400000000023"
	_bSqrtExponentElement, _ = new(big.Int).SetString(sqrtExponentElement, 16)
}

// Legendre returns the Legendre symbol of z (either +1, -1, or 0.)
func (z *Element) Legendre() int {
	var l Element
	// z^((q-1)/2)
	l.Exp(*z, _bLegendreExponentElement)

	if l.IsZero() {
		return 0
	}

	// if l == 1
	if (l[11] == 23071597697427581) && (l[10] == 15539704305423854047) && (l[9] == 5009280847225881135) && (l[8] == 8887388221587179644) && (l[7] == 2545351818702954755) && (l[6] == 12055474021000362245) && (l[5] == 13899911246788437003) && (l[4] == 17071399330169272331) && (l[3] == 15738672438262922740) && (l[2] == 11428286765660613342) && (l[1] == 6509995272855063783) && (l[0] == 144959613005956565) {
		return 1
	}
	return -1
}

// Sqrt z = √x mod q
// if the square root doesn't exist (x is not a square mod q)
// Sqrt leaves z unchanged and returns nil
func (z *Element) Sqrt(x *Element) *Element {
	// q ≡ 3 (mod 4)
	// using  z ≡ ± x^((p+1)/4) (mod q)
	var y, square Element
	y.Exp(*x, _bSqrtExponentElement)
	// as we didn't compute the legendre symbol, ensure we found y such that y * y = x
	square.Square(&y)
	if square.Equal(x) {
		return z.Set(&y)
	}
	return nil
}

// Inverse z = x^-1 mod q
// Algorithm 16 in "Efficient Software-Implementation of Finite Fields with Applications to Cryptography"
// if x == 0, sets and returns z = x
func (z *Element) Inverse(x *Element) *Element {
	if x.IsZero() {
		return z.Set(x)
	}

	// initialize u = q
	var u = Element{
		17626244516597989515,
		16614129118623039618,
		1588918198704579639,
		10998096788944562424,
		8204665564953313070,
		9694500593442880912,
		274362232328168196,
		8105254717682411801,
		5945444129596489281,
		13341377791855249032,
		15098257552581525310,
		81882988782276106,
	}

	// initialize s = r^2
	var s = Element{
		14305184132582319705,
		8868935336694416555,
		9196887162930508889,
		15486798265448570248,
		5402985275949444416,
		10893197322525159598,
		3204916688966998390,
		12417238192559061753,
		12426306557607898622,
		1305582522441154384,
		10311846026977660324,
		48736111365249031,
	}

	// r = 0
	r := Element{}

	v := *x

	var carry, borrow, t, t2 uint64
	var bigger, uIsOne, vIsOne bool

	for !uIsOne && !vIsOne {
		for v[0]&1 == 0 {

			// v = v >> 1
			t2 = v[11] << 63
			v[11] >>= 1
			t = t2
			t2 = v[10] << 63
			v[10] = (v[10] >> 1) | t
			t = t2
			t2 = v[9] << 63
			v[9] = (v[9] >> 1) | t
			t = t2
			t2 = v[8] << 63
			v[8] = (v[8] >> 1) | t
			t = t2
			t2 = v[7] << 63
			v[7] = (v[7] >> 1) | t
			t = t2
			t2 = v[6] << 63
			v[6] = (v[6] >> 1) | t
			t = t2
			t2 = v[5] << 63
			v[5] = (v[5] >> 1) | t
			t = t2
			t2 = v[4] << 63
			v[4] = (v[4] >> 1) | t
			t = t2
			t2 = v[3] << 63
			v[3] = (v[3] >> 1) | t
			t = t2
			t2 = v[2] << 63
			v[2] = (v[2] >> 1) | t
			t = t2
			t2 = v[1] << 63
			v[1] = (v[1] >> 1) | t
			t = t2
			v[0] = (v[0] >> 1) | t

			if s[0]&1 == 1 {

				// s = s + q
				s[0], carry = bits.Add64(s[0], 17626244516597989515, 0)
				s[1], carry = bits.Add64(s[1], 16614129118623039618, carry)
				s[2], carry = bits.Add64(s[2], 1588918198704579639, carry)
				s[3], carry = bits.Add64(s[3], 10998096788944562424, carry)
				s[4], carry = bits.Add64(s[4], 8204665564953313070, carry)
				s[5], carry = bits.Add64(s[5], 9694500593442880912, carry)
				s[6], carry = bits.Add64(s[6], 274362232328168196, carry)
				s[7], carry = bits.Add64(s[7], 8105254717682411801, carry)
				s[8], carry = bits.Add64(s[8], 5945444129596489281, carry)
				s[9], carry = bits.Add64(s[9], 13341377791855249032, carry)
				s[10], carry = bits.Add64(s[10], 15098257552581525310, carry)
				s[11], _ = bits.Add64(s[11], 81882988782276106, carry)

			}

			// s = s >> 1
			t2 = s[11] << 63
			s[11] >>= 1
			t = t2
			t2 = s[10] << 63
			s[10] = (s[10] >> 1) | t
			t = t2
			t2 = s[9] << 63
			s[9] = (s[9] >> 1) | t
			t = t2
			t2 = s[8] << 63
			s[8] = (s[8] >> 1) | t
			t = t2
			t2 = s[7] << 63
			s[7] = (s[7] >> 1) | t
			t = t2
			t2 = s[6] << 63
			s[6] = (s[6] >> 1) | t
			t = t2
			t2 = s[5] << 63
			s[5] = (s[5] >> 1) | t
			t = t2
			t2 = s[4] << 63
			s[4] = (s[4] >> 1) | t
			t = t2
			t2 = s[3] << 63
			s[3] = (s[3] >> 1) | t
			t = t2
			t2 = s[2] << 63
			s[2] = (s[2] >> 1) | t
			t = t2
			t2 = s[1] << 63
			s[1] = (s[1] >> 1) | t
			t = t2
			s[0] = (s[0] >> 1) | t

		}
		for u[0]&1 == 0 {

			// u = u >> 1
			t2 = u[11] << 63
			u[11] >>= 1
			t = t2
			t2 = u[10] << 63
			u[10] = (u[10] >> 1) | t
			t = t2
			t2 = u[9] << 63
			u[9] = (u[9] >> 1) | t
			t = t2
			t2 = u[8] << 63
			u[8] = (u[8] >> 1) | t
			t = t2
			t2 = u[7] << 63
			u[7] = (u[7] >> 1) | t
			t = t2
			t2 = u[6] << 63
			u[6] = (u[6] >> 1) | t
			t = t2
			t2 = u[5] << 63
			u[5] = (u[5] >> 1) | t
			t = t2
			t2 = u[4] << 63
			u[4] = (u[4] >> 1) | t
			t = t2
			t2 = u[3] << 63
			u[3] = (u[3] >> 1) | t
			t = t2
			t2 = u[2] << 63
			u[2] = (u[2] >> 1) | t
			t = t2
			t2 = u[1] << 63
			u[1] = (u[1] >> 1) | t
			t = t2
			u[0] = (u[0] >> 1) | t

			if r[0]&1 == 1 {

				// r = r + q
				r[0], carry = bits.Add64(r[0], 17626244516597989515, 0)
				r[1], carry = bits.Add64(r[1], 16614129118623039618, carry)
				r[2], carry = bits.Add64(r[2], 1588918198704579639, carry)
				r[3], carry = bits.Add64(r[3], 10998096788944562424, carry)
				r[4], carry = bits.Add64(r[4], 8204665564953313070, carry)
				r[5], carry = bits.Add64(r[5], 9694500593442880912, carry)
				r[6], carry = bits.Add64(r[6], 274362232328168196, carry)
				r[7], carry = bits.Add64(r[7], 8105254717682411801, carry)
				r[8], carry = bits.Add64(r[8], 5945444129596489281, carry)
				r[9], carry = bits.Add64(r[9], 13341377791855249032, carry)
				r[10], carry = bits.Add64(r[10], 15098257552581525310, carry)
				r[11], _ = bits.Add64(r[11], 81882988782276106, carry)

			}

			// r = r >> 1
			t2 = r[11] << 63
			r[11] >>= 1
			t = t2
			t2 = r[10] << 63
			r[10] = (r[10] >> 1) | t
			t = t2
			t2 = r[9] << 63
			r[9] = (r[9] >> 1) | t
			t = t2
			t2 = r[8] << 63
			r[8] = (r[8] >> 1) | t
			t = t2
			t2 = r[7] << 63
			r[7] = (r[7] >> 1) | t
			t = t2
			t2 = r[6] << 63
			r[6] = (r[6] >> 1) | t
			t = t2
			t2 = r[5] << 63
			r[5] = (r[5] >> 1) | t
			t = t2
			t2 = r[4] << 63
			r[4] = (r[4] >> 1) | t
			t = t2
			t2 = r[3] << 63
			r[3] = (r[3] >> 1) | t
			t = t2
			t2 = r[2] << 63
			r[2] = (r[2] >> 1) | t
			t = t2
			t2 = r[1] << 63
			r[1] = (r[1] >> 1) | t
			t = t2
			r[0] = (r[0] >> 1) | t

		}

		// v >= u
		bigger = !(v[11] < u[11] || (v[11] == u[11] && (v[10] < u[10] || (v[10] == u[10] && (v[9] < u[9] || (v[9] == u[9] && (v[8] < u[8] || (v[8] == u[8] && (v[7] < u[7] || (v[7] == u[7] && (v[6] < u[6] || (v[6] == u[6] && (v[5] < u[5] || (v[5] == u[5] && (v[4] < u[4] || (v[4] == u[4] && (v[3] < u[3] || (v[3] == u[3] && (v[2] < u[2] || (v[2] == u[2] && (v[1] < u[1] || (v[1] == u[1] && (v[0] < u[0])))))))))))))))))))))))

		if bigger {

			// v = v - u
			v[0], borrow = bits.Sub64(v[0], u[0], 0)
			v[1], borrow = bits.Sub64(v[1], u[1], borrow)
			v[2], borrow = bits.Sub64(v[2], u[2], borrow)
			v[3], borrow = bits.Sub64(v[3], u[3], borrow)
			v[4], borrow = bits.Sub64(v[4], u[4], borrow)
			v[5], borrow = bits.Sub64(v[5], u[5], borrow)
			v[6], borrow = bits.Sub64(v[6], u[6], borrow)
			v[7], borrow = bits.Sub64(v[7], u[7], borrow)
			v[8], borrow = bits.Sub64(v[8], u[8], borrow)
			v[9], borrow = bits.Sub64(v[9], u[9], borrow)
			v[10], borrow = bits.Sub64(v[10], u[10], borrow)
			v[11], _ = bits.Sub64(v[11], u[11], borrow)

			// r >= s
			bigger = !(r[11] < s[11] || (r[11] == s[11] && (r[10] < s[10] || (r[10] == s[10] && (r[9] < s[9] || (r[9] == s[9] && (r[8] < s[8] || (r[8] == s[8] && (r[7] < s[7] || (r[7] == s[7] && (r[6] < s[6] || (r[6] == s[6] && (r[5] < s[5] || (r[5] == s[5] && (r[4] < s[4] || (r[4] == s[4] && (r[3] < s[3] || (r[3] == s[3] && (r[2] < s[2] || (r[2] == s[2] && (r[1] < s[1] || (r[1] == s[1] && (r[0] < s[0])))))))))))))))))))))))

			if bigger {

				// s = s + q
				s[0], carry = bits.Add64(s[0], 17626244516597989515, 0)
				s[1], carry = bits.Add64(s[1], 16614129118623039618, carry)
				s[2], carry = bits.Add64(s[2], 1588918198704579639, carry)
				s[3], carry = bits.Add64(s[3], 10998096788944562424, carry)
				s[4], carry = bits.Add64(s[4], 8204665564953313070, carry)
				s[5], carry = bits.Add64(s[5], 9694500593442880912, carry)
				s[6], carry = bits.Add64(s[6], 274362232328168196, carry)
				s[7], carry = bits.Add64(s[7], 8105254717682411801, carry)
				s[8], carry = bits.Add64(s[8], 5945444129596489281, carry)
				s[9], carry = bits.Add64(s[9], 13341377791855249032, carry)
				s[10], carry = bits.Add64(s[10], 15098257552581525310, carry)
				s[11], _ = bits.Add64(s[11], 81882988782276106, carry)

			}

			// s = s - r
			s[0], borrow = bits.Sub64(s[0], r[0], 0)
			s[1], borrow = bits.Sub64(s[1], r[1], borrow)
			s[2], borrow = bits.Sub64(s[2], r[2], borrow)
			s[3], borrow = bits.Sub64(s[3], r[3], borrow)
			s[4], borrow = bits.Sub64(s[4], r[4], borrow)
			s[5], borrow = bits.Sub64(s[5], r[5], borrow)
			s[6], borrow = bits.Sub64(s[6], r[6], borrow)
			s[7], borrow = bits.Sub64(s[7], r[7], borrow)
			s[8], borrow = bits.Sub64(s[8], r[8], borrow)
			s[9], borrow = bits.Sub64(s[9], r[9], borrow)
			s[10], borrow = bits.Sub64(s[10], r[10], borrow)
			s[11], _ = bits.Sub64(s[11], r[11], borrow)

		} else {

			// u = u - v
			u[0], borrow = bits.Sub64(u[0], v[0], 0)
			u[1], borrow = bits.Sub64(u[1], v[1], borrow)
			u[2], borrow = bits.Sub64(u[2], v[2], borrow)
			u[3], borrow = bits.Sub64(u[3], v[3], borrow)
			u[4], borrow = bits.Sub64(u[4], v[4], borrow)
			u[5], borrow = bits.Sub64(u[5], v[5], borrow)
			u[6], borrow = bits.Sub64(u[6], v[6], borrow)
			u[7], borrow = bits.Sub64(u[7], v[7], borrow)
			u[8], borrow = bits.Sub64(u[8], v[8], borrow)
			u[9], borrow = bits.Sub64(u[9], v[9], borrow)
			u[10], borrow = bits.Sub64(u[10], v[10], borrow)
			u[11], _ = bits.Sub64(u[11], v[11], borrow)

			// s >= r
			bigger = !(s[11] < r[11] || (s[11] == r[11] && (s[10] < r[10] || (s[10] == r[10] && (s[9] < r[9] || (s[9] == r[9] && (s[8] < r[8] || (s[8] == r[8] && (s[7] < r[7] || (s[7] == r[7] && (s[6] < r[6] || (s[6] == r[6] && (s[5] < r[5] || (s[5] == r[5] && (s[4] < r[4] || (s[4] == r[4] && (s[3] < r[3] || (s[3] == r[3] && (s[2] < r[2] || (s[2] == r[2] && (s[1] < r[1] || (s[1] == r[1] && (s[0] < r[0])))))))))))))))))))))))

			if bigger {

				// r = r + q
				r[0], carry = bits.Add64(r[0], 17626244516597989515, 0)
				r[1], carry = bits.Add64(r[1], 16614129118623039618, carry)
				r[2], carry = bits.Add64(r[2], 1588918198704579639, carry)
				r[3], carry = bits.Add64(r[3], 10998096788944562424, carry)
				r[4], carry = bits.Add64(r[4], 8204665564953313070, carry)
				r[5], carry = bits.Add64(r[5], 9694500593442880912, carry)
				r[6], carry = bits.Add64(r[6], 274362232328168196, carry)
				r[7], carry = bits.Add64(r[7], 8105254717682411801, carry)
				r[8], carry = bits.Add64(r[8], 5945444129596489281, carry)
				r[9], carry = bits.Add64(r[9], 13341377791855249032, carry)
				r[10], carry = bits.Add64(r[10], 15098257552581525310, carry)
				r[11], _ = bits.Add64(r[11], 81882988782276106, carry)

			}

			// r = r - s
			r[0], borrow = bits.Sub64(r[0], s[0], 0)
			r[1], borrow = bits.Sub64(r[1], s[1], borrow)
			r[2], borrow = bits.Sub64(r[2], s[2], borrow)
			r[3], borrow = bits.Sub64(r[3], s[3], borrow)
			r[4], borrow = bits.Sub64(r[4], s[4], borrow)
			r[5], borrow = bits.Sub64(r[5], s[5], borrow)
			r[6], borrow = bits.Sub64(r[6], s[6], borrow)
			r[7], borrow = bits.Sub64(r[7], s[7], borrow)
			r[8], borrow = bits.Sub64(r[8], s[8], borrow)
			r[9], borrow = bits.Sub64(r[9], s[9], borrow)
			r[10], borrow = bits.Sub64(r[10], s[10], borrow)
			r[11], _ = bits.Sub64(r[11], s[11], borrow)

		}
		uIsOne = (u[0] == 1) && (u[11]|u[10]|u[9]|u[8]|u[7]|u[6]|u[5]|u[4]|u[3]|u[2]|u[1]) == 0
		vIsOne = (v[0] == 1) && (v[11]|v[10]|v[9]|v[8]|v[7]|v[6]|v[5]|v[4]|v[3]|v[2]|v[1]) == 0
	}

	if uIsOne {
		z.Set(&r)
	} else {
		z.Set(&s)
	}

	return z
}