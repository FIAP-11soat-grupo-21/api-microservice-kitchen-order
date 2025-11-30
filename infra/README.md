# Descrição

Este diretório contém as informações para realizar a implantação de uma aplicação no ECS da AWS utilizando o Terraform.
Com essa configuração também é possível realizar a exposição do serviço no API Gateway da AWS.

# Pré-requisitos

- Conta na AWS com permissões para criar recursos no ECS, API Gateway, IAM, VPC, entre outros.
- Terraform instalado na máquina local.
- AWS CLI configurado com as credenciais da conta AWS.

> [!IMPORTANT]
> Certifique-se de que os serviços base necessários, como VPC, sub-redes e grupos de segurança, estejam configurados antes de iniciar a implantação. A implantação deles deve ser feita em [Infra Core](https://github.com/FIAP-11soat-grupo-21/infra-core)

# Estrutura do Diretório
- `main.tf`: Arquivo principal do Terraform que define os recursos a serem criados.
- `variables.tf`: Define as variáveis utilizadas na configuração do Terraform.
- `data.tf`: Define os data sources utilizados na configuração do Terraform.
- `providers.tf`: Configura os provedores necessários para a implantação.

# Instruções de Implantação
1. Realize a configuração do provider AWS no arquivo `providers.tf`, especificando a região desejada.
``` hcl
provider "aws" {
  region = "us-east-2"
}

terraform {
  backend "s3" {
    bucket = "fiap-tc-terraform-846874" //Bucket S3 onde o arquivo de estado será armazenado
    key    = "tech-challenge-project/<MICROSSERVICE>/terraform.tfstate" //Esse é o caminho dentro do bucket onde o arquivo de estado será armazenado
    region = "us-east-2"
  }
}
```
> [!IMPORTANT]
> Substitua `<MICROSSERVICE>` pelo nome do microsserviço que está sendo implantado (por exemplo, `customer`).
2. Defina quais recursos serão criados no arquivo `main.tf`. Para maiores informações dos módulos que podem ser utilizados, consulte a [documentação dos módulos](https://github.com/FIAP-11soat-grupo-21/infra-core/tree/main/modules).
3. Crie um arquivo `terraform.tfvars` para definir os valores das variáveis utilizadas na configuração do Terraform.
4. Inicialize o Terraform para baixar os provedores e módulos necessários:
``` bash
terraform init
```
5. Verifique o plano de execução para garantir que os recursos serão criados conforme esperado:
``` bash
terraform plan
```
6. Aplique a configuração para criar os recursos na AWS:
``` bash
terraform apply
```
7. Confirme a aplicação quando solicitado digitando `yes`.
8. Após a conclusão, verifique no console da AWS se os recursos foram criados corretamente.
9. Para destruir os recursos criados, utilize o comando:
``` bash
terraform destroy
```
