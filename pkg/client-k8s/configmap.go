package clientk8s

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func CreateConfigMap(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, data map[string]string) error {

	context := context.Background()

	cmSpec := &v1.ConfigMap{
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
		Data: data,
	}

	_, err = k8sClientSet.CoreV1().ConfigMaps(objectMeta.Namespace).Create(context, cmSpec, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating configmap: %v \n", err)
		return err
	}

	return nil

}

func GetConfigMap(name, namespace string) (*v1.ConfigMap, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().ConfigMaps(namespace).Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting configmap: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdateConfigMap(objConfigMap *v1.ConfigMap, data map[string]string) error {

	context := context.Background()

	objConfigMap.Data = data

	_, err := k8sClientSet.CoreV1().ConfigMaps(objConfigMap.ObjectMeta.Namespace).Update(context, objConfigMap, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating configMap: %v \n", err)
		return err
	}

	return nil

}

func ListConfigMap(namespace string) (*v1.ConfigMapList, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().ConfigMaps(namespace).List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list configMaps: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeleteConfigMap(name, namespace string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.CoreV1().ConfigMaps(namespace).Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete configMap: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdateConfigMap(typeMeta Metav1TypeMeta, objectMeta Metav1ObjectMeta, data map[string]string) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetConfigMap(objectMeta.Name, objectMeta.Namespace)
		if getErr != nil {
			log.Printf("Error getting configMap: %v \n", err)
			// return err
		}

		if resultGet != nil {
			err := CreateConfigMap(typeMeta, objectMeta, data)
			if err != nil {
				log.Printf("Error creating configMap: %v \n", err)
				return err
			}
		} else {
			err := UpdateConfigMap(resultGet, data)
			if err != nil {
				log.Printf("Error updating configMap: %v \n", err)
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
