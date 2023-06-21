REGISTRY_URL=https://code.db.cafe
PACKAGE_OWNER=pauloo27
FILE=$(word 2, $(MAKECMDGOALS))

.PHONY: pack
pack:
	helm package .

.PHONY: push
push:
	@curl --user ${PACKAGE_OWNER}:${REGISTRY_PASSWORD} -X POST --upload-file ./$(FILE) $(REGISTRY_URL)/api/packages/$(PACKAGE_OWNER)/helm/api/charts
