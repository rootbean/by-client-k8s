package clientk8s

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateServiceAccount
func CreateServiceAccount(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, secretsArrStr []string) error {

	context := context.Background()

	secretReferences := []v1.ObjectReference{}

	for _, s := range secretsArrStr {

		secretReferences = append(secretReferences, v1.ObjectReference{
			Name: s,
		})
	}

	specServiceAccount := &v1.ServiceAccount{
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
		Secrets: secretReferences,
	}

	_, err = k8sClientSet.CoreV1().ServiceAccounts(objectMeta.Namespace).Create(context, specServiceAccount, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Failed to create the service account: %v \n", err)
		return err
	}

	return nil

}

func GetServiceAccount(name, namespace string) (*v1.ServiceAccount, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().ServiceAccounts(namespace).Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting ServiceAccount: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateServiceAccount(objServiceAccount *v1.ServiceAccount, secretsArrStr []string) error {

	context := context.Background()

	secretReferences := []v1.ObjectReference{}

	for _, s := range secretsArrStr {

		secretReferences = append(secretReferences, v1.ObjectReference{
			Name: s,
		})
	}

	objServiceAccount.Secrets = secretReferences

	_, err := k8sClientSet.CoreV1().ServiceAccounts(objServiceAccount.ObjectMeta.Namespace).Update(context, objServiceAccount, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating ServiceAccount: %v \n", err)
		return err
	}

	return nil

}

func ListServiceAccount(namespace string) (*v1.ServiceAccountList, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().ServiceAccounts(namespace).List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list ServiceAccounts: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteServiceAccount(name, namespace string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.CoreV1().ServiceAccounts(namespace).Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete ServiceAccount: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateServiceAccount(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, secretsArrStr []string) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetServiceAccount(objectMeta.Name, objectMeta.Namespace)
		if getErr != nil {
			log.Printf("Error getting ServiceAccount: %v \n", err)
			return err
		}

		if resultGet != nil {
			err := CreateServiceAccount(typeMeta, objectMeta, secretsArrStr)
			if err != nil {
				log.Printf("Error creating ServiceAccount: %v \n", err)
				return err
			}
		} else {
			err := UpdateServiceAccount(resultGet, secretsArrStr)
			if err != nil {
				log.Printf("Error updating ServiceAccount: %v \n", err)
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
