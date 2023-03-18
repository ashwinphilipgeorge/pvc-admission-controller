package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", handleValidate)
	log.Fatal(http.ListenAndServeTLS(":443", "/certs/tls.crt", "/certs/tls.key", mux))
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	var admissionReview admission.AdmissionReview

	body, err := readRequestBody(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %v", err), http.StatusBadRequest)
		return
	}

	if _, _, err := deserializer.Decode(body, nil, &admissionReview); err != nil {
		http.Error(w, fmt.Sprintf("failed to deserialize request: %v", err), http.StatusBadRequest)
		return
	}

	admissionResponse := validate(admissionReview.Request)
	admissionReview.Response = admissionResponse

	res, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func validate(req *admission.AdmissionRequest) *admission.AdmissionResponse {
	var pvc corev1.PersistentVolumeClaim

	if err := json.Unmarshal(req.Object.Raw, &pvc); err != nil {
		return &admission.AdmissionResponse{
			UID:     req.UID,
			Allowed: false,
			Result: &metav1.Status{
				Message: fmt.Sprintf("failed to unmarshal PVC object: %v", err),
				Code:    http.StatusBadRequest,
			},
		}
	}

	size := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
	maxSize := resource.MustParse("10Gi")

	if size.Cmp(maxSize) > 0 {
		return &admission.AdmissionResponse{
			UID:     req.UID,
			Allowed: false,
			Result: &metav1.Status{
				Message: "PVC size exceeds 10GB limit",
				Code:    http.StatusForbidden,
			},
		}
	}

	return &admission.AdmissionResponse{
		UID:     req.UID,
		Allowed: true,
	}
}

func readRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("Request body is empty")
	}

	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}
