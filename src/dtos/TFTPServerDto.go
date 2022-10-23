package dtos

import "github.com/google/uuid"

//TFTPServerDto TFTP server dto
type TFTPServerDto struct {
	BaseDto[uuid.UUID]
	TFTPServerBaseDto
	//State of tftp server
	State string
}
