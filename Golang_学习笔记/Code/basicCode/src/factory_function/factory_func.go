package factory_function

func AddSuffix(suffix string) func(string) string{
	return func (name string) string{
		return name + suffix
	}
}
