package promptx

// import buffer package
//go:generate gogen import ./buffer -t Buffer -t Document -o gen_pkg_buffer.go

// import input package
//go:generate gogen import ./input -t WinSize -t ConsoleParser -t Key -t ASCIICode -o gen_pkg_input.go

// import output package
//go:generate gogen import ./output -t Color -t ConsoleWriter -o gen_pkg_output.go

// import completion
//go:generate gogen import ./completion -t Suggest -t CompletionManager -t Filter -o gen_pkg_completion.go
