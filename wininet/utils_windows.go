package wininet

import "github.com/mjwhitta/errors"

func convertFail(str string, e error) error {
	return errors.Newf(
		"failed to convert %s to Windows type: %w",
		str,
		e,
	)
}
