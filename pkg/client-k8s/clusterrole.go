package clientk8s

import (
	"context"
	"log"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateClusterRole ...
func CreateClusterRole(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, rules []Rbacv1PolicyRule) error {

	context := context.Background()

	var policyRules []rbacv1.PolicyRule

	for _, item := range rules {
		policyRules = append(policyRules, rbacv1.PolicyRule{
			Verbs:           item.Verbs,
			APIGroups:       item.APIGroups,
			Resources:       item.Resources,
			ResourceNames:   item.ResourceNames,
			NonResourceURLs: item.NonResourceURLs,
		})
	}

	roleSpec := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       typeMeta.Kind,
			APIVersion: typeMeta.APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        objectMeta.Name,
			Namespace:   objectMeta.Namespace,
			Labels:      objectMeta.Labels,
			Annotations: objectMeta.Annotations,
		},
		Rules: policyRules,
	}

	_, err = k8sClientSet.RbacV1().ClusterRoles().Create(context, roleSpec, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating cluster role: %v \n", err)
		return err
	}

	return nil

}

func GetClusterRole(name string) (*rbacv1.ClusterRole, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().ClusterRoles().Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting cluster role: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateClusterRole(objClusterRole *rbacv1.ClusterRole, rules []Rbacv1PolicyRule) error {

	context := context.Background()

	var policyRules []rbacv1.PolicyRule

	for _, item := range rules {
		policyRules = append(policyRules, rbacv1.PolicyRule{
			Verbs:           item.Verbs,
			APIGroups:       item.APIGroups,
			Resources:       item.Resources,
			ResourceNames:   item.ResourceNames,
			NonResourceURLs: item.NonResourceURLs,
		})
	}

	objClusterRole.Rules = policyRules

	_, err := k8sClientSet.RbacV1().ClusterRoles().Update(context, objClusterRole, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating cluster role: %v \n", err)
		return err
	}

	return nil

}

func ListClusterRole() (*rbacv1.ClusterRoleList, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().ClusterRoles().List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list clusters role: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteClusterRole(name string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.RbacV1().ClusterRoles().Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete cluster role: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateClusterRole(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, rules []Rbacv1PolicyRule) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetClusterRole(objectMeta.Name)
		if getErr != nil {
			log.Printf("Error getting cluster role: %v \n", err)
			// return err
		}

		if resultGet != nil {
			err := UpdateClusterRole(resultGet, rules)
			if err != nil {
				log.Printf("Error updating cluster role: %v \n", err)
				return err
			}
		} else {
			err := CreateClusterRole(typeMeta, objectMeta, rules)
			if err != nil {
				log.Printf("Error creating cluster role: %v \n", err)
				return err
			}
		}
		return nil
	})
	if retryErr != nil {
		return retryErr
	}
	return nil
}
