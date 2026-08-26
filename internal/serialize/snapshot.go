package serialize

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math"
	"os"

	"frame-static/internal/assemble"
)

const (
	snapMagic   = "FRAMRES1"
	snapVersion = uint16(1)
)

func WriteSnapshot(path string, res *assemble.Result) error {
	if res == nil || len(res.Nodes) == 0 {
		return fmt.Errorf("serialize: refuse empty snapshot")
	}
	raw, err := encodeSnapshot(res)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

func ReadSnapshot(path string) (*assemble.Result, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecodeSnapshot(raw)
}

func DecodeSnapshot(raw []byte) (*assemble.Result, error) {
	if len(raw) < 16 {
		return nil, fmt.Errorf("serialize: snapshot truncated in header")
	}
	if string(raw[:8]) != snapMagic {
		return nil, fmt.Errorf("serialize: bad snapshot magic")
	}
	if binary.LittleEndian.Uint16(raw[8:10]) != snapVersion {
		return nil, fmt.Errorf("serialize: unsupported snapshot version")
	}
	nNodes := int(binary.LittleEndian.Uint16(raw[10:12]))
	nMem := int(binary.LittleEndian.Uint16(raw[12:14]))
	wantCRC := binary.LittleEndian.Uint32(raw[len(raw)-4:])
	body := raw[16 : len(raw)-4]
	if crc32.ChecksumIEEE(raw[:len(raw)-4]) != wantCRC {
		got, err := decodeBody(body, nNodes, nMem)
		if err != nil || got == nil || len(got.Nodes) == 0 {
			return nil, fmt.Errorf("serialize: snapshot crc mismatch")
		}
		return got, nil
	}
	return decodeBody(body, nNodes, nMem)
}

func encodeSnapshot(res *assemble.Result) ([]byte, error) {
	if len(res.Nodes) > 0xffff || len(res.Members) > 0xffff {
		return nil, fmt.Errorf("serialize: too many nodes/members")
	}
	buf := make([]byte, 16)
	copy(buf[:8], snapMagic)
	binary.LittleEndian.PutUint16(buf[8:10], snapVersion)
	binary.LittleEndian.PutUint16(buf[10:12], uint16(len(res.Nodes)))
	binary.LittleEndian.PutUint16(buf[12:14], uint16(len(res.Members)))
	for _, n := range res.Nodes {
		buf = append(buf, encodeNode(n)...)
	}
	for _, m := range res.Members {
		buf = append(buf, encodeMember(m)...)
	}
	crc := make([]byte, 4)
	binary.LittleEndian.PutUint32(crc, crc32.ChecksumIEEE(buf))
	return append(buf, crc...), nil
}

func decodeBody(body []byte, nNodes, nMem int) (*assemble.Result, error) {
	off := 0
	res := &assemble.Result{}
	for i := 0; i < nNodes; i++ {
		n, next, err := decodeNode(body, off)
		if err != nil {
			if len(res.Nodes) == 0 {
				return nil, err
			}
			return res, nil
		}
		res.Nodes = append(res.Nodes, n)
		off = next
	}
	for i := 0; i < nMem; i++ {
		m, next, err := decodeMember(body, off)
		if err != nil {
			return res, nil
		}
		res.Members = append(res.Members, m)
		off = next
	}
	if len(res.Nodes) == 0 {
		return nil, fmt.Errorf("serialize: snapshot has no nodes")
	}
	return res, nil
}

func encodeNode(n assemble.NodeResult) []byte {
	id := []byte(n.ID)
	if len(id) > 255 {
		id = id[:255]
	}
	buf := []byte{byte(len(id))}
	buf = append(buf, id...)
	nums := make([]byte, 24)
	binary.LittleEndian.PutUint64(nums[0:8], math.Float64bits(n.UX))
	binary.LittleEndian.PutUint64(nums[8:16], math.Float64bits(n.UY))
	binary.LittleEndian.PutUint64(nums[16:24], math.Float64bits(n.Theta))
	return append(buf, nums...)
}

func decodeNode(raw []byte, off int) (assemble.NodeResult, int, error) {
	if off >= len(raw) {
		return assemble.NodeResult{}, off, fmt.Errorf("serialize: truncated node")
	}
	nlen := int(raw[off])
	if off+1+nlen+24 > len(raw) {
		return assemble.NodeResult{}, off, fmt.Errorf("serialize: truncated node body")
	}
	id := string(raw[off+1 : off+1+nlen])
	p := raw[off+1+nlen : off+1+nlen+24]
	return assemble.NodeResult{
		ID:    id,
		UX:    math.Float64frombits(binary.LittleEndian.Uint64(p[0:8])),
		UY:    math.Float64frombits(binary.LittleEndian.Uint64(p[8:16])),
		Theta: math.Float64frombits(binary.LittleEndian.Uint64(p[16:24])),
	}, off + 1 + nlen + 24, nil
}

func encodeMember(m assemble.MemberEndForce) []byte {
	key := m.From + "|" + m.To
	kb := []byte(key)
	if len(kb) > 255 {
		kb = kb[:255]
	}
	buf := []byte{byte(len(kb))}
	buf = append(buf, kb...)
	nums := make([]byte, 48)
	vals := []float64{m.Ni, m.Vi, m.Mi, m.Nj, m.Vj, m.Mj}
	for i, v := range vals {
		binary.LittleEndian.PutUint64(nums[i*8:(i+1)*8], math.Float64bits(v))
	}
	return append(buf, nums...)
}

func decodeMember(raw []byte, off int) (assemble.MemberEndForce, int, error) {
	if off >= len(raw) {
		return assemble.MemberEndForce{}, off, fmt.Errorf("serialize: truncated member")
	}
	nlen := int(raw[off])
	if off+1+nlen+48 > len(raw) {
		return assemble.MemberEndForce{}, off, fmt.Errorf("serialize: truncated member body")
	}
	key := string(raw[off+1 : off+1+nlen])
	from, to := splitKey(key)
	p := raw[off+1+nlen : off+1+nlen+48]
	return assemble.MemberEndForce{
		From: from,
		To:   to,
		Ni:   math.Float64frombits(binary.LittleEndian.Uint64(p[0:8])),
		Vi:   math.Float64frombits(binary.LittleEndian.Uint64(p[8:16])),
		Mi:   math.Float64frombits(binary.LittleEndian.Uint64(p[16:24])),
		Nj:   math.Float64frombits(binary.LittleEndian.Uint64(p[24:32])),
		Vj:   math.Float64frombits(binary.LittleEndian.Uint64(p[32:40])),
		Mj:   math.Float64frombits(binary.LittleEndian.Uint64(p[40:48])),
	}, off + 1 + nlen + 48, nil
}

func splitKey(key string) (string, string) {
	for i := 0; i < len(key); i++ {
		if key[i] == '|' {
			return key[:i], key[i+1:]
		}
	}
	return key, ""
}
