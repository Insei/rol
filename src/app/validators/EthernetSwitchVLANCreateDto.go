package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/app/errors"
	"rol/dtos"
)

//ValidateEthernetSwitchVLANCreateDto validates ethernet switch VLAN create dto
//
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchVLANCreateDto(dto dtos.EthernetSwitchVLANCreateDto) error {
	var (
		err           error
		notInSliceErr error
	)
	validationErr := validation.ValidateStruct(&dto,
		validation.Field(&dto.VlanID, []validation.Rule{
			validation.Required,
			validation.Min(1),
		}...),
		validation.Field(&dto.TaggedPorts, []validation.Rule{
			validation.By(uuidSliceElemUniqueness),
		}...),
		validation.Field(&dto.UntaggedPorts, []validation.Rule{
			validation.By(uuidSliceElemUniqueness),
		}...))
	notInSliceErr = uuidsUniqueWithinSlices(dto.TaggedPorts, dto.UntaggedPorts)
	if validationErr != nil {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		for key, value := range validationErr.(validation.Errors) {
			err = errors.AddErrorContext(err, key, value.Error())
		}
		if notInSliceErr != nil {
			err = errors.AddErrorContext(err, "TaggedPorts", notInSliceErr.Error())
		}
	} else {
		if notInSliceErr != nil {
			err = errors.Validation.New(errors.ValidationErrorMessage)
			err = errors.AddErrorContext(err, "TaggedPorts", notInSliceErr.Error())
		}
	}
	return err
}
