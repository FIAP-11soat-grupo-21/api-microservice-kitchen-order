module "kitchen_order_api" {
  source     = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/ECS-Service?ref=main"
  depends_on = [aws_lb_listener.listener]

  cluster_id            = data.terraform_remote_state.infra.outputs.ecs_cluster_id
  ecs_security_group_id = data.terraform_remote_state.infra.outputs.ecs_security_group_id

  cloudwatch_log_group     = data.terraform_remote_state.infra.outputs.ecs_cloudwatch_log_group
  ecs_container_image      = var.image_name
  ecs_container_name       = var.application_name
  ecs_container_port       = var.image_port
  ecs_service_name         = var.application_name
  ecs_desired_count        = var.desired_count
  registry_credentials_arn = data.terraform_remote_state.infra.outputs.ecr_registry_credentials_arn

  ecs_container_environment_variables = merge(var.container_environment_variables,
    {
      # Database configuration - usando a mesma instância RDS do infra-core
      DB_HOST : data.terraform_remote_state.infra.outputs.rds_address

      # SQS configuration
      AWS_SQS_KITCHEN_ORDERS_QUEUE : data.terraform_remote_state.infra.outputs.sqs_kitchen_orders_queue_url
      AWS_SQS_ORDERS_QUEUE : data.terraform_remote_state.infra.outputs.sqs_orders_queue_url
      AWS_SQS_KITCHEN_ORDERS_ERROR_QUEUE : data.terraform_remote_state.infra.outputs.sqs_kitchen_orders_order_error_queue_url
  })
  ecs_container_secrets = merge(var.container_secrets,
    {
      DB_PASSWORD : data.terraform_remote_state.infra.outputs.rds_secret_arn
  })

  private_subnet_ids      = data.terraform_remote_state.infra.outputs.private_subnet_id
  task_execution_role_arn = data.terraform_remote_state.infra.outputs.ecs_task_execution_role_arn
  task_role_policy_arns   = var.task_role_policy_arns
  alb_target_group_arn    = aws_alb_target_group.target_group.arn
  alb_security_group_id   = data.terraform_remote_state.infra.outputs.alb_security_group_id

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "GetKitchenOrderAPIRoute" {
  source     = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/API-Gateway-Routes?ref=main"
  depends_on = [module.kitchen_order_api]

  api_id       = data.terraform_remote_state.infra.outputs.api_gateway_id
  alb_proxy_id = aws_apigatewayv2_integration.alb_proxy.id

  endpoints = {
    get_kitchen_order = {
      route_key  = "GET /kitchen-orders/{id}"
      restricted = false
    },
    get_all_kitchen_orders = {
      route_key  = "GET /kitchen-orders"
      restricted = false
    },
    update_kitchen_order = {
      route_key  = "PUT /kitchen-orders/{id}"
      restricted = false
    },
    get_all_status = {
      route_key  = "GET /kitchen-orders/status"
      restricted = false
    },
  }
}