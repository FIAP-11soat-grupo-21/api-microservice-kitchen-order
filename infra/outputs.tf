output "sqs_queue_url" {
  description = "URL da fila SQS"
  value       = module.sqs_kitchen_orders.sqs_queue_url
}

output "sqs_queue_arn" {
  description = "ARN da fila SQS"
  value       = module.sqs_kitchen_orders.sqs_queue_arn
}

output "db_address" {
  description = "Endereço do banco de dados RDS do Kitchen Orders"
  value       = module.app_db.db_connection
}

output "db_secret_arn" {
  description = "ARN do segredo do banco de dados RDS do Kitchen Orders"
  value       = module.app_db.db_secret_password_arn
}

output "ecs_service_id" {
  description = "ID do serviço ECS do Kitchen Orders"
  value       = module.kitchen_order_api.service_id
}