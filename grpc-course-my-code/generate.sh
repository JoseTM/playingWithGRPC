protoc greet/greetpb/greet.proto  --go_out=plugins=grpc:greet/greetpb/
protoc calculator/calculatorpb/calculatorpb.proto --go_out=plugins=grpc:calculator/calculatorpb