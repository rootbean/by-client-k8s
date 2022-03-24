package clientk8s

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateSecret ...
func CreateSecret(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, typeSecret SecretTypeStruct, data map[string][]byte) error {

	context := context.Background()

	var typeSecretSelected corev1.SecretType

	switch typeSecret.SecretType {
	case "SecretTypeBasicAuth":
		typeSecretSelected = corev1.SecretTypeBasicAuth
	case "SecretTypeBootstrapToken":
		typeSecretSelected = corev1.SecretTypeBootstrapToken
	case "SecretTypeDockerConfigJson":
		typeSecretSelected = corev1.SecretTypeDockerConfigJson
	case "SecretTypeDockercfg":
		typeSecretSelected = corev1.SecretTypeDockercfg
	case "SecretTypeOpaque":
		typeSecretSelected = corev1.SecretTypeOpaque
	case "SecretTypeSSHAuth":
		typeSecretSelected = corev1.SecretTypeSSHAuth
	case "SecretTypeServiceAccountToken":
		typeSecretSelected = corev1.SecretTypeServiceAccountToken
	case "SecretTypeTLS":
		typeSecretSelected = corev1.SecretTypeTLS
	}

	secret := corev1.Secret{
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
		Type: typeSecretSelected,
		Data: data,
	}

	_, err = k8sClientSet.CoreV1().Secrets(objectMeta.Namespace).Create(context, &secret, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Failed to create the secret: %v \n", err)
		return err
	}

	return nil

}

func GetSecret(name, namespace string) (*v1.Secret, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().Secrets(namespace).Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting Secret: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateSecret(objSecret *v1.Secret, typeSecret SecretTypeStruct, data map[string][]byte) error {

	context := context.Background()

	var typeSecretSelected corev1.SecretType

	switch typeSecret.SecretType {
	case "SecretTypeBasicAuth":
		typeSecretSelected = corev1.SecretTypeBasicAuth
	case "SecretTypeBootstrapToken":
		typeSecretSelected = corev1.SecretTypeBootstrapToken
	case "SecretTypeDockerConfigJson":
		typeSecretSelected = corev1.SecretTypeDockerConfigJson
	case "SecretTypeDockercfg":
		typeSecretSelected = corev1.SecretTypeDockercfg
	case "SecretTypeOpaque":
		typeSecretSelected = corev1.SecretTypeOpaque
	case "SecretTypeSSHAuth":
		typeSecretSelected = corev1.SecretTypeSSHAuth
	case "SecretTypeServiceAccountToken":
		typeSecretSelected = corev1.SecretTypeServiceAccountToken
	case "SecretTypeTLS":
		typeSecretSelected = corev1.SecretTypeTLS
	}

	objSecret.Type = typeSecretSelected
	objSecret.Data = data

	_, err := k8sClientSet.CoreV1().Secrets(objSecret.ObjectMeta.Namespace).Update(context, objSecret, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating Secret: %v \n", err)
		return err
	}

	return nil

}

func ListSecret(namespace string) (*v1.SecretList, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().Secrets(namespace).List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list Secrets: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteSecret(name, namespace string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.CoreV1().Secrets(namespace).Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete Secret: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateSecret(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, typeSecret SecretTypeStruct, data map[string][]byte) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetSecret(objectMeta.Name, objectMeta.Namespace)
		if getErr != nil {
			log.Printf("Error getting Secret: %v \n", err)
			// return err
		}

		if resultGet != nil {
			err := UpdateSecret(resultGet, typeSecret, data)
			if err != nil {
				log.Printf("Error updating Secret: %v \n", err)
				return err
			}
		} else {
			err := CreateSecret(typeMeta, objectMeta, typeSecret, data)
			if err != nil {
				log.Printf("Error creating Secret: %v \n", err)
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
