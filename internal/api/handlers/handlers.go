package handlers

type Test struct{}

type HelloArgs struct {
	Name string
}

func (test *Test) Hello(args *HelloArgs, result *string) error {
	*result = "Hello " + args.Name
	// return &jsonrpc2.Error{42, "some issue", []string{"one", "two"}}
	return nil
}
