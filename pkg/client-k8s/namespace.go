package clientk8s

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta) error {

	context := context.Background()

	nsSpec := &v1.Namespace{
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
	}

	_, err = k8sClientSet.CoreV1().Namespaces().Create(context, nsSpec, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating namespace: %v \n", err)
		return err
	}

	return nil

}

func GetNamespace(name string) (*v1.Namespace, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().Namespaces().Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting Namespace: %v \n", err)
		return nil, err
	}

	return result, nil

}

func ListNamespace() (*v1.NamespaceList, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().Namespaces().List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list Namespaces: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteNamespace(name string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.CoreV1().Namespaces().Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete Namespace: %v \n", err)
		return err
	}

	return nil

}
