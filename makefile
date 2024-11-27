.PHONY: tools testgen

tools:
	go build -o ./bin/mvm.exe ./mvmtools/ 

testgen:
	@rm -rf ./examples/test-views/* 
	@go build -o ./bin/mvm.exe ./mvmtools/ 
	@.\bin\mvm.exe gen main, controller , view, model --name login -o .\examples\test-views\ && go run .\examples\test-views\
