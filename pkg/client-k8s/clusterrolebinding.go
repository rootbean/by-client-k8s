package clientk8s

import (
	"context"
	"log"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateClusterRoleBinding ...
func CreateClusterRoleBinding(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, subject []Rbacv1Subject, roleRef Rbacv1RoleRef) error {

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

	clusterRoleBindingSpec := &rbacv1.ClusterRoleBinding{
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

	_, err = k8sClientSet.RbacV1().ClusterRoleBindings().Create(context, clusterRoleBindingSpec, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating cluster role binding: %v \n", err)
		return err
	}

	return nil

}

func GetClusterRoleBinding(name string) (*rbacv1.ClusterRoleBinding, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().ClusterRoleBindings().Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting cluster role binding: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateClusterRoleBinding(objClusterRoleBinding *rbacv1.ClusterRoleBinding, subject []Rbacv1Subject) error {

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

	objClusterRoleBinding.Subjects = subjectItems

	_, err := k8sClientSet.RbacV1().ClusterRoleBindings().Update(context, objClusterRoleBinding, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating cluster role binding: %v \n", err)
		return err
	}

	return nil

}

func ListClusterRoleBinding() (*rbacv1.ClusterRoleBindingList, error) {

	context := context.Background()

	result, err := k8sClientSet.RbacV1().ClusterRoleBindings().List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list clusters role bindings: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteClusterRoleBinding(name string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.RbacV1().ClusterRoleBindings().Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete cluster role binding: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateClusterRoleBinding(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, subject []Rbacv1Subject, roleRef Rbacv1RoleRef) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetClusterRoleBinding(objectMeta.Name)
		if getErr != nil {
			log.Printf("Error getting cluster role binding: %v \n", err)
			// return err
		}

		if resultGet != nil {
			err := CreateClusterRoleBinding(typeMeta, objectMeta, subject, roleRef)
			if err != nil {
				log.Printf("Error creating cluster role binding: %v \n", err)
				return err
			}
		} else {
			err := UpdateClusterRoleBinding(resultGet, subject)
			if err != nil {
				log.Printf("Error updating cluster role binding: %v \n", err)
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
