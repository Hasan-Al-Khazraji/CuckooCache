package proto

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// Shape: Op (1B), KeyLen (2B), ValLen (4B), Key (KeyLen), Value (ValLen)

func ReadRequest(r *bufio.Reader, maxKey uint16, maxVal uint32) (*Request, error) {
	op, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if op != OpGet && op != OpSet {
		return nil, fmt.Errorf("invalid op: 0x%x", op)
	}

	var keylen uint16
	if err := binary.Read(r, binary.BigEndian, &keylen); err != nil {
		return nil, err
	}
	if keylen > maxKey {
		return nil, fmt.Errorf("key too large: %d > %d", keylen, maxKey)
	}

	var valLen uint32
	if err := binary.Read(r, binary.BigEndian, &valLen); err != nil {
		return nil, err
	}
	if valLen > maxVal {
		return nil, fmt.Errorf("value too large: %d > %d", valLen, maxVal)
	}
	if op == OpGet && valLen != 0 {
		return nil, fmt.Errorf("GET must have valLen=0, got %d", valLen)
	}

	key := make([]byte, keylen)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}

	value := make([]byte, valLen)
	if _, err := io.ReadFull(r, value); err != nil {
		return nil, err
	}

	return &Request{
		Op:    op,
		Key:   key,
		Value: value,
	}, nil
}

func WriteRequest(w *bufio.Writer, req *Request) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.Op != OpGet && req.Op != OpSet {
		return fmt.Errorf("invalid op: 0x%x", req.Op)
	}

	keyLen := len(req.Key)
	if keyLen > (1<<16 - 1) {
		return fmt.Errorf("key too large: %d > %d", keyLen, (1<<16 - 1))
	}

	valLen := len(req.Value)
	if req.Op == OpGet && valLen != 0 {
		return fmt.Errorf("GET must have empty value,, got %d bytes", valLen)
	}
	if valLen > int(^uint32(0)) {
		return fmt.Errorf("value too large: %d > %d", valLen, int(^uint32(0)))
	}

	if err := w.WriteByte(req.Op); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(keyLen)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint32(valLen)); err != nil {
		return err
	}
	if _, err := w.Write(req.Key); err != nil {
		return err
	}
	if _, err := w.Write(req.Value); err != nil {
		return err
	}

	return nil
}

// Shape: Status (1B), ValLen (4B), Value (ValLen)
func ReadResponse(r *bufio.Reader, maxVal uint32) (*Response, error) {
	status, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if status != StatusOK && status != StatusNotFound && status != StatusErr {
		return nil, fmt.Errorf("status is not a valid type: 0x%x", status)
	}

	var valLen uint32
	if err := binary.Read(r, binary.BigEndian, &valLen); err != nil {
		return nil, err
	}
	if valLen > maxVal {
		return nil, fmt.Errorf("length of variable is too large: %d > %d", valLen, maxVal)
	}

	value := make([]byte, valLen)
	if _, err := io.ReadFull(r, value); err != nil {
		return nil, err
	}
	return &Response{Status: status, Value: value}, nil
}

func WriteResponse(w *bufio.Writer, res *Response) error {
	if res == nil {
		return fmt.Errorf("nil response")
	}
	if res.Status != StatusOK && res.Status != StatusNotFound && res.Status != StatusErr {
		return fmt.Errorf("invalid status: 0x%x", res.Status)
	}

	valLen := len(res.Value)
	if valLen > int(^uint32(0)) {
		return fmt.Errorf("value too big: %d > %d", valLen, int(^uint32(0)))
	}

	if err := w.WriteByte(res.Status); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint32(valLen)); err != nil {
		return err
	}

	if _, err := w.Write(res.Value); err != nil {
		return err
	}
	return nil
}
