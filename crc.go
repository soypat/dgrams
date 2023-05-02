package dgrams

import "encoding/binary"

// CRC_RFC791 function as defined by RFC 791. The Checksum field for TCP+IP
// is the 16-bit ones' complement of the ones' complement sum of
// all 16-bit words in the header. In case of uneven number of octet the
// last word is LSB padded with zeros.
type CRC_RFC791 struct {
	sum      uint32
	excedent uint8
	needsPad bool
}

func (c *CRC_RFC791) Write(buff []byte) (n int, err error) {
	// automatic padding of uneven data
	if c.needsPad {
		c.sum += uint32(c.excedent)<<8 + uint32(buff[0])
		buff = buff[1:]
		c.needsPad = false
	}
	n = len(buff) / 2
	if len(buff)%2 != 0 {
		c.excedent = buff[len(buff)-1]
		buff = buff[:len(buff)-1]
		c.needsPad = true
	}
	for i := 0; i < n; i++ {
		c.sum += uint32(binary.BigEndian.Uint16(buff[i*2 : i*2+2]))
	}
	return len(buff), nil
}

func (c *CRC_RFC791) Sum() uint16 {
	if c.needsPad {
		c.sum += uint32(c.excedent) << 8
		c.needsPad = false
	}
	for c.sum > 0xffff {
		c.sum = c.sum&0xffff + c.sum>>16
	}
	return ^uint16(c.sum)
}

func (c *CRC_RFC791) Reset() {
	c.sum = 0
	c.excedent = 0
	c.needsPad = false
}
