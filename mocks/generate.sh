#!/usr/bin/env bash
set -x
mockgen -source ../company/company_repository.go -destination mock_company/mock_company_repository.go -package mock_company_repository
mockgen -source ../company/company_service.go -destination mock_company/mock_company_service.go -package mock_company_service
git add .