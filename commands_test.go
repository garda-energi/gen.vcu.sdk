package sdk

import "testing"

func TestCommandsHelper(t *testing.T) {
	t.Run("getCmdByInvoker", func(t *testing.T) {
		_, err := getCmdByInvoker("Generator")

		wantErr := errCmdNotFound
		if err != wantErr {
			t.Errorf("want %s, got %s", wantErr, err)
		}
	})
}
