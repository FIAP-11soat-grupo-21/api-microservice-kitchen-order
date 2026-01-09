output "sqs_kitchen_orders_queue_url" {
  description = "URL da fila SQS do Kitchen Orders (do infra-core)"
  value       = data.terraform_remote_state.infra.outputs.sqs_kitchen_orders_queue_url
}

output "sqs_kitchen_orders_order_error_queue_url" {
  description = "URL da fila SQS de erro do Kitchen Orders (do infra-core)"
  value       = data.terraform_remote_state.infra.outputs.sqs_kitchen_orders_order_error_queue_url
}

output "ecs_service_id" {
  description = "ID do serviço ECS do Kitchen Orders"
  value       = module.kitchen_order_api.service_id
}