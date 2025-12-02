 .PHONY: down up
down:
	podman kube play --down kube.yaml

up:
	podman kube play --build kube.yaml