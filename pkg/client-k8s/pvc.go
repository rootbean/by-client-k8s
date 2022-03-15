package clientk8s

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func CreatePVC(
	typeMeta Metav1TypeMeta,
	objectMeta Metav1ObjectMeta,
	volumeAccessMode PersistentVolumeAccessMode,
	storageClassName string,
	resourceMustParse string,
) error {

	context := context.Background()

	var persistentVolumeAccessModeItems []corev1.PersistentVolumeAccessMode

	if volumeAccessMode.ReadWriteOnce {
		persistentVolumeAccessModeItems = append(persistentVolumeAccessModeItems, corev1.ReadWriteOnce)
	}

	if volumeAccessMode.ReadOnlyMany {
		persistentVolumeAccessModeItems = append(persistentVolumeAccessModeItems, corev1.ReadOnlyMany)
	}

	if volumeAccessMode.ReadWriteMany {
		persistentVolumeAccessModeItems = append(persistentVolumeAccessModeItems, corev1.ReadWriteMany)
	}

	pvcSpec := corev1.PersistentVolumeClaimSpec{
		AccessModes: persistentVolumeAccessModeItems,
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(resourceMustParse),
			},
		},
		StorageClassName: &storageClassName,
	}

	pvc := corev1.PersistentVolumeClaim{
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
		Spec: pvcSpec,
	}

	_, err = k8sClientSet.CoreV1().PersistentVolumeClaims(objectMeta.Namespace).Create(context, &pvc, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Error creating pvc: %v \n", err)
		return err
	}

	return nil

}

func GetPVC(name, namespace string) (*v1.PersistentVolumeClaim, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context, name, metav1.GetOptions{})

	if err != nil {
		log.Printf("Error getting PVC: %v \n", err)
		return nil, err
	}

	return result, nil

}

func UpdatePVC(
	objPVC *v1.PersistentVolumeClaim,
	volumeAccessMode PersistentVolumeAccessMode,
	storageClassName string,
	resourceMustParse string,
) error {

	context := context.Background()

	var persistentVolumeAccessModeItems []corev1.PersistentVolumeAccessMode

	if volumeAccessMode.ReadWriteOnce {
		persistentVolumeAccessModeItems = append(persistentVolumeAccessModeItems, corev1.ReadWriteOnce)
	}

	if volumeAccessMode.ReadOnlyMany {
		persistentVolumeAccessModeItems = append(persistentVolumeAccessModeItems, corev1.ReadOnlyMany)
	}

	if volumeAccessMode.ReadWriteMany {
		persistentVolumeAccessModeItems = append(persistentVolumeAccessModeItems, corev1.ReadWriteMany)
	}

	pvcSpec := corev1.PersistentVolumeClaimSpec{
		AccessModes: persistentVolumeAccessModeItems,
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(resourceMustParse),
			},
		},
		StorageClassName: &storageClassName,
	}

	objPVC.Spec = pvcSpec

	_, err := k8sClientSet.CoreV1().PersistentVolumeClaims(objPVC.ObjectMeta.Namespace).Update(context, objPVC, metav1.UpdateOptions{})

	if err != nil {
		log.Printf("Error updating PVC: %v \n", err)
		return err
	}

	return nil

}

func ListPVC(namespace string) (*v1.PersistentVolumeClaimList, error) {

	context := context.Background()

	result, err := k8sClientSet.CoreV1().PersistentVolumeClaims(namespace).List(context, metav1.ListOptions{})

	if err != nil {
		log.Printf("Error list PVCs: %v \n", err)
		return nil, err
	}

	return result, nil

}

func DeletePVC(name, namespace string) error {

	context := context.Background()

	deletePolicy := metav1.DeletePropagationForeground

	err := k8sClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		log.Printf("Error delete PVC: %v \n", err)
		return err
	}

	return nil

}

func CreateOrUpdatePVC(
	typeMeta Metav1TypeMeta,
	objectMeta Metav1ObjectMeta,
	volumeAccessMode PersistentVolumeAccessMode,
	storageClassName string,
	resourceMustParse string,
) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		resultGet, getErr := GetPVC(objectMeta.Name, objectMeta.Namespace)
		if getErr != nil {
			log.Printf("Error getting PVC: %v \n", err)
			return err
		}

		if resultGet != nil {
			err := CreatePVC(typeMeta, objectMeta, volumeAccessMode, storageClassName, resourceMustParse)
			if err != nil {
				log.Printf("Error creating PVC: %v \n", err)
				return err
			}
		} else {
			err := UpdatePVC(resultGet, volumeAccessMode, storageClassName, resourceMustParse)
			if err != nil {
				log.Printf("Error updating PVC: %v \n", err)
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
