variable "application_name" {
  description = "Nome da aplicação ECS"
  type        = string
}

variable "image_name" {
  description = "Nome da imagem do container"
  type        = string
}

variable "image_port" {
  description = "Porta do container"
  type        = number
}

variable "desired_count" {
  description = "Número desejado de tarefas ECS"
  type        = number
  default     = 1
}

variable "container_environment_variables" {
  description = "Variáveis de ambiente do container"
  type        = map(string)
  default     = {}
}

variable "container_secrets" {
  description = "Segredos do container"
  type        = map(string)
  default     = {}
}

variable "health_check_path" {
  description = "Caminho de verificação de integridade do serviço"
  type        = string
  default     = "/health"
}

variable "task_role_policy_arns" {
  description = "Lista de ARNs de políticas para anexar à função da tarefa ECS"
  type        = list(string)
  default     = []
}

variable "alb_is_internal" {
  description = "Se o ALB é interno"
  type        = bool
  default     = true
}

variable "app_path_pattern" {
  description = "Lista de padrões de caminho para o listener rule do ALB"
  type        = list(string)
}

#########################################################
############### Variáveis do API Gateway ################
#########################################################

variable "apigw_integration_type" {
  description = "Tipo de integração do API Gateway"
  type        = string
  default     = "HTTP_PROXY"
}

variable "apigw_integration_method" {
  description = "Método de integração do API Gateway"
  type        = string
  default     = "ANY"
}

variable "apigw_payload_format_version" {
  description = "Versão do payload do API Gateway"
  type        = string
  default     = "1.0"
}

variable "apigw_connection_type" {
  description = "Tipo de conexão do API Gateway"
  type        = string
  default     = "VPC_LINK"
}

variable "authorization_name" {
  description = "Link do authorizer do API Gateway"
  type        = string
}

#########################################################
################## Variáveis do SQS ####################
#########################################################

variable "sqs_delay_seconds" {
  description = "Tempo em segundos que a entrega de todas as mensagens na fila será atrasada"
  type        = number
  default     = 0
}

variable "sqs_message_retention_seconds" {
  description = "Número de segundos que o Amazon SQS retém uma mensagem"
  type        = number
  default     = 86400 # 1 dia
}

variable "sqs_receive_wait_time_seconds" {
  description = "Tempo que uma chamada ReceiveMessage aguardará por uma mensagem chegar (long polling)"
  type        = number
  default     = 10
}

variable "sqs_visibility_timeout_seconds" {
  description = "Timeout de visibilidade para a fila"
  type        = number
  default     = 30
}