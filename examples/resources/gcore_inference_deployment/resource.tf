resource "gcore_inference_deployment" "inf" {
        project_id = 184550
        name = "terra-inf"
        image = "nginx:latest"
        listening_port = 80
        flavor_name = "inference-1vcpu-1gib"
        timeout = 60
        containers {
                region_id = 4
                cooldown_period = 60
                scale_min = 2
                scale_max = 2
                triggers_cpu_threshold = 80
        }
}