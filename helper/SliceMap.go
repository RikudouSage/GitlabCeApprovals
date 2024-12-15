package helper

func SliceMap[TInput any, TOutput any](input []TInput, callback func(item TInput) TOutput) []TOutput {
	output := make([]TOutput, 0)
	for _, value := range input {
		output = append(output, callback(value))
	}

	return output
}
