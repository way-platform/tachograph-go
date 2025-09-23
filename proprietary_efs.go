package tachograph

import (
	"encoding/binary"
	"sync"
)

// Package-level storage for proprietary EFs during roundtrip processing
var (
	proprietaryEFsMap   = make(map[uintptr]*ProprietaryEFs)
	proprietaryEFsMutex sync.RWMutex
)

// ProprietaryEF represents a proprietary Elementary File that we don't fully parse
type ProprietaryEF struct {
	FID  uint16 // File Identifier
	Data []byte // Raw data
}

// ProprietaryEFs holds a collection of proprietary EFs found in a card file
type ProprietaryEFs struct {
	EFs []ProprietaryEF
}

// AddProprietaryEF adds a proprietary EF to the collection
func (p *ProprietaryEFs) AddProprietaryEF(fid uint16, data []byte) {
	p.EFs = append(p.EFs, ProprietaryEF{
		FID:  fid,
		Data: make([]byte, len(data)), // Make a copy
	})
	copy(p.EFs[len(p.EFs)-1].Data, data)
}

// AppendProprietaryEFs appends all proprietary EFs to the output data
func (p *ProprietaryEFs) AppendProprietaryEFs(data []byte) []byte {
	for _, ef := range p.EFs {
		// Write data block (FID + appendix 0x00 + length + data)
		data = binary.BigEndian.AppendUint16(data, ef.FID)
		data = append(data, 0x00) // appendix for data
		data = binary.BigEndian.AppendUint16(data, uint16(len(ef.Data)))
		data = append(data, ef.Data...)

		// Write signature block (FID + appendix 0x01 + 128 bytes signature)
		data = binary.BigEndian.AppendUint16(data, ef.FID)
		data = append(data, 0x01)                       // appendix for signature
		data = binary.BigEndian.AppendUint16(data, 128) // signature length
		signature := make([]byte, 128)                  // zeros for now
		data = append(data, signature...)
	}
	return data
}

// UnmarshalProprietaryEF handles unmarshalling of a proprietary EF
func UnmarshalProprietaryEF(fid uint16, value []byte, proprietaryEFs *ProprietaryEFs) {
	proprietaryEFs.AddProprietaryEF(fid, value)
}

// StoreProprietaryEFs stores proprietary EFs for a specific file pointer
func StoreProprietaryEFs(filePtr uintptr, proprietaryEFs *ProprietaryEFs) {
	proprietaryEFsMutex.Lock()
	defer proprietaryEFsMutex.Unlock()
	proprietaryEFsMap[filePtr] = proprietaryEFs
}

// GetProprietaryEFs retrieves proprietary EFs for a specific file pointer
func GetProprietaryEFs(filePtr uintptr) *ProprietaryEFs {
	proprietaryEFsMutex.RLock()
	defer proprietaryEFsMutex.RUnlock()
	return proprietaryEFsMap[filePtr]
}

// ClearProprietaryEFs removes proprietary EFs for a specific file pointer
func ClearProprietaryEFs(filePtr uintptr) {
	proprietaryEFsMutex.Lock()
	defer proprietaryEFsMutex.Unlock()
	delete(proprietaryEFsMap, filePtr)
}
