#####################################################################
# Rules for local environment                                       #
#####################################################################

kind-% kind2-% gardener-%: export IPFAMILY := $(IPFAMILY)
# KUBECONFIG
kind-up kind-down gardener-up gardener-dev gardener-debug gardener-down: export KUBECONFIG = $(GARDENER_LOCAL_KUBECONFIG)
test-e2e-local-simple test-e2e-local-migration test-e2e-local-workerless test-e2e-local ci-e2e-kind ci-e2e-kind-upgrade: export KUBECONFIG = $(GARDENER_LOCAL_KUBECONFIG)
kind2-up kind2-down gardenlet-kind2-up gardenlet-kind2-dev gardenlet-kind2-debug gardenlet-kind2-down: export KUBECONFIG = $(GARDENER_LOCAL2_KUBECONFIG)
kind-extensions-up kind-extensions-down gardener-extensions-up gardener-extensions-down: export KUBECONFIG = $(GARDENER_EXTENSIONS_KUBECONFIG)
kind-ha-single-zone-up kind-ha-single-zone-down gardener-ha-single-zone-up gardener-ha-single-zone-down: export KUBECONFIG = $(GARDENER_LOCAL_HA_SINGLE_ZONE_KUBECONFIG)
test-e2e-local-ha-single-zone test-e2e-local-migration-ha-single-zone ci-e2e-kind-ha-single-zone ci-e2e-kind-ha-single-zone-upgrade: export KUBECONFIG = $(GARDENER_LOCAL_HA_SINGLE_ZONE_KUBECONFIG)
kind2-ha-single-zone-up kind2-ha-single-zone-down gardenlet-kind2-ha-single-zone-up gardenlet-kind2-ha-single-zone-dev gardenlet-kind2-ha-single-zone-debug gardenlet-kind2-ha-single-zone-down: export KUBECONFIG = $(GARDENER_LOCAL2_HA_SINGLE_ZONE_KUBECONFIG)
kind-ha-multi-zone-up kind-ha-multi-zone-down gardener-ha-multi-zone-up: export KUBECONFIG = $(GARDENER_LOCAL_HA_MULTI_ZONE_KUBECONFIG)
test-e2e-local-ha-multi-zone ci-e2e-kind-ha-multi-zone ci-e2e-kind-ha-multi-zone-upgrade: export KUBECONFIG = $(GARDENER_LOCAL_HA_MULTI_ZONE_KUBECONFIG)
kind-operator-up kind-operator-down operator-up operator-dev operator-debug operator-down test-e2e-local-operator ci-e2e-kind-operator: export KUBECONFIG = $(GARDENER_LOCAL_OPERATOR_KUBECONFIG)
# CLUSTER_NAME
kind-up kind-down: export CLUSTER_NAME = gardener-local
kind2-up kind2-down: export CLUSTER_NAME = gardener-local2
kind-ha-single-zone-up kind-ha-single-zone-down: export CLUSTER_NAME = gardener-local-ha-single-zone
kind2-ha-single-zone-up kind2-ha-single-zone-down: export CLUSTER_NAME = gardener-local2-ha-single-zone
kind-ha-multi-zone-up kind-ha-multi-zone-down: export CLUSTER_NAME = gardener-local-ha-multi-zone
# KIND_KUBECONFIG
kind-up kind-down: export KIND_KUBECONFIG = $(KIND_ROOT)/example/provider-local/seed-kind/base/kubeconfig
kind2-up kind2-down: export KIND_KUBECONFIG = $(KIND_ROOT)/example/provider-local/seed-kind2/base/kubeconfig
kind-ha-single-zone-up kind-ha-single-zone-down: export KIND_KUBECONFIG = $(KIND_ROOT)/example/provider-local/seed-kind-ha-single-zone/base/kubeconfig
kind2-ha-single-zone-up kind2-ha-single-zone-down: export KIND_KUBECONFIG = $(KIND_ROOT)/example/provider-local/seed-kind2-ha-single-zone/base/kubeconfig
kind-ha-multi-zone-up kind-ha-multi-zone-down: export KIND_KUBECONFIG = $(KIND_ROOT)/example/provider-local/seed-kind-ha-multi-zone/base/kubeconfig
# CLUSTER_VALUES
kind-up kind-down: export CLUSTER_VALUES = $(KIND_ROOT)/example/gardener-local/kind/local/values.yaml
kind2-up kind2-down: export CLUSTER_VALUES = $(KIND_ROOT)/example/gardener-local/kind/local2/values.yaml
kind-ha-single-zone-up kind-ha-single-zone-down: export CLUSTER_VALUES = $(KIND_ROOT)/example/gardener-local/kind/ha-single-zone/values.yaml
kind2-ha-single-zone-up kind2-ha-single-zone-down: export CLUSTER_VALUES = $(KIND_ROOT)/example/gardener-local/kind/ha-single-zone2/values.yaml
kind-ha-multi-zone-up kind-ha-multi-zone-down: export CLUSTER_VALUES = $(KIND_ROOT)/example/gardener-local/kind/ha-multi-zone/values.yaml
# ADDITIONAL_PARAMETERS
kind2-up kind2-ha-single-zone-up: export ADDITIONAL_PARAMETERS = --skip-registry
kind2-down: export ADDITIONAL_PARAMETERS = --keep-backupbuckets-dir
kind-ha-multi-zone-up: export ADDITIONAL_PARAMETERS = --multi-zonal

kind-up kind2-up kind-ha-single-zone-up kind2-ha-single-zone-up kind-ha-multi-zone-up: $(KIND) $(KUBECTL) $(HELM) $(YQ) $(KIND_ROOT)/example
	$(HACK_DIR)/kind-up.sh --cluster-name $(CLUSTER_NAME) --path-kubeconfig $(KIND_KUBECONFIG) --path-cluster-values $(CLUSTER_VALUES) $(ADDITIONAL_PARAMETERS)
kind-down kind2-down kind-ha-single-zone-down kind2-ha-single-zone-down kind-ha-multi-zone-down: $(KIND) $(KIND_ROOT)/example
	$(HACK_DIR)/kind-down.sh --cluster-name $(CLUSTER_NAME) --path-kubeconfig $(KIND_KUBECONFIG) $(ADDITIONAL_PARAMETERS)

kind-extensions-up: $(KIND) $(KUBECTL)
	REPO_ROOT=$(REPO_ROOT) $(HACK_DIR)/kind-extensions-up.sh
kind-extensions-down: $(KIND)
	docker stop gardener-extensions-control-plane
kind-extensions-clean:
	$(HACK_DIR)/kind-down.sh --cluster-name gardener-extensions --path-kubeconfig $(REPO_ROOT)/example/provider-extensions/garden/kubeconfig

kind-operator-up: $(KIND) $(KUBECTL) $(HELM) $(YQ)
	$(HACK_DIR)/kind-up.sh --cluster-name gardener-operator-local --path-kubeconfig $(REPO_ROOT)/example/gardener-local/kind/operator/kubeconfig --path-cluster-values $(REPO_ROOT)/example/gardener-local/kind/operator/values.yaml
	mkdir -p $(REPO_ROOT)/dev/local-backupbuckets/gardener-operator
kind-operator-down: $(KIND)
	$(HACK_DIR)/kind-down.sh --cluster-name gardener-operator-local --path-kubeconfig $(REPO_ROOT)/example/gardener-local/kind/operator/kubeconfig
	# We need root privileges to clean the backup bucket directory, see https://github.com/gardener/gardener/issues/6752
	docker run --user root:root -v $(REPO_ROOT)/dev/local-backupbuckets:/dev/local-backupbuckets alpine rm -rf /dev/local-backupbuckets/gardener-operator


# speed-up skaffold deployments by building all images concurrently
export SKAFFOLD_BUILD_CONCURRENCY = 0
export SKAFFOLD_FILENAME := $(REPO_ROOT)/skaffold/skaffold.yaml
gardener%up gardener%dev gardener%debug gardenlet%up gardenlet%dev gardenlet%debug operator-up operator-dev operator-debug: export SKAFFOLD_DEFAULT_REPO = localhost:5001
gardener%up gardener%dev gardener%debug gardenlet%up gardenlet%dev gardenlet%debug operator-up operator-dev operator-debug: export SKAFFOLD_PUSH = true
gardener%up gardener%dev gardener%debug gardenlet%up gardenlet%dev gardenlet%debug operator-up operator-dev operator-debug: export SOURCE_DATE_EPOCH = $(shell date +%s)
# use static label for skaffold to prevent rolling all gardener components on every `skaffold` invocation
gardener%up gardener%dev gardener%debug gardener%down gardenlet%up gardenlet%dev gardenlet%debug gardenlet%down: export SKAFFOLD_LABEL = skaffold.dev/run-id=gardener-local
# skaffold dev and debug clean up deployed modules by default, disable this
gardener%dev gardener%debug gardenlet%dev gardenlet%debug operator-dev operator-debug: export SKAFFOLD_CLEANUP = false
# skaffold dev triggers new builds and deployments immediately on file changes by default,
# this is too heavy in a large project like gardener, so trigger new builds and deployments manually instead.
gardener%dev gardenlet%dev operator-dev: export SKAFFOLD_TRIGGER = manual
# Artifacts might be already built when you decide to start debugging.
# However, these artifacts do not include the gcflags which `skaffold debug` sets automatically, so delve would not work.
# Disabling the skaffold cache for debugging ensures that you run artifacts with gcflags required for debugging.
gardener%debug gardenlet%debug operator-debug: export SKAFFOLD_CACHE_ARTIFACTS = false
# Health checks may be disabled on a per-controller basis. To boot a garden from scratch, some controllers may be launched before
# their attendant resources are installed. Those controllers will need to be restarted by health checks after the resources are available.
# Developers who prefer to do this manually may provide comma-separated controller names from {admission,controller,scheduler,gardenlet}.
gardener%dev gardener%debug gardenlet%dev gardenlet%debug operator-dev operator-debug: export SKAFFOLD_DISABLE_HEALTH_CHECKS := ""

gardener-ha-single-zone%: export SKAFFOLD_PROFILE=ha-single-zone
gardener-ha-multi-zone%: export SKAFFOLD_PROFILE=ha-multi-zone

gardener-up gardener-ha-single-zone-up gardener-ha-multi-zone-up: $(SKAFFOLD) $(HELM) $(KUBECTL) $(YQ)
	env | grep SKAFFOLD_ | tr '\n' ';'
	$(SKAFFOLD) run
gardener-dev gardener-ha-single-zone-dev gardener-ha-multi-zone-dev: $(SKAFFOLD) $(HELM) $(KUBECTL) $(YQ)
	env | grep SKAFFOLD_ | tr '\n' ';'
	$(SKAFFOLD) dev
gardener-debug gardener-ha-single-zone-debug gardener-ha-multi-zone-debug: $(SKAFFOLD) $(HELM) $(KUBECTL) $(YQ)
	env | grep SKAFFOLD_ | tr '\n' ';'
	$(SKAFFOLD) debug
gardener-down gardener-ha-single-zone-down gardener-ha-multi-zone-down: $(SKAFFOLD) $(HELM) $(KUBECTL)
	./hack/gardener-down.sh

gardener-extensions-%: export SKAFFOLD_LABEL = skaffold.dev/run-id=gardener-extensions

gardener-extensions-up: $(SKAFFOLD) $(HELM) $(KUBECTL) $(YQ)
	$(HACK_DIR)/gardener-extensions-up.sh --path-garden-kubeconfig $(REPO_ROOT)/example/provider-extensions/garden/kubeconfig --path-seed-kubeconfig $(SEED_KUBECONFIG) --seed-name $(SEED_NAME)
gardener-extensions-down: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(HACK_DIR)/gardener-extensions-down.sh --path-garden-kubeconfig $(REPO_ROOT)/example/provider-extensions/garden/kubeconfig --path-seed-kubeconfig $(SEED_KUBECONFIG) --seed-name $(SEED_NAME)

gardenlet-kind2-up gardenlet-kind2-dev gardenlet-kind2-debug gardenlet-kind2-down: export SKAFFOLD_PREFIX_NAME = kind2
gardenlet-kind2-ha-single-zone-up gardenlet-kind2-ha-single-zone-dev gardenlet-kind2-ha-single-zone-debug gardenlet-kind2-ha-single-zone-down: export SKAFFOLD_PREFIX_NAME = kind2-ha-single-zone
gardenlet-kind2-up gardenlet-kind2-dev gardenlet-kind2-debug gardenlet-kind2-down: export SKAFFOLD_COMMAND_KUBECONFIG := $(GARDENER_LOCAL_KUBECONFIG)
gardenlet-kind2-ha-single-zone-up gardenlet-kind2-ha-single-zone-dev gardenlet-kind2-ha-single-zone-debug gardenlet-kind2-ha-single-zone-down: export SKAFFOLD_COMMAND_KUBECONFIG := $(GARDENER_LOCAL_HA_SINGLE_ZONE_KUBECONFIG)

gardenlet-kind2-up gardenlet-kind2-ha-single-zone-up: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) deploy -m $(SKAFFOLD_PREFIX_NAME)-env -p $(SKAFFOLD_PREFIX_NAME) --kubeconfig=$(SKAFFOLD_COMMAND_KUBECONFIG)
	@# define GARDENER_LOCAL_KUBECONFIG so that it can be used by skaffold when checking whether the seed managed by this gardenlet is ready
	GARDENER_LOCAL_KUBECONFIG=$(SKAFFOLD_COMMAND_KUBECONFIG) $(SKAFFOLD) run -m gardenlet -p $(SKAFFOLD_PREFIX_NAME)
gardenlet-kind2-dev gardenlet-kind2-ha-single-zone-dev: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) deploy -m $(SKAFFOLD_PREFIX_NAME)-env -p $(SKAFFOLD_PREFIX_NAME) --kubeconfig=$(SKAFFOLD_COMMAND_KUBECONFIG)
	@# define GARDENER_LOCAL_KUBECONFIG so that it can be used by skaffold when checking whether the seed managed by this gardenlet is ready
	GARDENER_LOCAL_KUBECONFIG=$(SKAFFOLD_COMMAND_KUBECONFIG) $(SKAFFOLD) dev -m gardenlet -p $(SKAFFOLD_PREFIX_NAME)
gardenlet-kind2-debug gardenlet-kind2-ha-single-zone-debug: $(SKAFFOLD) $(HELM)
	$(SKAFFOLD) deploy -m $(SKAFFOLD_PREFIX_NAME)-env -p $(SKAFFOLD_PREFIX_NAME) --kubeconfig=$(SKAFFOLD_COMMAND_KUBECONFIG)
	@# define GARDENER_LOCAL_KUBECONFIG so that it can be used by skaffold when checking whether the seed managed by this gardenlet is ready
	GARDENER_LOCAL_KUBECONFIG=$(SKAFFOLD_COMMAND_KUBECONFIG) $(SKAFFOLD) debug -m gardenlet -p $(SKAFFOLD_PREFIX_NAME)
gardenlet-kind2-down gardenlet-kind2-ha-single-zone-down: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) delete -m $(SKAFFOLD_PREFIX_NAME)-env -p $(SKAFFOLD_PREFIX_NAME) --kubeconfig=$(SKAFFOLD_COMMAND_KUBECONFIG)
	$(SKAFFOLD) delete -m gardenlet,$(SKAFFOLD_PREFIX_NAME)-env -p $(SKAFFOLD_PREFIX_NAME)

operator-%: export SKAFFOLD_FILENAME = $(REPO_ROOT)/skaffold/skaffold-operator.yaml

operator-up: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) run
operator-dev: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) dev
operator-debug: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) debug
operator-down: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(KUBECTL) annotate garden --all confirmation.gardener.cloud/deletion=true
	$(KUBECTL) delete garden --all --ignore-not-found --wait --timeout 5m
	$(SKAFFOLD) delete
