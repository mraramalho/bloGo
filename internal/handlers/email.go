package handlers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/joho/godotenv"
)

var (
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsRegion          string
	sender             string
	recipient          string
)

func init() {
	// Loads .env variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("[Warning]: Failed to load .env file: %v", err)
	}

	awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion = os.Getenv("AWS_REGION")
	sender = os.Getenv("SES_SENDER")
	recipient = os.Getenv("SES_RECIPIENT")

	// Verifies if the required environment variables are set
	if awsAccessKeyID == "" || awsSecretAccessKey == "" || awsRegion == "" || sender == "" || recipient == "" {
		log.Println("[Warning]: Uma ou mais variáveis de ambiente necessárias para o envio de e-mail não estão definidas")
	}
}

// sendEmail sends e-mails using AWS SES
func sendEmail(name, email, message string) error {

	// Verifies if the required environment variables are set
	if awsAccessKeyID == "" || awsSecretAccessKey == "" || awsRegion == "" || sender == "" || recipient == "" {
		return fmt.Errorf("configuração de e-mail incompleta: verifique as variáveis de ambiente")
	}

	// Criar um provedor de credenciais estático
	credProvider := credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")

	// Carregar a configuração da AWS com credenciais personalizadas
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credProvider)),
	)
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração da AWS: %v", err)
	}

	// Criar o cliente SES
	sesClient := ses.NewFromConfig(cfg)

	// Configurar os detalhes do e-mail
	input := &ses.SendEmailInput{
		Source: &sender, // O e-mail precisa estar verificado no SES
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String("Prospecto via Blog"),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data:    aws.String(fmt.Sprintf("Nome: %s\nEmail: %s\nMensagem: %s", name, email, message)),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	// Enviar o e-mail
	resp, err := sesClient.SendEmail(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("erro ao enviar e-mail: %v", err)
	}

	log.Printf("E-mail enviado com sucesso! \nID: %v\nNome: %s\nEmail: %s\nMensagem: %s\n", *resp.MessageId, name, email, message)
	return nil
}
