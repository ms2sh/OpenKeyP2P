package math

import "fmt"

// Uint24 speichert einen 24-Bit-Unsigned-Integer als Array von 3 Bytes.
type Uint24 [3]byte

// Uint32 gibt den Wert als uint32 zurück.
func (u Uint24) Uint32() uint32 {
	return uint32(u[0])<<16 | uint32(u[1])<<8 | uint32(u[2])
}

// NewUint24 erstellt einen Uint24 aus einem uint32-Wert.
// Es wird überprüft, ob der Wert in 24 Bit passt.
func NewUint24(x uint32) Uint24 {
	if x > 0xFFFFFF {
		panic("Wert zu groß für uint24")
	}
	return Uint24{
		byte(x >> 16), // höchstwertiges Byte
		byte(x >> 8),  // mittleres Byte
		byte(x),       // niederwertiges Byte
	}
}

// Bytes wandelt den Uint24 in ein Slice von 3 Bytes um.
func (u Uint24) Bytes() []byte {
	// Hier wird der zugrundeliegende Array als Slice zurückgegeben.
	return u[:]
}

// Uint24FromBytes erstellt einen Uint24 aus einem Byte-Slice.
// Der Slice muss genau 3 Byte lang sein.
func Uint24FromBytes(b []byte) (Uint24, error) {
	if len(b) != 3 {
		return Uint24{}, fmt.Errorf("ungültige Länge für Uint24: erwartet 3, erhalten %d", len(b))
	}
	var u Uint24
	copy(u[:], b)
	return u, nil
}
