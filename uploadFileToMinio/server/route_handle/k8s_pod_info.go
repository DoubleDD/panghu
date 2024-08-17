package routehandle

import (
	"minioUploadFile/k8s"
	"minioUploadFile/server/common"
	"net/http"
)

func K8sPodInfo(w http.ResponseWriter, r *http.Request) {
	pod := k8s.PodInfo()
	common.OK(w, pod)
}
