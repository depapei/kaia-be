package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GenerateModels() {
	// Generate Models
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_TIMEZONE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Initialize the generator
	g := gen.NewGenerator(gen.Config{
		OutPath: "./models",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	// Generate structs from all tables of the current database
	// The generated models will be placed in the specified OutPath
	users := g.GenerateModel("users")

	admins := g.GenerateModel("admins")

	// productslices := g.GenerateModel("productslice")

	// products := g.GenerateModel("products",
	// 	gen.FieldRelate(field.HasOne, "Admin", admins, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"created_by"}},
	// 	}),
	// 	gen.FieldRelate(field.HasMany, "ProductSlices", productslices, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"product_id"}},
	// 	}),
	// )

	// wishlists := g.GenerateModel("wishlist",
	// 	gen.FieldRelate(field.HasOne, "Product", products, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"product_id"}},
	// 	}),
	// 	gen.FieldRelate(field.HasOne, "User", users, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"user_id"}},
	// 	}),
	// )

	// detail_transaction := g.GenerateModel("detail_transaction",
	// 	gen.FieldRelate(field.HasOne, "ProductSlice", productslices, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"productslice_id"}},
	// 	}),
	// )

	// transaction := g.GenerateModelAs("transactions", "HeaderTransaction",
	// 	gen.FieldRelate(field.HasMany, "DetailTransaction", detail_transaction, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"transaction_id"}},
	// 	}),
	// 	gen.FieldRelate(field.HasOne, "Customer", users, &field.RelateConfig{
	// 		GORMTag: field.GormTag{"foreignKey": []string{"created_by"}},
	// 	}),
	// )

	g.ApplyBasic(
		users,
		admins,
		// products,
		// wishlists,
		// productslices,
		// detail_transaction,
		// transaction
	)

	// g.GenerateAllTable()

	// Execute the generation
	g.Execute()
}
