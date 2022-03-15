package clientk8s

import (
	"context"
	"log"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateRole ...
func CreateRole(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, rules []Rbacv1PolicyRule) error {

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

	roleSpec := &rbacv1.Role{
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

	_, err = k8sClientSet.RbacV1().Roles(objectMeta.Namespace).Create(context, roleSpec, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating role: %v \n", err)
		return err
	}

	return nil

}

func GetRole(name, namespace string) (*rbacv1.Role, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().Roles(namespace).Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting role: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateRole(objRole *rbacv1.Role, rules []Rbacv1PolicyRule) error {

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

	objRole.Rules = policyRules

	_, err := k8sClientSet.RbacV1().Roles(objRole.ObjectMeta.Namespace).Update(context, objRole, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating role: %v \n", err)
		return err
	}

	return nil

}

func ListRole(namespace string) (*rbacv1.RoleList, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().Roles(namespace).List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list roles: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteRole(name, namespace string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.RbacV1().Roles(namespace).Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete role: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateRole(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, rules []Rbacv1PolicyRule) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetRole(objectMeta.Name, objectMeta.Namespace)
		if getErr != nil {
			log.Printf("Error getting role: %v \n", err)
			return err
		}

		if resultGet != nil {
			err := CreateRole(typeMeta, objectMeta, rules)
			if err != nil {
				log.Printf("Error creating role: %v \n", err)
				return err
			}
		} else {
			err := UpdateRole(resultGet, rules)
			if err != nil {
				log.Printf("Error updating role: %v \n", err)
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
