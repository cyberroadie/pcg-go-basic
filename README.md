# PCG Random Number Generation, Minimal Go Edition

[PCG-Random website]: http://www.pcg-random.org

This code is a rewrite in Go of the C implementation found here https://github.com/imneme/pcg-c-basic

This code provides a minimal implementation of one member of the PCG family
of random number generators, which are fast, statistically excellent,
and offer a number of useful features.

Full details can be found at the [PCG-Random website].  This version
of the code provides a single family member and skips some useful features
(such as jump-ahead/jump-back) 

## Building and running
go get git@github.com:cyberroadie/pcg-go-basic/pcg32
cd src/github.com/cyberroadie/pcg-go-basic/demo/

        go run pcg32-demo.go

Global initializer for RNG

        go run pcg32-demo.go -global=true 

non deterministic seed

        go run pcg32-demo.go -r=true

