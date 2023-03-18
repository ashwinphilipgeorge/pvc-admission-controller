# PVC Admission Controller

This repository contains the source code and Kubernetes manifests for a custom validating admission controller that restricts the size of Persistent Volume Claims (PVCs) in a Kubernetes cluster. The admission controller will reject any PVC creation or update request if the requested size is greater than 10GB.

To understand more about admission controllers, you can following along the related Medium article [here](https://medium.com/@ashwinphilip96/kubernetes-admission-controllers-enhance-security-and-ensure-compliance-6b61e85d6f24)

## Quickstart

Deploy the latest version of cert-manager into your cluster:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
```

Clone this repository to your local machine

```bash
git clone https://github.com/your-repo/pvc-admission-controller.git
cd pvc-admission-controller
```

Deploy the manifest files with the following command:

```bash
kubectl apply -f manifests
```

Wait for all the controller components to be in a ready and running state, and youre done! The controller will now capture PVC requests to the API server and deny them if theyre greater than 10GB.

You can test if your deployments are working wit hthe following command.

```bash
kubectl apply -f tests
```
## Optimizations

The provided admission controller is a simple example that demonstrates the process of implementing and deploying a custom admission controller to enforce PVC size restrictions. Various optimizations and improvements can be made to enhance the functionality and security of the controller. Check the **Optimizations** section of the Medium article for more information on possible enhancements.

