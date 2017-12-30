package validate

//go:generate genny -in=comparable.gen -out=comparable.go gen "T1=float64,int64,time.Duration"
//go:generate genny -in=validations_with_args.gen -out=validations_with_args.go gen "T1=BUILTINS,interface{} T2=bool,string,float64,int64,interface{}"
//go:generate genny -in=validations_simple.gen -out=validations_simple.go gen "T1=BUILTINS,interface{}"
