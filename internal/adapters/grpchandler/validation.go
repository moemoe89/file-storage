package grpchandler

import "fmt"

// validateContentType validates the file Content-Type with expected values from the request.
func validateContentType(targetTypes []string, contentType string) error {
	if len(targetTypes) == 0 {
		return nil
	}

	matchType := false

	for _, targetType := range targetTypes {
		if targetType == contentType {
			matchType = true
			break
		}
	}

	if !matchType {
		return fmt.Errorf("file mime type: %s doesn't match with the target: %v", contentType, targetTypes)
	}

	return nil
}

// validateSize validates the file size with expected max size from the request.
func validateSize(serverOffset, clientOffset, targetSize int64) error {
	if clientOffset != serverOffset {
		return fmt.Errorf("unexpected offset, got %d, want %d", serverOffset, clientOffset)
	}

	if targetSize > 0 && serverOffset > targetSize {
		return fmt.Errorf("file size: %d exceed the target: %d", serverOffset, targetSize)
	}

	return nil
}
