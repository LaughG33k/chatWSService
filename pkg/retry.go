package pkg

import "time"

func RetrySmth(fn func() error, attempts int, timeWait time.Duration) (err error) {

	for attempts > 0 {

		if err = fn(); err != nil {
			attempts--
			time.Sleep(timeWait)
			continue
		}

		return nil

	}

	return err

}
