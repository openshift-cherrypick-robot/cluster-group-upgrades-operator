package controllers

import (
	"errors"

	ranv1alpha1 "github.com/openshift-kni/cluster-group-upgrades-operator/api/v1alpha1"
	"github.com/openshift-kni/cluster-group-upgrades-operator/controllers/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// extractOpenshiftImagePlatformFromPolicies validates that there's ClusterVersion policy, validates the content of ClusterVersion and extracts Image if needed
func (r *ClusterGroupUpgradeReconciler) extractOpenshiftImagePlatformFromPolicies(
	policies []*unstructured.Unstructured) (string, error) {

	var (
		upstream                        string
		channel                         string
		version                         string
		image                           string
		foundClusterVersionGroupVersion bool
	)

	// validate ClusterVersionGroupVersionKind and keep track to upstream, channel, version, image
	for _, policy := range policies {
		objects, err := r.stripPolicy(policy.Object)
		if err != nil {
			return "", err
		}
		for _, object := range objects {
			kind := object["kind"]
			switch kind {
			case utils.ClusterVersionGroupVersionKind().Kind:
				_, foundSpec := object["spec"]
				if !foundSpec || object["spec"] == nil {
					continue
				}

				if object["spec"].(map[string]interface{})["upstream"] != nil {
					nextUpstream := object["spec"].(map[string]interface{})["upstream"].(string)
					if upstream == "" {
						upstream = nextUpstream
					} else if upstream != nextUpstream {
						return "", errors.New("platform image defined more then once with conflicting upstream values")
					}
				}

				if object["spec"].(map[string]interface{})["channel"] != nil {
					nextChannel := object["spec"].(map[string]interface{})["channel"].(string)
					if channel == "" {
						channel = nextChannel
					} else if channel != nextChannel {
						return "", errors.New("platform image defined more then once with conflicting channel values")
					}
				}

				if object["spec"].(map[string]interface{})["desiredUpdate"] != nil {
					if object["spec"].(map[string]interface{})["desiredUpdate"].(map[string]interface{})["version"] != nil {
						nextVersion := object["spec"].(map[string]interface{})["desiredUpdate"].(map[string]interface{})["version"].(string)
						if version == "" {
							version = nextVersion
						} else if version != nextVersion {
							return "", errors.New("platform image defined more then once with conflicting version values")
						}
					}
					if object["spec"].(map[string]interface{})["desiredUpdate"].(map[string]interface{})["image"] != nil {
						nextImage := object["spec"].(map[string]interface{})["desiredUpdate"].(map[string]interface{})["image"].(string)
						if image == "" {
							image = nextImage
						} else if image != nextImage {
							return "", errors.New("platform image defined more then once with conflicting image values")
						}
					}
				}

				foundClusterVersionGroupVersion = true
			default:
				continue
			}
		}
	}

	if !foundClusterVersionGroupVersion {
		return "", nil
	}

	// return early if policy is valid and user provided .Spec.DesiredUpdate.image in ClusterVersion
	if image != "" {
		return image, nil
	}

	// check for all the required variables needed to make http call and retrieve image
	if upstream == "" || channel == "" || version == "" {
		return "", errors.New("policy with ClusterVersion must have upstream, channel, and version when image is not provided")
	}
	image, err := r.getImageForVersionFromUpdateGraph(upstream, channel, version)
	if err != nil {
		return "", err
	}
	if image == "" {
		return "", errors.New("unable to find platform image for specified upstream, channel, and version")
	}

	return image, nil
}

func (r *ClusterGroupUpgradeReconciler) validateOpenshiftUpgradeVersion(
	clusterGroupUpgrade *ranv1alpha1.ClusterGroupUpgrade, policies []*unstructured.Unstructured) error {

	_, err := r.extractOpenshiftImagePlatformFromPolicies(policies)
	if err != nil {
		utils.SetStatusCondition(
			&clusterGroupUpgrade.Status.Conditions,
			utils.ConditionTypes.Validated,
			utils.ConditionReasons.InvalidPlatformImage,
			metav1.ConditionFalse,
			err.Error(),
		)
		return err
	}

	return nil
}
