package luks

import (
	"github.com/fleetdm/fleet/v4/orbit/pkg/dialog"
)

type KeyEscrower interface {
	SendLinuxKeyEscrowResponse(LuksResponse) error
}

type LuksRunner struct {
	escrower KeyEscrower
	notifier dialog.Dialog
}

type LuksResponse struct {
	// Passphrase is a newly created passphrase generated by fleetd for securing the LUKS volume.
	// This passphrase will be securely escrowed to the server.
	Passphrase string

	// KeySlot specifies the LUKS key slot where this new passphrase was created.
	// It is currently not used, but may be useful in the future for passphrase rotation.
	KeySlot *uint

	// Salt is the salt used to generate the LUKS key.
	Salt string

	// Err is the error message that occurred during the escrow process.
	Err string
}

func New(escrower KeyEscrower, notifier dialog.Dialog) *LuksRunner {
	return &LuksRunner{
		escrower: escrower,
		notifier: notifier,
	}
}
