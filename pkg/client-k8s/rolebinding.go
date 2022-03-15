package clientk8s

import (
	"context"
	"log"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateRoleBinding ...
func CreateRoleBinding(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, subject []Rbacv1Subject, roleRef Rbacv1RoleRef) error {

	context := context.Background()

	var subjectItems []rbacv1.Subject

	for _, item := range subject {
		subjectItems = append(subjectItems, rbacv1.Subject{
			Kind:      item.Kind,
			APIGroup:  item.APIGroup,
			Name:      item.Name,
			Namespace: item.Namespace,
		})
	}

	roleBindingSpec := &rbacv1.RoleBinding{
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
		Subjects: subjectItems,
		RoleRef: rbacv1.RoleRef{
			APIGroup: roleRef.APIGroup,
			Kind:     roleRef.Kind,
			Name:     roleRef.Name,
		},
	}

	_, err = k8sClientSet.RbacV1().RoleBindings(objectMeta.Namespace).Create(context, roleBindingSpec, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating role binding: %v \n", err)
		return err
	}

	return nil

}

func GetRoleBinding(name, namespace string) (*rbacv1.RoleBinding, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().RoleBindings(namespace).Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting role binding: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateRoleBinding(objRoleBinding *rbacv1.RoleBinding, subject []Rbacv1Subject) error {

	context := context.Background()

	var subjectItems []rbacv1.Subject

	for _, item := range subject {
		subjectItems = append(subjectItems, rbacv1.Subject{
			Kind:      item.Kind,
			APIGroup:  item.APIGroup,
			Name:      item.Name,
			Namespace: item.Namespace,
		})
	}

	objRoleBinding.Subjects = subjectItems

	_, err := k8sClientSet.RbacV1().RoleBindings(objRoleBinding.ObjectMeta.Namespace).Update(context, objRoleBinding, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating role binding: %v \n", err)
		return err
	}

	return nil

}

func ListRoleBinding(namespace string) (*rbacv1.RoleBindingList, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().RoleBindings(namespace).List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list role bindings: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteRoleBinding(name, namespace string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.RbacV1().RoleBindings(namespace).Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete role binding: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateRoleBinding(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, subject []Rbacv1Subject, roleRef Rbacv1RoleRef) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetRoleBinding(objectMeta.Name, objectMeta.Namespace)
		if getErr != nil {
			log.Printf("Error getting role binding: %v \n", err)
			// return err
		}

		if resultGet != nil {
			err := CreateRoleBinding(typeMeta, objectMeta, subject, roleRef)
			if err != nil {
				log.Printf("Error creating role binding: %v \n", err)
				return err
			}
		} else {
			err := UpdateRoleBinding(resultGet, subject)
			if err != nil {
				log.Printf("Error updating role binding: %v \n", err)
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
