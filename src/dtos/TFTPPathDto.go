package dtos

import "github.com/google/uuid"

//TFTPPathDto TFTP path dto
type TFTPPathDto struct {
	BaseDto[uuid.UUID]
	TFTPPathBaseDto
}
