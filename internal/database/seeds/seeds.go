package seeds

import (
	"log"
	"money_map/internal/modules/user"

	"gorm.io/gorm"

	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// SeedUsers insere usuários de exemplo no banco de dados
func SeedUsers(db *gorm.DB) error {
	// Verificar se já existem usuários
	var count int64
	db.Model(&user.User{}).Count(&count)
	if count > 0 {
		log.Println("Usuários já existem no banco, pulando seed...")
		return nil
	}

	// Hash das senhas
	password1, _ := HashPassword("senha123")
	password2, _ := HashPassword("senha456")
	password3, _ := HashPassword("senha789")

	// Criar usuários de exemplo
	users := []user.User{
		{
			ID:       cuid.New(),
			Email:    "joao.silva@example.com",
			Password: password1,
			Name:     "João Silva",
		},
		{
			ID:       cuid.New(),
			Email:    "maria.santos@example.com",
			Password: password2,
			Name:     "Maria Santos",
		},
		{
			ID:       cuid.New(),
			Email:    "pedro.oliveira@example.com",
			Password: password3,
			Name:     "Pedro Oliveira",
		},
	}

	// Inserir usuários no banco
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Erro ao criar usuário %s: %v", user.Email, err)
			return err
		}
		log.Printf("Usuário criado: %s (%s)", user.Name, user.Email)
	}

	log.Println("Seed de usuários concluída com sucesso!")
	return nil
}

// SeedAll executa todas as seeds
func SeedAll(db *gorm.DB) error {
	log.Println("Iniciando seed do banco de dados...")

	if err := SeedUsers(db); err != nil {
		return err
	}

	log.Println("Todas as seeds foram executadas com sucesso!")
	return nil
}
