package mode_s
import "time"

func Fuzz(data []byte) int {
	_, err := DecodeString(string(data), time.Now());
	if err == nil {
		return 1;
	}

	return 0;
}
