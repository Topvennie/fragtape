package capture

type withRollbackStruct[T any] struct {
	items  []T           // All items
	do     func(T) error // The function to use on each item
	revert func(T)       // The function called on every parsed item if one fails
}

func withRollback[T any](s withRollbackStruct[T]) error {
	done := make([]T, 0, len(s.items))
	var err error

	for i := range s.items {
		err = s.do(s.items[i])
		if err != nil {
			break
		}

		done = append(done, s.items[i])
	}

	if err != nil {
		for i := range done {
			s.revert(done[i])
		}

		return err
	}

	return nil
}
