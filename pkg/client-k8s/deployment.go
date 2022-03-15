package clientk8s

import (
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateJSONDeployment(
	typeMeta Metav1TypeMeta,
	objectMeta Metav1ObjectMeta,
	deploymentContainer []DeploymentContainerStruct,
	replicas int32,
) *appsv1.Deployment {

	var containerList []apiv1.Container
	var containerPortList []apiv1.ContainerPort
	var envFromList []apiv1.EnvFromSource
	var envList []apiv1.EnvVar
	var volumenMountsList []apiv1.VolumeMount
	var volumenDevicesList []apiv1.VolumeDevice

	for _, item := range deploymentContainer {

		for _, itemPortList := range item.ContainerPorts {
			containerPortList = append(containerPortList, apiv1.ContainerPort{
				Name:          itemPortList.Name,
				HostPort:      itemPortList.HostPort,
				ContainerPort: itemPortList.ContainerPort,
				Protocol:      v1.Protocol(itemPortList.Protocol),
				HostIP:        itemPortList.HostIP,
			})
		}

		for _, itemEnvFromList := range item.ContainerEnvFrom {
			envFromList = append(envFromList, apiv1.EnvFromSource{
				Prefix: itemEnvFromList.Prefix,
				ConfigMapRef: &apiv1.ConfigMapEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: itemEnvFromList.ConfigMapRef,
					},
				},
				SecretRef: &v1.SecretEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: itemEnvFromList.SecretRef,
					},
				},
			})
		}

		for _, itemEnvList := range item.ContainerEnvVar {
			envList = append(envList, apiv1.EnvVar{
				Name:  itemEnvList.Name,
				Value: itemEnvList.Value,
			})
		}

		for _, itemVolumenMount := range item.ContainerVolumeMounts {
			volumenMountsList = append(volumenMountsList, apiv1.VolumeMount{
				Name:      itemVolumenMount.Name,
				ReadOnly:  itemVolumenMount.ReadOnly,
				MountPath: itemVolumenMount.MountPath,
				SubPath:   itemVolumenMount.SubPath,
			})
		}

		for _, itemVolumenDevice := range item.ContainerVolumeDevices {
			volumenDevicesList = append(volumenDevicesList, apiv1.VolumeDevice{
				Name:       itemVolumenDevice.Name,
				DevicePath: itemVolumenDevice.DevicePath,
			})
		}

		resources := v1.ResourceRequirements{}
		if item.ContainerResource.resourcesLimitsCPU != "" || item.ContainerResource.resourcesLimitsMemory != "" {
			resources.Limits = make(v1.ResourceList)
		}
		if item.ContainerResource.resourcesLimitsCPU != "" {
			resources.Limits[v1.ResourceCPU] = resource.MustParse(item.ContainerResource.resourcesLimitsCPU)
		}
		if item.ContainerResource.resourcesLimitsMemory != "" {
			resources.Limits[v1.ResourceMemory] = resource.MustParse(item.ContainerResource.resourcesLimitsMemory)
		}

		if item.ContainerResource.resourcesRequestsCPU != "" || item.ContainerResource.resourcesRequestsMemory != "" {
			resources.Requests = make(v1.ResourceList)
		}
		if item.ContainerResource.resourcesRequestsCPU != "" {
			resources.Requests[v1.ResourceCPU] = resource.MustParse(item.ContainerResource.resourcesRequestsCPU)
		}
		if item.ContainerResource.resourcesRequestsMemory != "" {
			resources.Requests[v1.ResourceMemory] = resource.MustParse(item.ContainerResource.resourcesRequestsMemory)
		}

		containerList = append(containerList, apiv1.Container{
			Name:          item.ContainerName,
			Image:         item.ContainerImage,
			Command:       item.ContainerCommand,
			Args:          item.ContainerArgs,
			WorkingDir:    item.ContainerWorkingDir,
			Ports:         containerPortList,
			EnvFrom:       envFromList,
			Env:           envList,
			VolumeMounts:  volumenMountsList,
			VolumeDevices: volumenDevicesList,
			Resources:     resources,
		})
	}

	deployment := &appsv1.Deployment{
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
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: objectMeta.Labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: objectMeta.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: containerList,
				},
			},
		},
	}

	log.Println("DEPLOYMENT: ", deployment)

	return deployment

}
