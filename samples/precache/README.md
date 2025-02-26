# Pre-caching manual test #
This directory contains samples for manual standalone pre-caching testing.
All the commands provided here must be executed from the main project directory with KUBECONFIG environment variable initialized and pointing to the ACM hub cluster.
## Mocks and workarounds ##
1. Deploy the CRD from [config/crd/bases/ran.openshift.io_clustergroupupgrades.yaml](../../config/crd/bases/ran.openshift.io_clustergroupupgrades.yaml)
2. If the CSV object is not present on the hub, it must be mocked. 
3. Copy the autogenerated CSV from [bundle/manifests/cluster-group-upgrades-operator.clusterserviceversion.yaml](../../bundle/manifests/cluster-group-upgrades-operator.clusterserviceversion.yaml). 
4. Change the namespace to `openshift-operators` and adjust the pre-caching image pull spec to your desired location. 
5. Deploy the mock policies and objects by
```bash
oc apply -k samples/precache
```
### Note on catalogsource policy ###
When working with a disconnected registry, cluster administrators will have to create such a policy for the cluster to function. For connected environments, spoke clusters are preconfigured with default catalog sources (such as redhat-operators), and no catalog source policy will be necessary. 
Precaching, however, requires a catalog source policy to be explicitly configured on the hub cluster. The example catalog source policy for working with default registry is provided in [catsrc-policy.yaml](catsrc-policy.yaml)

## Running the controller ##
```bash
make run
```
