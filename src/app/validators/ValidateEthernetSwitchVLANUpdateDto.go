package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/app/errors"
	"rol/dtos"
)

//ValidateEthernetSwitchVLANUpdateDto validates ethernet switch VLAN update dto
//
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchVLANUpdateDto(dto dtos.EthernetSwitchVLANUpdateDto) error {
	var notInSliceErr error
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.TaggedPorts, []validation.Rule{
			validation.By(uuidSliceElemUniqueness),
		}...),
		validation.Field(&dto.UntaggedPorts, []validation.Rule{
			validation.By(uuidSliceElemUniqueness),
		}...))
	notInSliceErr = uuidsUniqueWithinSlices(dto.TaggedPorts, dto.UntaggedPorts)
	if notInSliceErr != nil {
		if err == nil {
			err = errors.Validation.New(errors.ValidationErrorMessage)
		}
		err = errors.AddErrorContext(err, "TaggedPorts", notInSliceErr.Error())
	}
	return convertOzzoErrorToValidationError(err)
}
