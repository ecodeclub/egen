# 单元测试
ut:
	go test -race ./...

setup:
	sh ./script/setup.sh

# e2e 测试
e2e:
	sh ./script/integrate_test.sh

e2e_up:
	docker compose -f ./script/docker-compose.yml up -d

e2e_down:
	docker compose -f ./script/docker-compose.yml down
