.PHONY: buf-generate
generate:
	buf mod update
	buf generate --template ./buf.gen.yaml --path ./user -o ../types
	buf generate --template ./buf.gen.yaml --path ./config -o ../types
