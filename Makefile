APP_EXECUTABLE="out/nabitu-service"

serve-api:
	export ENVIRONMENT=local && \
	go run main.go serve-rest-api

serve-scheduler:
	export ENVIRONMENT=local && \
	go run main.go serve-scheduler

serve-worker:
	export ENVIRONMENT=local && \
	go run main.go serve-worker

serve-task-queue-worker:
	export ENVIRONMENT=local && \
	go run main.go serve-task-queue-worker

compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)

test:
	go test -v ./...

test-race:
	go test -race -v ./...
test-failfast:
	go test -failfast -v ./...
test-coverage:
	./coverage.sh;

static-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

lint:
	@golangci-lint run -E gofmt

format:
	@$(MAKE) fmt
	@$(MAKE) imports

fmt:
	@echo "Formatting code style..."
	gofmt -w -s cmd/.. \
		internal/..

	@echo "[DONE] Formatting code style..."

imports:
	@echo "Formatting imports..."
	goimports -w -local go-core-nabitu \
		cmd/.. \
		config/.. \
		internal/.. \
		pkg/..
	@echo "[DONE] Formatting imports..."

migrate-up:
	@echo "Running migration up..."
	goose --dir=gen/migrations postgres "user=$(user) password=$(password) dbname=$(dbname) host=$(host) port=$(port) sslmode=disable" up
	@echo $(user) $(password) $(dbname)
	@echo "[DONE] Running migration up.."

migrate-down:
	@echo "Running migration down..."
	goose --dir=gen/migrations postgres "user=$(user) password=$(password) dbname=$(dbname) host=$(host) port=$(port) sslmode=disable" down
	@echo $(user) $(password) $(dbname)
	@echo "[DONE] Running migration down.."

gen-mocks:
	@echo "  >  Rebuild Mocking..."
	mockgen --package mocks --source=pkg/storage/storage.go --destination=pkg/storage/mocks/mock_storage.go

	# driver
	@echo "  >  Driver Mocking..."
	mockgen  --package mocks --source=pkg/http_client/http/http_client.go --destination=pkg/http_client/http/mocks/mock_http_client.go

	# common
	@echo "  >  Common Mocking..."
	mockgen  --package mocks --source=internal/common/privy/domain/repository.go --destination=internal/common/privy/domain/mocks/mock_repository.go
	mockgen  --package mocks --source=internal/common/whatsapp/domain/repository.go --destination=internal/common/whatsapp/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/storage/domain/repository.go --destination=internal/common/storage/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/bsi_nabitu/domain/repository.go --destination=internal/common/bsi_nabitu/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/mailer/domain/repository.go --destination=internal/common/mailer/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/notification/domain/repository.go --destination=internal/common/notification/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/content/domain/repository.go --destination=internal/common/content/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/user/domain/repository.go --destination=internal/common/user/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/pdfconverter/domain/repository.go --destination=internal/common/pdfconverter/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/media/domain/repository.go --destination=internal/common/media/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/bi/domain/repository.go --destination=internal/common/bi/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/bi/domain/repository.go --destination=internal/common/bi/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/auth/domain/repository.go --destination=internal/common/auth/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/document/domain/repository.go --destination=internal/common/document/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/queue/domain/repository.go --destination=internal/common/queue/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/authtoken/domain/repository.go --destination=internal/common/authtoken/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/cache/domain/repository.go --destination=internal/common/cache/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/broadcastlog/domain/repository.go --destination=internal/common/broadcastlog/domain/mocks/mock_repository.go


	# modul
	@echo "  >  Usecase and Repository Mocking..."
	mockgen --package mocks --source=internal/module/member/auth/domain/repository.go --destination=internal/module/member/auth/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/auth/domain/usecase.go --destination=internal/module/member/auth/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/otp/domain/repository.go --destination=internal/module/member/otp/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/otp/domain/usecase.go --destination=internal/module/member/otp/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/register_individual/domain/usecase.go --destination=internal/module/member/register_individual/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/register_individual/domain/repository.go --destination=internal/module/member/register_individual/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/register_individual_global/domain/usecase.go --destination=internal/module/member/register_individual_global/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/register_individual_global/domain/repository.go --destination=internal/module/member/register_individual_global/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/register/domain/repository.go --destination=internal/module/member/register/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/register/domain/usecase.go --destination=internal/module/member/register/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/register_issuer/domain/usecase.go --destination=internal/module/member/register_issuer/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/register_issuer/domain/repository.go --destination=internal/module/member/register_issuer/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/referral/domain/usecase.go --destination=internal/module/member/referral/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/referral/domain/repository.go --destination=internal/module/member/referral/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/register_company/domain/usecase.go --destination=internal/module/member/register_company/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/register_company/domain/repository.go --destination=internal/module/member/register_company/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/user_setting/domain/repository.go --destination=internal/module/member/user_setting/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/user_setting/domain/usecase.go --destination=internal/module/member/user_setting/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/dashboard/domain/repository.go --destination=internal/module/member/dashboard/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/dashboard/domain/usecase.go --destination=internal/module/member/dashboard/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/portfolio/domain/repository.go --destination=internal/module/member/portfolio/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/portfolio/domain/usecase.go --destination=internal/module/member/portfolio/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/privy_registration/domain/repository.go --destination=internal/module/member/privy_registration/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/privy_registration/domain/usecase.go --destination=internal/module/member/privy_registration/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/expired_order/domain/usecase.go --destination=internal/module/member/expired_order/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/expired_order/domain/repository.go --destination=internal/module/member/expired_order/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/order_mutation/domain/repository.go --destination=internal/module/member/order_mutation/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/order_mutation/domain/usecase.go --destination=internal/module/member/order_mutation/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/support/domain/repository.go --destination=internal/module/member/support/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/support/domain/usecase.go --destination=internal/module/member/support/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/exchange/domain/repository.go --destination=internal/module/member/exchange/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/exchange/domain/usecase.go --destination=internal/module/member/exchange/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/project/domain/repository.go --destination=internal/module/member/project/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/project/domain/usecase.go --destination=internal/module/member/project/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/project_submission/domain/repository.go --destination=internal/module/member/project_submission/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/project_submission/domain/usecase.go --destination=internal/module/member/project_submission/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/payment_method/domain/repository.go --destination=internal/module/member/payment_method/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/payment_method/domain/usecase.go --destination=internal/module/member/payment_method/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/issuer_dashboard/domain/repository.go --destination=internal/module/member/issuer_dashboard/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/issuer_dashboard/domain/usecase.go --destination=internal/module/member/issuer_dashboard/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/document/domain/usecase.go --destination=internal/module/member/document/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/document/domain/repository.go --destination=internal/module/member/document/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/member/rate/domain/usecase.go --destination=internal/module/member/rate/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/member/rate/domain/repository.go --destination=internal/module/member/rate/domain/mocks/mock_repository.go

	# admin
	mockgen --package mocks --source=internal/module/admin/dashboard/domain/repository.go --destination=internal/module/admin/dashboard/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/dashboard/domain/usecase.go --destination=internal/module/admin/dashboard/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project/domain/repository.go --destination=internal/module/admin/project/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project/domain/usecase.go --destination=internal/module/admin/project/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/team/domain/repository.go --destination=internal/module/admin/team/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/team/domain/usecase.go --destination=internal/module/admin/team/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/capitalpotential/domain/repository.go --destination=internal/module/admin/capitalpotential/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/capitalpotential/domain/usecase.go --destination=internal/module/admin/capitalpotential/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/settings/domain/repository.go --destination=internal/module/admin/settings/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/settings/domain/usecase.go --destination=internal/module/admin/settings/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/support/domain/repository.go --destination=internal/module/admin/support/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/support/domain/usecase.go --destination=internal/module/admin/support/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/activity/domain/repository.go --destination=internal/module/admin/activity/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/activity/domain/usecase.go --destination=internal/module/admin/activity/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member/domain/repository.go --destination=internal/module/admin/member/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member/domain/usecase.go --destination=internal/module/admin/member/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_login_activity/domain/repository.go --destination=internal/module/admin/member_login_activity/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_login_activity/domain/usecase.go --destination=internal/module/admin/member_login_activity/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_profile/domain/repository.go --destination=internal/module/admin/member_profile/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_profile/domain/usecase.go --destination=internal/module/admin/member_profile/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_apuppt/domain/repository.go --destination=internal/module/admin/member_apuppt/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_apuppt/domain/service.go --destination=internal/module/admin/member_apuppt/domain/mocks/mock_service.go
	mockgen --package mocks --source=internal/module/admin/member_business/domain/repository.go --destination=internal/module/admin/member_business/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_business/domain/usecase.go --destination=internal/module/admin/member_business/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/broadcast/domain/repository.go --destination=internal/module/admin/broadcast/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/broadcast/domain/usecase.go --destination=internal/module/admin/broadcast/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_bank/domain/repository.go --destination=internal/module/admin/member_bank/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_bank/domain/service.go --destination=internal/module/admin/member_bank/domain/mocks/mock_service.go
	mockgen --package mocks --source=internal/module/admin/member_referral/domain/repository.go --destination=internal/module/admin/member_referral/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_referral/domain/usecase.go --destination=internal/module/admin/member_referral/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_address/domain/repository.go --destination=internal/module/admin/member_address/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_address/domain/service.go --destination=internal/module/admin/member_address/domain/mocks/mock_service.go
	mockgen --package mocks --source=internal/module/admin/member_portfolio/domain/repository.go --destination=internal/module/admin/member_portfolio/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_portfolio/domain/service.go --destination=internal/module/admin/member_portfolio/domain/mocks/mock_service.go
	mockgen --package mocks --source=internal/module/admin/broadcastlog/domain/repository.go --destination=internal/module/admin/broadcastlog/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/broadcastlog/domain/usecase.go --destination=internal/module/admin/broadcastlog/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projecteditprofit/domain/repository.go --destination=internal/module/admin/projecteditprofit/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projecteditprofit/domain/usecase.go --destination=internal/module/admin/projecteditprofit/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projectaddinstallment/domain/repository.go --destination=internal/module/admin/projectaddinstallment/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projectaddinstallment/domain/usecase.go --destination=internal/module/admin/projectaddinstallment/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projectinstallment/domain/repository.go --destination=internal/module/admin/projectinstallment/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projectinstallment/domain/usecase.go --destination=internal/module/admin/projectinstallment/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_stage_change/domain/repository.go --destination=internal/module/admin/member_stage_change/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_stage_change/domain/usecase.go --destination=internal/module/admin/member_stage_change/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/member_project_submission/domain/repository.go --destination=internal/module/admin/member_project_submission/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/member_project_submission/domain/usecase.go --destination=internal/module/admin/member_project_submission/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projectinstallmentbyphase/domain/repository.go --destination=internal/module/admin/projectinstallmentbyphase/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projectinstallmentbyphase/domain/usecase.go --destination=internal/module/admin/projectinstallmentbyphase/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projectinstallmentupdatedatebyphase/domain/repository.go --destination=internal/module/admin/projectinstallmentupdatedatebyphase/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projectinstallmentupdatedatebyphase/domain/usecase.go --destination=internal/module/admin/projectinstallmentupdatedatebyphase/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project_simulation/domain/repository.go --destination=internal/module/admin/project_simulation/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project_simulation/domain/usecase.go --destination=internal/module/admin/project_simulation/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project_referrer/domain/repository.go --destination=internal/module/admin/project_referrer/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project_referrer/domain/usecase.go --destination=internal/module/admin/project_referrer/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project_document/domain/repository.go --destination=internal/module/admin/project_document/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project_document/domain/usecase.go --destination=internal/module/admin/project_document/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/users/domain/repository.go --destination=internal/module/admin/users/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/users/domain/usecase.go --destination=internal/module/admin/users/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project_refund/domain/repository.go --destination=internal/module/admin/project_refund/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project_refund/domain/usecase.go --destination=internal/module/admin/project_refund/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project_report/domain/repository.go --destination=internal/module/admin/project_report/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project_report/domain/usecase.go --destination=internal/module/admin/project_report/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/project_final_report/domain/repository.go --destination=internal/module/admin/project_final_report/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/project_final_report/domain/usecase.go --destination=internal/module/admin/project_final_report/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projectbroadcast/domain/repository.go --destination=internal/module/admin/projectbroadcast/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projectbroadcast/domain/usecase.go --destination=internal/module/admin/projectbroadcast/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/category/domain/repository.go --destination=internal/module/admin/category/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/category/domain/usecase.go --destination=internal/module/admin/category/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/auth/domain/repository.go --destination=internal/module/admin/auth/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/auth/domain/usecase.go --destination=internal/module/admin/auth/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/userinfo/domain/repository.go --destination=internal/module/admin/userinfo/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/userinfo/domain/usecase.go --destination=internal/module/admin/userinfo/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/notifyavailableslot/domain/repository.go --destination=internal/module/admin/notifyavailableslot/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/notifyavailableslot/domain/usecase.go --destination=internal/module/admin/notifyavailableslot/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/admin/projectstatus/domain/repository.go --destination=internal/module/admin/projectstatus/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/admin/projectstatus/domain/usecase.go --destination=internal/module/admin/projectstatus/domain/mocks/mock_usecase.go
