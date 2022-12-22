// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package whirlpool

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// SetProtocolFeeRate is the `setProtocolFeeRate` instruction.
type SetProtocolFeeRate struct {
	ProtocolFeeRate *uint16

	// [0] = [] whirlpoolsConfig
	//
	// [1] = [WRITE] whirlpool
	//
	// [2] = [SIGNER] feeAuthority
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewSetProtocolFeeRateInstructionBuilder creates a new `SetProtocolFeeRate` instruction builder.
func NewSetProtocolFeeRateInstructionBuilder() *SetProtocolFeeRate {
	nd := &SetProtocolFeeRate{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetProtocolFeeRate sets the "protocolFeeRate" parameter.
func (inst *SetProtocolFeeRate) SetProtocolFeeRate(protocolFeeRate uint16) *SetProtocolFeeRate {
	inst.ProtocolFeeRate = &protocolFeeRate
	return inst
}

// SetWhirlpoolsConfigAccount sets the "whirlpoolsConfig" account.
func (inst *SetProtocolFeeRate) SetWhirlpoolsConfigAccount(whirlpoolsConfig ag_solanago.PublicKey) *SetProtocolFeeRate {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(whirlpoolsConfig)
	return inst
}

// GetWhirlpoolsConfigAccount gets the "whirlpoolsConfig" account.
func (inst *SetProtocolFeeRate) GetWhirlpoolsConfigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetWhirlpoolAccount sets the "whirlpool" account.
func (inst *SetProtocolFeeRate) SetWhirlpoolAccount(whirlpool ag_solanago.PublicKey) *SetProtocolFeeRate {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(whirlpool).WRITE()
	return inst
}

// GetWhirlpoolAccount gets the "whirlpool" account.
func (inst *SetProtocolFeeRate) GetWhirlpoolAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetFeeAuthorityAccount sets the "feeAuthority" account.
func (inst *SetProtocolFeeRate) SetFeeAuthorityAccount(feeAuthority ag_solanago.PublicKey) *SetProtocolFeeRate {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(feeAuthority).SIGNER()
	return inst
}

// GetFeeAuthorityAccount gets the "feeAuthority" account.
func (inst *SetProtocolFeeRate) GetFeeAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

func (inst SetProtocolFeeRate) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_SetProtocolFeeRate,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetProtocolFeeRate) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetProtocolFeeRate) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.ProtocolFeeRate == nil {
			return errors.New("ProtocolFeeRate parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.WhirlpoolsConfig is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Whirlpool is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.FeeAuthority is not set")
		}
	}
	return nil
}

func (inst *SetProtocolFeeRate) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetProtocolFeeRate")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("ProtocolFeeRate", *inst.ProtocolFeeRate))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("whirlpoolsConfig", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("       whirlpool", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("    feeAuthority", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (obj SetProtocolFeeRate) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `ProtocolFeeRate` param:
	err = encoder.Encode(obj.ProtocolFeeRate)
	if err != nil {
		return err
	}
	return nil
}
func (obj *SetProtocolFeeRate) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `ProtocolFeeRate`:
	err = decoder.Decode(&obj.ProtocolFeeRate)
	if err != nil {
		return err
	}
	return nil
}

// NewSetProtocolFeeRateInstruction declares a new SetProtocolFeeRate instruction with the provided parameters and accounts.
func NewSetProtocolFeeRateInstruction(
	// Parameters:
	protocolFeeRate uint16,
	// Accounts:
	whirlpoolsConfig ag_solanago.PublicKey,
	whirlpool ag_solanago.PublicKey,
	feeAuthority ag_solanago.PublicKey) *SetProtocolFeeRate {
	return NewSetProtocolFeeRateInstructionBuilder().
		SetProtocolFeeRate(protocolFeeRate).
		SetWhirlpoolsConfigAccount(whirlpoolsConfig).
		SetWhirlpoolAccount(whirlpool).
		SetFeeAuthorityAccount(feeAuthority)
}
