package wgpu

// BufferMapAsyncStatus is the Go-side map callback status. Values are mapped from
// WGPUMapAsyncStatus in v29.go; they are not a 1:1 copy of the C enum.
type BufferMapAsyncStatus uint32

const (
	BufferMapAsyncStatusSuccess                 BufferMapAsyncStatus = 0x00000000
	BufferMapAsyncStatusValidationError         BufferMapAsyncStatus = 0x00000001
	BufferMapAsyncStatusUnknown                 BufferMapAsyncStatus = 0x00000002
	BufferMapAsyncStatusDeviceLost              BufferMapAsyncStatus = 0x00000003
	BufferMapAsyncStatusDestroyedBeforeCallback BufferMapAsyncStatus = 0x00000004
	BufferMapAsyncStatusUnmappedBeforeCallback  BufferMapAsyncStatus = 0x00000005
	BufferMapAsyncStatusMappingAlreadyPending   BufferMapAsyncStatus = 0x00000006
	BufferMapAsyncStatusOffsetOutOfRange        BufferMapAsyncStatus = 0x00000007
	BufferMapAsyncStatusSizeOutOfRange          BufferMapAsyncStatus = 0x00000008
)

func (v BufferMapAsyncStatus) String() string {
	switch v {
	case BufferMapAsyncStatusSuccess:
		return "success"
	case BufferMapAsyncStatusValidationError:
		return "validation-error"
	case BufferMapAsyncStatusUnknown:
		return "unknown"
	case BufferMapAsyncStatusDeviceLost:
		return "device-lost"
	case BufferMapAsyncStatusDestroyedBeforeCallback:
		return "destroyed-before-callback"
	case BufferMapAsyncStatusUnmappedBeforeCallback:
		return "unmapped-before-callback"
	case BufferMapAsyncStatusMappingAlreadyPending:
		return "mapping-already-pending"
	case BufferMapAsyncStatusOffsetOutOfRange:
		return "offset-out-of-range"
	case BufferMapAsyncStatusSizeOutOfRange:
		return "size-out-of-range"
	default:
		return ""
	}
}
