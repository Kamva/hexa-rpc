package hrpc

import "github.com/kamva/hexa"

// ErrDetailsFromErrList converts list of errors to ErrorDetails.
func ErrDetailsFromErrList(ctx hexa.Context, errors []error) []*ErrorDetails {
	errorsPb := make([]*ErrorDetails, len(errors))
	for i, err := range errors {
		if err != nil {
			errorsPb[i] = NewErrorDetailsFromRawError(ctx, err)
		}
	}

	return errorsPb
}
