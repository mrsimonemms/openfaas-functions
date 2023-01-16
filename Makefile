K3D_CONFIG ?= k3d.config.yaml
REGISTRY_ADDRESS ?= registry.localhost:5000

deps:
	@for chart in $$(ls -d ./charts/*); do \
		helm dependency update $$chart ; \
	done
.PHONY: deps

destroy:
	@skaffold config unset local-cluster
	@skaffold config unset default-repo
	@skaffold config unset kube-context
	@k3d cluster delete --config ${K3D_CONFIG}
.PHONY: destroy

openfaas_login:
	@kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 -d | faas-cli login \
		--username "$(shell kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-user}" | base64 -d)" \
		--password-stdin

	@echo "==="
	@echo "Host:     http://$(shell kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-user}" | base64 -d):$(shell kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 -d)@127.0.0.1:8080"
	@echo "Username: $(shell kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-user}" | base64 -d)"
	@echo "Password: $(shell kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 -d)"
	@echo "==="
.PHONY: openfaas_login

provision:
	@k3d cluster create --config ${K3D_CONFIG} || true
	@kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml
	@skaffold config set local-cluster false
	@skaffold config set default-repo ${REGISTRY_ADDRESS}

	@helm upgrade \
		--atomic \
		--cleanup-on-fail \
		--create-namespace \
		--install \
		--namespace="openfaas" \
		--reset-values \
		--wait \
		openfaas \
		./charts/openfaas/

	@$(MAKE) openfaas_login
.PHONY: provision

serve:
	@skaffold dev
.PHONY: serve

templates:
	@faas-cli template store pull golang-middleware
.PHONY: templates
