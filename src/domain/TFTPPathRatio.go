package domain

import "github.com/google/uuid"

//TFTPPathRatio TFTP path ratio entity
type TFTPPathRatio struct {
	Entity
	//TFTPConfigID TFTP config ID
	TFTPConfigID uuid.UUID
	//ActualPath actual file path
	ActualPath string
	//VirtualPath virtual file path
	VirtualPath string
}
