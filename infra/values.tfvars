application_name = "kitchen-order-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8082
app_path_pattern = ["/kitchen-orders*", "/kitchen-orders/*"]

# =======================================================
# Configurações do ECS Service
# =======================================================
container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8082"
  API_HOST : "0.0.0.0"
  AWS_REGION : "us-east-2"

  DB_RUN_MIGRATIONS : "true"
  DB_NAME : "postgres"
  DB_PORT : "5432"
  DB_USERNAME : "postgres"

  MESSAGE_BROKER_TYPE : "sqs"
}

container_secrets = {}

health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonRDSFullAccess",
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
  "arn:aws:iam::aws:policy/AmazonCognitoPowerUser",
]
alb_is_internal = true

# =======================================================
# Configurações do API Gateway
# =======================================================
apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"

authorization_name = "CognitoAuthorizer"

# =======================================================
# Configurações do SQS
# =======================================================
sqs_delay_seconds              = 0
sqs_message_retention_seconds  = 86400 # 1 dia
sqs_receive_wait_time_seconds  = 10
sqs_visibility_timeout_seconds = 30

