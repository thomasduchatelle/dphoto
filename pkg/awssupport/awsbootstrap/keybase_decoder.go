package awsbootstrap

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
)

func KeybaseDecode(b64Secret string) (string, error) {
	cmd := fmt.Sprintf("echo '%s' | base64 --decode | keybase pgp decrypt", b64Secret)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return string(out), errors.Wrapf(err, "keybase failed to decrypt the secret key [%s]", out)
}
