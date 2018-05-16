format :
	clang-format --style=google -i *.proto
gen :
	protoc -I. --go_out=. *.proto
	mv types_test.pb.go types_test.go
fuzz :
	go-fuzz-build github.com/utilitywarehouse/protoid
	go-fuzz -workdir ./fuzz -bin protoid-fuzz.zip

.PHONY : format gen fuzz
