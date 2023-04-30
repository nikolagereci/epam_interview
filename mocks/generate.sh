#!/usr/bin/env bash
set -x
mockgen -source ../company/company_repository.go -destination mock_company/repository/mock_company_repository.go -package mock_company_repository
mockgen -source ../company/company_service.go -destination mock_company/service/mock_company_service.go -package mock_company_service
mockgen -source ../event/kafka.go -destination mock_company/event/mock_kafka.go -package mock_kafka
git add .